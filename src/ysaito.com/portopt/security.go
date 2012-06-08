//
// Created by Yaz Saito on 06/05/12.
//

package portopt
import "fmt"
import "time"
// import "log"

type Security struct {
	Ticker string
	Weight float64

	db *Database
	minDate time.Time
	maxDate time.Time

	// Date (UNIX time) -> Adjusted closing price, as reported by yahoo.
	priceMap map[int64]float64

};

var allSecurities map[string]*Security;

func init() {
	allSecurities = make(map[string]*Security)
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

func intersectDateRange(
	min1 time.Time, max1 time.Time,
	min2 time.Time, max2 time.Time) (time.Time, time.Time) {
	min := min1
	if min2.After(min) { min = min2 }
	max := max1
	if max2.Before(max) { max = max2 }
	return min, max
}

func NewSecurity(
	db *Database, ticker string, weight float64) *Security {
	s, found := allSecurities[ticker]
	if found { return s }

	s = new(Security)
	s.Ticker = ticker
	s.db = db
	s.Weight = weight
	s.priceMap = make(map[int64]float64)

	for date := quantizedNow; ; date = date.Add(-priceSamplingInterval) {
		price := searchPrice(db, ticker, date, priceSamplingInterval)
		if price >= 0.0 {
			if s.minDate.IsZero() || s.minDate.After(date) {
				s.minDate = date
			}
			if s.maxDate.IsZero() || s.maxDate.Before(date) {
				s.maxDate = date
			}
			s.priceMap[date.Unix()] = price
		} else {
			break;
		}
	}
	allSecurities[ticker] = s
	return s;
}

