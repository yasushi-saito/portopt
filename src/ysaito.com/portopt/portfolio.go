//
// Created by Yaz Saito on 06/08/12.
// Copyright (C) 2012 Upthere Inc.
//

package portopt

type portfolioEntry struct {
	weight float64
}

type Portfolio struct {
	db *Database
	securities map[string]*portfolioEntry
	totalWeight float64
	dateRange *dateRange
	finalized bool
}

func NewPortfolio(db *Database, dateRange *dateRange) (p *Portfolio) {
	p = new(Portfolio)
	p.db = db
	p.securities = make(map[string]*portfolioEntry);
	p.dateRange = dateRange
	return p
}

func (p *Portfolio) DateRange() (*dateRange) {
	return p.dateRange
}

func (p *Portfolio) TotalWeight() (float64) {
	return p.totalWeight
}

func (p *Portfolio) Weight(ticker string) (float64) {
	return p.securities[ticker].weight
}

func (p *Portfolio) Tickers() ([]string) {
	var x []string
	for key, _ := range(p.securities) {
		x = append(x, key)
	}
	return x
}

func (p *Portfolio) Add(ticker string, weight float64) (error) {
	if p.finalized { panic("Portfolio already finalized") }
	e := new(portfolioEntry)
	e.weight = weight
	p.securities[ticker] = e
	return nil
}

func (p *Portfolio) Finalize() {
	p.finalized = true
	p.totalWeight = 0.0
	for _, e := range p.securities {
		p.totalWeight += e.weight
	}
}
