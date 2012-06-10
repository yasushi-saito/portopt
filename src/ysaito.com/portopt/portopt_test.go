package portopt;

import "log"
import "os"
import "fmt"
import "testing"
import "time"
import "math"

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

	r := db.GetDateRange("FOOBAR")
	log.Print("FOO range: [", r.Start(), "..", r.End(), "]")
	if !r.Empty()  { t.Error("Range: ", r.Start(), r.End()) }
}

func TestSecurity(t *testing.T) {
	db := newDb(t)
	err := db.FillFromCsv("C.bak", "C")
	if err != nil { t.Fatal(err) }

	corr, err := db.Correlation("C", "C")
	if err != nil { t.Fail() }
	if corr < 0.99 || corr >= 1.01 { t.Fatal(corr) }
}

func TestDateRange(t *testing.T) {
	d := NewDateRange(
		time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		time.Date(1980, time.Month(12), 31, 0, 0, 0, 0, time.UTC),
		time.Duration(time.Hour * 24 * 30))

	for i := d.Begin(); !i.Done(); i.Next() {
		log.Print("ITER: ", i.Time())
	}
}

func TestEff(t *testing.T) {
	err := os.MkdirAll("/tmp/portopt_test", 0700)
	if err != nil { t.Fatal(err) }

	path := "/tmp/portopt_test/eff.db"
	db := CreateDb(path)

	dateRange := NewDateRange(time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		time.Duration(time.Hour * 24 * 30))
	portfolio := NewPortfolio(db, dateRange, map[string]float64{
		"VCADX": 1.0,  // CA interm bond
 		"VTMGX": 1.0, // tax-managed intl
		"VPCCX": 1.0, // Primecap core

		"VVIAX": 1.0, // value index adm
		"VMVIX": 1.0, // mid-cap value index inv
		"VMVAX": 1.0, // mid-cap value index adm
		"VIMSX": 1.0, // mid-cap index inv
		"VIMAX": 1.0, // mid-cap index adm
		"VTMSX": 1.0, // T-M smallcap
		"VGHAX": 1.0, // Healthcare adm
		"VGHCX": 1.0, // Healthcare inv
		"VSS": 1.0,  // Ex-us smallcap ETF
		"DLS": 1.0,  // Wisdomtree intl smallcap dividend
		"VWO": 1.0,  // Emerging market ETF
		"VSIAX": 1.0, // Small-cap value index adm
		"VISVX": 1.0, // small-cap value index inv
	})

	combinedVariance := 0.0
	combinedMean := 0.0
	for _, e1 := range portfolio.List() {
		w1 := e1.weight / portfolio.TotalWeight()
		stats1, err := db.Stats(e1.ticker, portfolio.DateRange())
		if err != nil { panic(err) }

		combinedMean += w1 * stats1.Mean
		for _, e2 := range portfolio.List() {
			w2 := e2.weight / portfolio.TotalWeight()
			stats2, err := db.Stats(e2.ticker, portfolio.DateRange())
			if err != nil { panic(err) }

			corr, err := db.Correlation(e1.ticker, e2.ticker)
			if err != nil { panic(err) }

			combinedVariance += w1 * w2 * corr * stats1.Stddev * stats2.Stddev
			log.Print("CORR(", e1.ticker, ", ", e2.ticker, ")=", corr)
		}
	}
	x, _ := db.Stats("VCADX", portfolio.DateRange())
	log.Print("VCADX: ", x.Mean, ", ", x.Stddev)
	log.Print("Combined: ", combinedVariance, " ", math.Sqrt(combinedVariance))
	log.Print("Mean: ", combinedMean)
}

