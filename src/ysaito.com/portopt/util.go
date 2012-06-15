package portopt
import "math"
import "log"

func doAssert(b bool, message... interface{}) {
	if !b {
		log.Print(message...)
		panic("Assertion failed")
	}
}

type statsAccumulator struct {
	numItems int
	min float64
	max float64
	total float64
	totalSquared float64
}

func (s *statsAccumulator) Add(v float64) {
	if s.numItems == 0 || v < s.min {
		s.min = v
	}
	if s.numItems == 0 || v > s.max {
		s.max = v
	}
	s.numItems++
	s.total += v
	s.totalSquared += v * v
}

func (s *statsAccumulator) NumItems() (int) {
	return s.numItems
}

func (s *statsAccumulator) Min() (float64) {
	return s.min
}

func (s *statsAccumulator) Max() (float64) {
	return s.max
}

func (s *statsAccumulator) Mean() (float64) {
	return s.total / float64(s.numItems)
}

func (s *statsAccumulator) StdDev() (float64) {
	mean := s.Mean()
	return math.Sqrt(s.totalSquared / float64(s.numItems) - mean * mean)
}
