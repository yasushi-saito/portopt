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
	label string
	frozen bool
	data []float64
	perPeriodReturn float64
	stddev float64
}

func newStatsAccumulator(label string) *statsAccumulator {
	return &statsAccumulator{label: label, data: make([]float64, 0, 100)}
}

func (s *statsAccumulator) Add(v float64) {
	doAssert(!s.frozen)
	l := len(s.data)
	cc := cap(s.data)
	if l >= cc {
		n := make([]float64, cap(s.data) * 2)
		copy(n, s.data)
		s.data = n
	}

	s.data = s.data[0:l + 1]
	s.data[l] = v
}

func (s *statsAccumulator) freeze() {
	if s.frozen {
		return
	}
	s.frozen = true

	n := float64(len(s.data))
	firstValue := s.data[0]
	lastValue := s.data[len(s.data) - 1]
	if lastValue >= firstValue {
		s.perPeriodReturn = math.Pow((lastValue - firstValue), 1.0 / (n-1)) - 1
	} else {
		s.perPeriodReturn = math.Pow((firstValue - lastValue), 1.0 / (n-1)) - 1
	}

	delta2 := 0.0
	for i, v := range s.data {
		expected := firstValue * math.Pow((1 + s.perPeriodReturn), float64(i))
		delta2 += (v - expected) * (v - expected)
	}
	s.stddev = math.Sqrt(delta2 / n)
/*	log.Print("STAT: ", s.label, " values=", firstValue, ",", lastValue,
		" x=", s.perPeriodReturn)*/

}

func (s *statsAccumulator) NumItems() int {
	return len(s.data)
}

func (s *statsAccumulator) DeltaForPeriod(i int) float64 {
	s.freeze()
	firstValue := s.data[0]

	expected := firstValue * math.Pow((1 + s.perPeriodReturn), float64(i))
	return s.data[i] - expected
}

func (s *statsAccumulator) PerPeriodReturn() float64 {
	s.freeze()
	return s.perPeriodReturn
}

func (s *statsAccumulator) Mean() float64 {
	s.freeze()
	return s.perPeriodReturn
}

func (s *statsAccumulator) StdDev() float64 {
	s.freeze()
	return s.stddev
}
