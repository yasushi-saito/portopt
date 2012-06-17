//
// Created by Yaz Saito on 06/08/12.
// Copyright (C) 2012 Upthere Inc.
//

package portopt
import "log"
import "fmt"
import "time"

const minInterval = time.Duration(time.Hour * 24 * 30) // 30 days
type dateRange struct {
	start time.Time
	end time.Time
	samplingInterval time.Duration
}
type dateRangeIterator struct {
	r *dateRange
	t time.Time
}

func roundInterval(d time.Duration) time.Duration {
	if d <= minInterval {
		return minInterval
	}
	mod := d % minInterval
	if mod < minInterval / 2 {
		return d / minInterval * minInterval
	}
	return (d / minInterval + 1) * minInterval
}

func NewDateRange(
	start time.Time,
	end time.Time,
	desiredInterval time.Duration) (*dateRange) {
	interval := roundInterval(desiredInterval)
	log.Print("INTERVAL: ", desiredInterval, " ", interval)
	intervalSecs := int64(interval / (1000 * 1000 * 1000))
	quantizedStart := (start.Unix() / intervalSecs) * intervalSecs
	quantizedEnd := (end.Unix() / intervalSecs) * intervalSecs

	r := new(dateRange)
	r.start = time.Unix(quantizedStart, 0)
	r.end = time.Unix(quantizedEnd, 0)
	r.samplingInterval = time.Duration(intervalSecs * (1000 * 1000 * 1000))
	return r
}

func (d *dateRange) String() string {
	return fmt.Sprintf("[%v,%v,%v]", d.start, d.end, d.samplingInterval)
}

func (d *dateRange) Empty() bool {
	return d.start.After(d.end)
}

func (d *dateRange) Start() (time.Time) { return d.start }
func (d *dateRange) End() (time.Time) { return d.end }

func (d1 *dateRange) Inside(d2 *dateRange) (bool) {
	if d1.start.Before(d2.start) { return false }
	if d1.end.After(d2.end) { return false }
	return true
}

func (d1 *dateRange) Intersect(d2 *dateRange) (*dateRange) {
	if (d1.samplingInterval != d2.samplingInterval) {
		panic("Trying to intersect ranges with different interval")
	}
	n := new(dateRange)
	n.samplingInterval = d1.samplingInterval
	n.start = d1.start
	if d2.start.After(n.start) { n.start = d2.start }

	n.end = d1.end
	if d2.end.Before(n.end) { n.end = d2.end }
	return n
}

func (d *dateRange) Begin() (*dateRangeIterator) {
	i := new(dateRangeIterator)
	i.r = d
	i.t = d.start
	return i
}

func (i *dateRangeIterator) Time() (time.Time) {
	return i.t
}

func (i *dateRangeIterator) Done() (bool) {
	return i.t.After(i.r.end)
}

func (i *dateRangeIterator) Next() {
	if (i.Done()) { panic("done") }
	i.t = i.t.Add(i.r.samplingInterval)
}

