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
};

func PanicOnError(err error, params... interface{}) {
	if err != nil {
		pp := params;
		log.Fatal(append(pp, ": ", err));
		log.Fatal(err)
	}
}

func mustParseDecimal(s string) (int) {
	value, err := strconv.ParseInt(s, 10, 0)
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

func (db Database) FillFromYahoo(ticker string) (error) {
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
	log.Print("End")
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

func (db Database) GetDateRange(ticker string) (minDate time.Time, maxDate time.Time) {
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


func CreateDb(path string) (*Database) {
	sqlite3.Initialize()
	db, err := sqlite3.Open(path, sqlite3.O_CREATE | sqlite3.O_READWRITE);
	if err != nil {
		log.Panic("Failed to create db: ", path, ": ", err)
	}

	var d *Database = new(Database);
	d.Path = path
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
