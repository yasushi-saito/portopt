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

	portfolio := NewPortfolio(db,
		NewDateRange(time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		time.Duration(time.Hour * 24 * 30)))
	portfolio.Add("VCADX", 1.0)  // CA interm bond
/*
 	portfolio.Add("VTMGX", 1.0) // tax-managed intl
	portfolio.Add("VPCCX", 1.0) // Primecap core

	portfolio.Add("VVIAX", 1.0) // value index adm
	portfolio.Add("VMVIX", 1.0) // mid-cap value index inv
	portfolio.Add("VMVAX", 1.0) // mid-cap value index adm
	portfolio.Add("VIMSX", 1.0) // mid-cap index inv
	portfolio.Add("VIMAX", 1.0) // mid-cap index adm
	portfolio.Add("VTMSX", 1.0) // T-M smallcap
	portfolio.Add("VGHAX", 1.0) // Healthcare adm
	portfolio.Add("VGHCX", 1.0) // Healthcare inv
	portfolio.Add("VSS", 1.0)  // Ex-us smallcap ETF
	portfolio.Add("VWO", 1.0)  // Emerging market ETF
	portfolio.Add("DLS", 1.0)  // Wisdomtree intl smallcap dividend
	portfolio.Add("VSIAX", 1.0) // Small-cap value index adm
	portfolio.Add("VISVX", 1.0) // small-cap value index inv
*/
	portfolio.Finalize()

	combinedVariance := 0.0
	combinedMean := 0.0
	for _, t1 := range portfolio.Tickers() {
		w1 := portfolio.Weight(t1) / portfolio.TotalWeight()
		stats1, err := db.Stats(t1, portfolio.DateRange())
		if err != nil { panic(err) }

		combinedMean += w1 * stats1.Mean
		for _, t2 := range portfolio.Tickers() {
			w2 := portfolio.Weight(t2) / portfolio.TotalWeight()
			corr, err := db.Correlation(t1, t2)
			if err != nil { panic(err) }

			stats2, err := db.Stats(t2, portfolio.DateRange())
			if err != nil { panic(err) }
			combinedVariance += w1 * w2 * corr * stats1.Stddev * stats2.Stddev
			log.Print("CORR(", t1, ", ", t2, ")=", corr)
		}
	}
	x, _ := db.Stats("VCADX", portfolio.DateRange())
	log.Print("VCADX: ", x.Mean, ", ", x.Stddev)
	log.Print("Combined: ", combinedVariance, " ", math.Sqrt(combinedVariance))
	log.Print("Mean: ", combinedMean)
}

