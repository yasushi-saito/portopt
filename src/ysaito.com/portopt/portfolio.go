//
// Created by Yaz Saito on 06/08/12.
// Copyright (C) 2012 Upthere Inc.
//

package portopt
import "math"
import "math/rand"

type portfolioEntry struct {
	ticker string
	weight float64
}

type PortfolioStats struct {
	perPeriodReturn float64
	stddev float64
}

type Portfolio struct {
	db *Database
	securities []portfolioEntry
	totalWeight float64
	dateRange *dateRange
	cachedStats PortfolioStats
}

func NewPortfolio(db *Database,
	dateRange *dateRange,
	securities map[string]float64) (p *Portfolio) {
	p = new(Portfolio)
	p.db = db
	p.securities = make([]portfolioEntry, len(securities))
	p.dateRange = dateRange
	p.cachedStats.perPeriodReturn = -1.0  // sentinel
	p.cachedStats.stddev = -1.0
	n := 0
	for s, w := range securities {
		p.securities[n].ticker = s
		p.securities[n].weight = w
		p.totalWeight += w
		n++
	}
	return p
}

func (p *Portfolio) Stats() PortfolioStats {
	if p.cachedStats.perPeriodReturn < 0.0 {
		variance := 0.0
		perPeriodReturn := 0.0
		arithMean := 0.0
		db := p.Db()
		for _, e1 := range p.List() {
			w1 := e1.weight / p.TotalWeight()
			stats1, err := db.Stats(e1.ticker, p.DateRange())
			if err != nil { panic(err) }

			perPeriodReturn += w1 * stats1.PerPeriodReturn
			arithMean += w1 * stats1.ArithmeticMean

			for _, e2 := range p.List() {
				corr, err := db.Correlation(e1.ticker, e2.ticker)
				if err != nil { panic(err) }
				w2 := e2.weight / p.TotalWeight()
				stats2, err := db.Stats(e2.ticker, p.DateRange())
				if err != nil { panic(err) }
				variance += w1 * w2 * corr * stats1.Stddev * stats2.Stddev
			}
		}
		var stddev float64
		if (variance <= 0) {
			stddev = 0
		} else {
			stddev = math.Sqrt(variance) / arithMean
		}
		p.cachedStats.perPeriodReturn = perPeriodReturn
		p.cachedStats.stddev = stddev
	}
	return p.cachedStats
}

func (p *Portfolio) RandomMutate() (*Portfolio) {
	n := len(p.securities)
	q := Portfolio{
		db: p.db,
		securities: make([]portfolioEntry, n),
	        totalWeight: 0.0, // filled later
		dateRange: p.dateRange,
	        cachedStats: PortfolioStats{-1.0, -1.0},
	}
	for i, e := range p.securities {
		q.securities[i] = e
	}

	for i := 0; i < 10; i++ {
		delta := p.totalWeight * 0.01;
		tmp := delta
		for true {
			i := rand.Intn(n)
			if q.securities[i].weight > tmp {
				q.securities[i].weight -= tmp
				break
			} else {
				tmp -= q.securities[i].weight
				q.securities[i].weight = 0
			}
		}
		i := rand.Intn(n)
		q.securities[i].weight += delta
	}
	// Recompute the total weight
	for _, e := range q.securities {
		q.totalWeight += e.weight
	}
	return &q
}

func (p *Portfolio) Db() (*Database) { return p.db }

func (p *Portfolio) DateRange() (*dateRange) {
	return p.dateRange
}

func (p *Portfolio) TotalWeight() (float64) {
	return p.totalWeight
}

func (p *Portfolio) Weight(ticker string) (float64) {
	for _, e := range p.securities {
		if (e.ticker == ticker) { return e.weight }
	}
	return 0.0
}

func (p *Portfolio) List() ([]portfolioEntry) {
	return p.securities
}

