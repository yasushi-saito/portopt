//
// Created by Yaz Saito on 06/08/12.
// Copyright (C) 2012 Upthere Inc.
//

package portopt
import "math/rand"

type portfolioEntry struct {
	ticker string
	weight float64
}

type Portfolio struct {
	db *Database
	securities []portfolioEntry
	totalWeight float64
	dateRange *dateRange
}

func NewPortfolio(db *Database,
	dateRange *dateRange,
	securities map[string]float64) (p *Portfolio) {
	p = new(Portfolio)
	p.db = db
	p.securities = make([]portfolioEntry, len(securities))
	p.dateRange = dateRange
	n := 0
	for s, w := range securities {
		p.securities[n].ticker = s
		p.securities[n].weight = w
		p.totalWeight += w
		n++
	}
	return p
}

func (p *Portfolio) RandomMutate() (*Portfolio) {
	n := len(p.securities)
	delta := p.totalWeight * 0.02;
	q := Portfolio{
		db: p.db,
		securities: make([]portfolioEntry, n),
	        totalWeight: 0.0, // filled later
	        dateRange: p.dateRange,
	}
	for i, e := range p.securities {
		q.securities[i] = e
	}

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

