package portopt;

import "log"
import "os"
import "fmt"
import "testing"
import "time"

var pathSeq int = 0;

func newDb(t *testing.T) (db *Database) {
	err := os.MkdirAll("/tmp/portopt_test", 0700)
	if err != nil { t.Fatal(err) }

	pathSeq += 1
	path := fmt.Sprintf("/tmp/portopt_test/db%d.db", pathSeq)
	os.Remove(path)  // ignore error
	return CreateDb(path)
}

func TestCreate(t *testing.T) {
	db := newDb(t)
	err := db.FillFromCsv("C.bak", "C")
	if err != nil { t.Fatal(err) }

	min, max := db.GetDateRange("FOOBAR")
	log.Print("FOO range: [", min, "..", max, "]")
	if !min.IsZero() || !max.IsZero() { t.Error("Range: ", min, max) }

	min, max = db.GetDateRange("C")
	if !min.Equal(time.Date(1977, time.Month(1), 3, 0, 0, 0, 0, time.UTC)) {
		t.Error("Min: ", min.Unix())
	}
	if !max.Equal(time.Date(2012, time.Month(6), 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("Max: ", max)
	}
}

func TestSecurity(t *testing.T) {
	db := newDb(t)
	err := db.FillFromCsv("C.bak", "C")
	if err != nil { t.Fatal(err) }

	minDate := time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(1980, time.Month(12), 31, 0, 0, 0, 0, time.UTC)

	c := NewSecurity(db, "C", minDate, maxDate, time.Hour * 24 * 30)
	log.Print("C=", c, c.Mean(), c.Stddev())
	log.Print("COV=", Covariance(c, c))

	corr := Correlation(c, c)
	if corr < 0.99 || corr >= 1.01 { t.Fatal(corr) }
}

/*
func TestDownload(t *testing.T) {
	db := newDb(t)
	err := db.FillFromYahoo("GOOG")
	if err != nil { panic(err) }

	min, max := db.GetDateRange("GOOG")
	log.Print("GOOG range: [", min, "..", max, "]")
	// goog := NewSecurity(db, "GOOG", min, max, time.Hour * 24 * 30)
}
*/
