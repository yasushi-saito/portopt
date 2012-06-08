package portopt

type TickerPair struct {
	ticker1 string
	ticker2 string
}

var correlationCache map[TickerPair]float64

func init() {
	correlationCache = make(map[TickerPair]float64)
}

func Correlation (s1 *Security, s2 *Security) (float64) {
	p := TickerPair{ ticker1: s1.Ticker, ticker2: s2.Ticker }
	if p.ticker1 == p.ticker2 { return 1.0 }
	if p.ticker1 > p.ticker2 {
		p.ticker1, p.ticker2 = p.ticker2, p.ticker1
	}

	corr, found := correlationCache[p]
	if found { return corr }

	minDate, maxDate := intersectDateRange(
		s1.minDate, s1.maxDate,
		s2.minDate, s2.maxDate)

	stats1 := new(statsAccumulator);
	stats2 := new(statsAccumulator);

	for d := minDate; d.Equal(maxDate) || d.Before(maxDate);  d = d.Add(priceSamplingInterval) {
		stats1.Add(s1.priceMap[d.Unix()])
		stats2.Add(s2.priceMap[d.Unix()])
	}

	var diffTotal float64 = 0.0
	for d := minDate; d.Equal(maxDate) || d.Before(maxDate);  d = d.Add(priceSamplingInterval) {
		value1 := s1.priceMap[d.Unix()]
		value2 := s2.priceMap[d.Unix()]
		diffTotal += (value1 - stats1.Mean()) * (value2 - stats2.Mean())
	}

	corr = diffTotal / float64(stats1.NumItems()) / stats1.StdDev() / stats2.StdDev()
	correlationCache[p] = corr
	return corr
}


