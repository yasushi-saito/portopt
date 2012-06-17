//
// Created by Yaz Saito on 06/05/12.
//

package portopt;

import "github.com/kuroneko/gosqlite3"
import "time"
import "os"
import "net/http"
import "log"
import "encoding/csv"
import "regexp"
import "strconv"
import "fmt"

var dateRe *regexp.Regexp

type Database struct {
	Path string
	db *sqlite3.Database
	cachedSecurities map[string]*Security
	correlationCache map[TickerPair]float64
};

// Cache of database entry.
type Security struct {
	Ticker string
	Weight float64

	// The date range for which the priceMap is defined
	priceDateRange *dateRange

	// Date (UNIX time) -> Adjusted closing price, as reported by yahoo.
	priceMap map[int64]float64

	// Cache of previously computed security stats. The key must
	// be a subrange of this.dateRange.
	statsCache map[*dateRange]SecurityStats
};

type TickerPair struct {
	ticker1 string
	ticker2 string
}

type SecurityStats struct {
	Mean float64
	Stddev float64
}

func (db *Database) Stats(ticker string, r *dateRange) (SecurityStats, error) {
	var stats SecurityStats;
	acc := newStatsAccumulator(ticker);

	s1, err := db.FindSecurity(ticker)
	if err != nil { return stats, err }

	stats, found := s1.statsCache[r]
	if found {
		return stats, nil
	}

	for i := r.Begin(); !i.Done(); i.Next() {
		t := i.Time().Unix()
		acc.Add(s1.priceMap[t])
	}
	stats.Mean = acc.Mean()
	stats.Stddev = acc.StdDev()
	s1.statsCache[r] = stats
	return stats, nil
}

func (db *Database) Correlation (ticker1 string, ticker2 string) (float64, error) {
	p := TickerPair{ ticker1: ticker1, ticker2 : ticker2 }
	// if p.ticker1 == p.ticker2 { return 1.0, nil }
	if p.ticker1 > p.ticker2 {
		p.ticker1, p.ticker2 = p.ticker2, p.ticker1
	}
	corr, found := db.correlationCache[p]
	if found { return corr, nil }

	s1, err := db.FindSecurity(ticker1)
	if err != nil { return -1.0, err }

	s2, err := db.FindSecurity(ticker2)
	if err != nil { return -1.0, err }

	dateRange := s1.priceDateRange.Intersect(s2.priceDateRange)
	log.Print("RANGE: ", s1.priceDateRange, s2.priceDateRange, dateRange)
	stats1 := newStatsAccumulator(ticker1)
	stats2 := newStatsAccumulator(ticker2)

	for i := dateRange.Begin(); !i.Done(); i.Next() {
		t := i.Time().Unix()
		stats1.Add(s1.priceMap[t])
		stats2.Add(s2.priceMap[t])
	}

	var diffTotal float64 = 0.0
	for period := 0; period < stats1.NumItems(); period++ {
		diffTotal += stats1.DeltaForPeriod(period) * stats2.DeltaForPeriod(period)
	}
	corr = diffTotal / float64(stats1.NumItems()) / stats1.StdDev() / stats2.StdDev()
	db.correlationCache[p] = corr
	return corr, nil
}

func searchPrice(db *Database,
	ticker string,
	date time.Time,
	interval time.Duration) (float64) {
	price := -1.0
	limitDate := date.Add(interval)
	db.MustRunQuery(fmt.Sprintf(
		"SELECT adjclose from price WHERE ticker = '%s' AND date >= %d AND date < %d ORDER BY date LIMIT 1",
		ticker, date.Unix(), limitDate.Unix()),
		func(val... interface{}) {
		price = val[0].(float64)
	})
	return price
}

func (db *Database) FindSecurity(ticker string) (*Security, error) {
	s, found := db.cachedSecurities[ticker]
	if found {
		// TODO: check staleness
		return s, nil
	}

	now := time.Now()
	r := db.GetDateRange(ticker)
	if r.Empty() || (now.Sub(r.End()) >= 24 * 3600 * 30) {
		log.Print("Filling ", ticker, " from interweb, range=", r.String(), now.Sub(r.End()))
		err := db.fillFromYahoo(ticker)
		if err != nil { return nil, err }
	}
	s = new(Security)
	s.Ticker = ticker
	s.priceMap = make(map[int64]float64)
	s.statsCache = make(map[*dateRange]SecurityStats)

	var minDate time.Time
	var maxDate time.Time
	for i := r.Begin(); !i.Done(); i.Next() {
		price := searchPrice(db, ticker, i.Time(), minInterval)
		if price >= 0.0 {
			if minDate.IsZero() || minDate.After(i.Time()) {
				minDate = i.Time()
			}
			if maxDate.IsZero() || maxDate.Before(i.Time()) {
				maxDate = i.Time()
			}
			s.priceMap[i.Time().Unix()] = price
		} else {
			break;
		}
	}
	s.priceDateRange = NewDateRange(minDate, maxDate, minInterval)
	db.cachedSecurities[ticker] = s
	return s, nil;
}

func PanicOnError(err error, params... interface{}) {
	if err != nil {
		pp := params;
		log.Fatal(append(pp, ": ", err));
		log.Fatal(err)
	}
}

func mustParseDecimal(s string) (int) {
	value, err := strconv.ParseInt(s, 10, 64)
	PanicOnError(err, "Failed to parse ", s, " as a decimal string")
	return int(value)
}

func mustParseFloat(s string) (float64) {
	value, err := strconv.ParseFloat(s, 64)
	PanicOnError(err, "Failed to parse ", s, " as a float string")
	return value
}

func (db Database) MustUpdate(sql string) {
	st, err := db.db.Prepare(sql)
	PanicOnError(err, "failed to prepare SQL: ", sql)
	st.Step()
	PanicOnError(st.Finalize(), "blah")
}

func (db Database) MustPrepare(sql string) (*sqlite3.Statement) {
	st, err := db.db.Prepare(sql)
	PanicOnError(err, "select " + sql)
	return st
}

func (db Database) MustRunQuery(sql string, cb func(val... interface{})) {
	st := db.MustPrepare(sql)
	_, err := st.All(func(st *sqlite3.Statement, val... interface{}) {
		cb(val...)
	})
	PanicOnError(err, "Failed to run query "+ sql)
}

func (db Database) FillFromCsv(path string, ticker string) (error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	r, err := reader.ReadAll()
	if err != nil {
		return err
	}
	db.MustUpdate("BEGIN TRANSACTION");
	for _, line := range r {
		matches := dateRe.FindStringSubmatch(line[0])
		if (matches == nil) {
			continue;
		}
		date := time.Date(mustParseDecimal(matches[1]),
			time.Month(mustParseDecimal(matches[2])),
			mustParseDecimal(matches[3]),
			0, 0, 0, 0, time.UTC)
		sql := fmt.Sprintf("INSERT INTO price values('C', %d, %f, %f, %f, %f, %d, %f)",
			date.Unix(),
			mustParseFloat(line[1]),
			mustParseFloat(line[2]),
			mustParseFloat(line[3]),
			mustParseFloat(line[4]),
			mustParseDecimal(line[5]),
			mustParseFloat(line[6]));
		db.MustUpdate(sql)
	}
	db.MustUpdate("COMMIT TRANSACTION");
	return nil
}

func (db Database) fillFromYahoo(ticker string) (error) {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=00&b=0&c=1980&d=01&e=1&f=2015&g=d&ignore=.csv", ticker)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	r, err := reader.ReadAll()
	if err != nil {
		return err
	}
	db.MustUpdate("BEGIN TRANSACTION");
	for _, line := range r {
		matches := dateRe.FindStringSubmatch(line[0])
		if (matches == nil) {
			continue;
		}
		date := time.Date(mustParseDecimal(matches[1]),
			time.Month(mustParseDecimal(matches[2])),
			mustParseDecimal(matches[3]),
			0, 0, 0, 0, time.UTC)
		sql := fmt.Sprintf("INSERT INTO price values('%s', %d, %f, %f, %f, %f, %d, %f)",
			ticker,
			date.Unix(),
			mustParseFloat(line[1]),
			mustParseFloat(line[2]),
			mustParseFloat(line[3]),
			mustParseFloat(line[4]),
			mustParseDecimal(line[5]),
			mustParseFloat(line[6]));
		db.MustUpdate(sql)
	}
	db.MustUpdate("COMMIT TRANSACTION");
	return nil
}

func (db Database) TableExists(table string) (bool) {
	found := false;
	db.MustRunQuery(
		fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", table),
		func(val... interface{}) {
		if val[0] == table {
			found = true
		}
	})
	return found
}

func (db Database) GetDateRange(ticker string) (*dateRange) {
	minDate := time.Now()
	var maxDate time.Time
	// maxDate is zero by default
	db.MustRunQuery(
		fmt.Sprintf("SELECT MIN(date), MAX(date) FROM price WHERE ticker='%s'", ticker),
		func(val... interface{}) {
		if val[0] == nil {
			// No row found for the ticker
		} else {
			minDate = time.Unix(val[0].(int64), 0)
			maxDate = time.Unix(val[1].(int64), 0)
		}
	})
	return NewDateRange(minDate, maxDate, time.Hour * 24)
}

func (db Database) FillCorrelationIfNecessary(ticker1 string, ticker2 string) (error) {
	return nil
}

func CreateDb(path string) (*Database) {
	sqlite3.Initialize()
	db, err := sqlite3.Open(path, sqlite3.O_CREATE | sqlite3.O_READWRITE);
	if err != nil {
		log.Panic("Failed to create db: ", path, ": ", err)
	}

	var d *Database = new(Database);
	d.Path = path
	d.db = db
	d.cachedSecurities = make(map[string]*Security)
	d.correlationCache = make(map[TickerPair]float64)
	if !d.TableExists("dividend") {
		d.MustUpdate("CREATE TABLE dividend (ticker VARCHAR(10), date INTEGER, dividend REAL)");
	}
	if !d.TableExists("correlation") {
		d.MustUpdate("CREATE TABLE correlation (ticker1 VARCHAR(10), ticker2 VARCHAR(10), corr REAL, lastUpdateDate INTEGER)");
		d.MustUpdate("CREATE INDEX correlation_index ON correlation (ticker1, ticker2)")
	}
	if !d.TableExists("price") {
		d.MustUpdate("CREATE TABLE price (ticker VARCHAR(10), date INTEGER, open REAL, high REAL, low REAL, close REAL, volume INTEGER, adjclose REAL)");
		d.MustUpdate("CREATE INDEX price_index ON price (ticker, date)")
	}
	return d;
}

func ShutdownDb() {
	sqlite3.Shutdown()
}

func init() {
	dateRe = regexp.MustCompile("(\\d\\d\\d\\d)-(\\d\\d)-(\\d\\d)")
}

