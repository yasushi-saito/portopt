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
	if !max.IsZero()  { t.Error("Range: ", min, max) }

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

	c := NewSecurity(db, "C", 1.0)
	corr := Correlation(c, c)
	if corr < 0.99 || corr >= 1.01 { t.Fatal(corr) }
}


func TestEff(t *testing.T) {
	err := os.MkdirAll("/tmp/portopt_test", 0700)
	if err != nil { t.Fatal(err) }

	path := "/tmp/portopt_test/eff.db"
	db := CreateDb(path)

	portfolio := map[string]float64 {
		"VCADX" : 1.0,  // CA interm bond
		"VTMGX" : 1.0,  // tax-managed intel
		"VPCCX" : 1.0,  // Primecap core
		"VVIAX" : 1.0,  // value index adm
		"VMVIX" : 1.0,  // mid-cap value index inv
		"VMVAX" : 1.0,  // mid-cap value index adm
		"VIMSX" : 1.0,  // mid-cap index inv
		"VIMAX" : 1.0,  // mid-cap index adm
		"VTMSX" : 1.0,  // T-M smallcap
		"VGHAX" : 1.0, // Healthcare adm
		"VGHCX" : 1.0, // Healthcare inv
		"VSS" : 1.0,  // Ex-us smallcap ETF
		"VWO" : 1.0,  // Emerging market ETF
		"DLS" : 1.0,  // Wisdomtree intl smallcap dividend
		"VSIAX" : 1.0, // Small-cap value index adm
		"VISVX" : 1.0, // small-cap value index inv
	}


	for ticker, _ := range portfolio {
		err = db.FillFromYahooIfNecessary(ticker)
		if err != nil { panic(err) }
	}

	tickers := make(map[string]*Security)

	for ticker, weight := range portfolio {
		tickers[ticker] = NewSecurity(db, ticker, weight)
	}

	for t1, s1 := range tickers {
		for t2, s2 := range tickers {
			log.Print("CORR(", t1, ", ", t2, ")=",
				Correlation(s1, s2))
		}
	}
}

