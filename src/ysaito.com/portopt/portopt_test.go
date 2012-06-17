package portopt;

import "log"
import "os"
import "fmt"
import "testing"
import "time"
import "github.com/yasushi-saito/fifo_queue"
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

func TestCorrelation(t *testing.T) {
	err := os.MkdirAll("/tmp/portopt_test", 0700)
	if err != nil { t.Fatal(err) }

	path := "/tmp/portopt_test/corr.db"
	db := CreateDb(path)

	corr, err := db.Correlation("C", "C")
	fmt.Print("CORR1=", corr, err, "\n")
	corr, err = db.Correlation("C", "GOOG")
	log.Print("CORR2=", corr, err, "\n")

	corr, err = db.Correlation("VBMFX", "VGTSX")
	log.Print("CORR3=", corr, err, "\n")
}

func TestEff(t *testing.T) {
	err := os.MkdirAll("/tmp/portopt_test", 0700)
	if err != nil { t.Fatal(err) }

	path := "/tmp/portopt_test/eff.db"
	db := CreateDb(path)

	dateRange := NewDateRange(time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		time.Duration(time.Hour * 24 * 90))
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
/*	portfolio := NewPortfolio(db, dateRange, map[string]float64{
		"^GSPC": 1.0,  // S&P 500 index
		"VBMFX" : 1.0,  // Vanguard total bond market index
		"VGTSX" : 1.0,  // Vanguard total intl index
	})*/
/*
	portfolio := NewPortfolio(db, dateRange, map[string]float64{
		"^GSPC": 1.0,  // S&P 500 index
		"VFSTX" : 1.0,  // Vanguard short-term investment grade
		"VGTSX" : 1.0,  // Vanguard total intl index
	})*/

	frontier := newFrontier()
	fifo := fifo_queue.NewQueue()
	fifo.PushBack(portfolio)

	for fifo.Len() > 0 {
		p := fifo.PopFront().(*Portfolio)
		stats := p.Stats()
		maxTries := 20
		if stats.perPeriodReturn >= frontier.MaxX() {
			// Try many times to find a better return
			maxTries = 100
			// fmt.Print(fifo.Len(), " New: Mean: ", mean, " Stddev: ", stddev, "\n")
		} else {
			// fmt.Print(fifo.Len(), " Ins: Mean: ", mean, " Stddev: ", stddev, "\n")
		}

		for i := 0; i < maxTries; i++ {
			newP := p.RandomMutate()
			stats := newP.Stats()
			maxX := frontier.MaxX()
			inserted := frontier.Insert(stats.perPeriodReturn, stats.stddev, newP)
			if inserted {
				fifo.PushBack(newP)
				if stats.perPeriodReturn > maxX {
					// Found a portfolio with the
					// best return so far. We'll
					// start searching from newP
					// with a large maxTries
					// later, so shortcut the
					// search from p now.
					break
				}
			}
		}
	}

	fmt.Print(frontier.String())
}

