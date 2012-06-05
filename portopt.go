package main;

import "os"
import "net/http"
import "log"
import "encoding/csv"
import "github.com/kuroneko/gosqlite3"
import "regexp"
import "strconv"
import "time"
import "fmt"

var dateRe *regexp.Regexp

type Db struct {
	db *sqlite3.Database
};

var db Db;

func PanicOnError(err error, params... interface{}) {
	if err != nil {
		pp := params;
		log.Fatal(append(pp, ": ", err));
		log.Fatal(err)
	}
}

func Download(ticker string) (error) {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=00&b=0&c=1980&d=01&e=1&f=2015&g=d&ignore=.csv", ticker)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	/*r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}*/
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
		date := time.Date(MustParseDecimal(matches[1]),
			time.Month(MustParseDecimal(matches[2])),
			MustParseDecimal(matches[3]),
			0, 0, 0, 0, time.UTC)
		sql := fmt.Sprintf("INSERT INTO price values('%s', %d, %f, %f, %f, %f, %d, %f)",
			ticker,
			date.Unix(),
			MustParseFloat(line[1]),
			MustParseFloat(line[2]),
			MustParseFloat(line[3]),
			MustParseFloat(line[4]),
			MustParseDecimal(line[5]),
			MustParseFloat(line[6]));
		db.MustUpdate(sql)
	}
	db.MustUpdate("COMMIT TRANSACTION");
	log.Print("End")
	return nil
}

func ReadC() (error) {
	file, err := os.Open("C.bak") // For read access.
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	r, err := reader.ReadAll()
	if err != nil {
		return err
	}
	log.Print("Start, DB=", db)
	db.MustUpdate("BEGIN TRANSACTION");
	for _, line := range r {
		matches := dateRe.FindStringSubmatch(line[0])
		if (matches == nil) {
			continue;
		}
		date := time.Date(MustParseDecimal(matches[1]),
			time.Month(MustParseDecimal(matches[2])),
			MustParseDecimal(matches[3]),
			0, 0, 0, 0, time.UTC)
		sql := fmt.Sprintf("INSERT INTO price values('C', %d, %f, %f, %f, %f, %d, %f)",
			date.Unix(),
			MustParseFloat(line[1]),
			MustParseFloat(line[2]),
			MustParseFloat(line[3]),
			MustParseFloat(line[4]),
			MustParseDecimal(line[5]),
			MustParseFloat(line[6]));
		db.MustUpdate(sql)
	}
	db.MustUpdate("COMMIT TRANSACTION");
	log.Print("End")
	return nil
}

func (db Db) MustUpdate(sql string) {
	st, err := db.db.Prepare(sql)
	PanicOnError(err, "failed to prepare SQL: ", sql)
	st.Step()
	PanicOnError(st.Finalize(), "blah")
}

func (db Db) MustPrepare(sql string) (*sqlite3.Statement) {
	st, err := db.db.Prepare(sql)
	PanicOnError(err, "select " + sql)
	return st
}

func (db Db) MustRunQuery(sql string, cb func(val... interface{})) {
	st := db.MustPrepare(sql)
	_, err := st.All(func(st *sqlite3.Statement, val... interface{}) {
		cb(val...)
	})
	PanicOnError(err, "Failed to run query "+ sql)
}

func (db Db) TableExists(table string) (bool) {
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

func (db Db) GetDateRange(ticker string) (minDate time.Time, maxDate time.Time) {
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
	return minDate, maxDate
}

func CreateDb() (Db) {
	dbName := "./quotes.db"
	db, err := sqlite3.Open(dbName, sqlite3.O_CREATE | sqlite3.O_READWRITE);
	if err != nil {
		log.Panic("Failed to create db: ", dbName, ": ", err)
	}

	var d Db;
	d.db = db
	if !d.TableExists("dividend") {
		d.MustUpdate("CREATE TABLE dividend (ticker VARCHAR(10), date INTEGER, dividend REAL)");
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

func MustParseDecimal(s string) (int) {
	value, err := strconv.ParseInt(s, 10, 0)
	PanicOnError(err, "Failed to parse ", s, " as a decimal string")
	return int(value)
}

func MustParseFloat(s string) (float64) {
	value, err := strconv.ParseFloat(s, 64)
	PanicOnError(err, "Failed to parse ", s, " as a float string")
	return value
}

func Search(ticker string, year int, month int) {
	date := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC);
	date2 := time.Date(year, time.Month(month + 1), 0, 0, 0, 0, 0, time.UTC);

	db.MustRunQuery(fmt.Sprintf(
		"SELECT * from price WHERE ticker = '%s' AND date >= %d AND date <= %d ORDER BY date LIMIT 1",
		ticker, date.Unix(), date2.Unix()),
		func(val... interface{}) {
		log.Print("Got value: ", val)
	})
}

func main() {
	sqlite3.Initialize()
	db = CreateDb()
	defer ShutdownDb();
	log.Print("Created DB=", db)
/*
	err := Download("GOOG")
	if err != nil {
		log.Panic("Failed to read: ", err)
	}
	err = Download("C")
	if err != nil {
		log.Panic("Failed to read: ", err)
	}
*/
	min, max := db.GetDateRange("FOOBAR")
	log.Print("FOO range: [", min, "..", max, "]")
	if min.IsZero() {
		log.Print("FOO ZERO")
	}
	min, max = db.GetDateRange("GOOG")
	log.Print("GOOG range: [", min, "..", max, "]")
	min, max = db.GetDateRange("C")
	log.Print("C range: [", min, "..", max, "]")
/*
	for y := 2000; y < 2012; y++ {
		for m := 1; m < 12; m++ {
			Search("C", y, m);
			Search("GOOG", y, m);
		}
	}
*/
}

