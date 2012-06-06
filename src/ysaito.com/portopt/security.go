//
// Created by Yaz Saito on 06/05/12.
//

package portopt
import "fmt"
import "log"
import "math"
import "time"

type Security struct {
	Ticker string
	db *Database
	mean float64
	stddev float64

	// Date (UNIX time) -> Adjusted closing price, as reported by yahoo.
	priceMap map[int64]float64
	covarianceCache map[string]float64
};

func searchPrice(db *Database, ticker string, date time.Time) (float64) {
	price := -1.0
	db.MustRunQuery(fmt.Sprintf(
		"SELECT adjclose from price WHERE ticker = '%s' AND date >= %d ORDER BY date LIMIT 1",
		ticker, date.Unix()),
		func(val... interface{}) {
		price = val[0].(float64)
		log.Print("Got value: ", price)
	})
	return price
}

func Correlation (s1 *Security, s2 *Security) (float64) {
	return Covariance(s1, s2) / (s1.stddev * s2.stddev)
}

func Covariance (s1 *Security, s2 *Security) (float64) {
	cov, found := s1.covarianceCache[s2.Ticker]
	if found {return cov}

	var diffTotal float64 = 0.0
	for date, value1 := range s1.priceMap {
		value2 := s2.priceMap[date]
		diffTotal += (value1 - s1.mean) * (value2 - s2.mean)
	}
	cov = diffTotal / float64(len(s1.priceMap))
	s1.covarianceCache[s2.Ticker] = cov
	s2.covarianceCache[s1.Ticker] = cov
	return cov
}

func (s Security) Mean() (float64) {
	return s.mean
}

func (s Security) Stddev() (float64) {
	return s.stddev
}

func NewSecurity(
	db *Database, ticker string,
	minDate time.Time,
	maxDate time.Time,
	interval time.Duration) *Security {
	s := new(Security)
	s.Ticker = ticker
	s.db = db
	s.priceMap = make(map[int64]float64)

	for date := minDate; date == maxDate || date.Before(maxDate); date = date.Add(interval) {
		price := searchPrice(db, ticker, date)
		if price < 0.0 { panic(price) }
		s.priceMap[date.Unix()] = price
	}

	// Compute the mean and stddev
	n := 0
	total := 0.0;
	total2 := 0.0;
	for _, value := range s.priceMap {
		total += value
		total2 += value * value;
		n++
	}
	s.mean = total / float64(n);
	s.stddev = math.Sqrt(total2 / float64(n) - (s.mean * s.mean))
	s.covarianceCache = make(map[string]float64)
	return s;
}

