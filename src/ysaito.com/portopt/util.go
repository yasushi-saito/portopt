package portopt
import "log"
import "math"

func doAssert(b bool, message... interface{}) {
	if !b {
		log.Print(message...)
		panic("Assertion failed")
	}
}

type statsAccumulator struct {
	label string
	frozen bool     // true once any stats accessor is called
	data []float64  // list of adjusted price quotes
	perPeriodReturn float64
	arithmeticMean float64
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

	lastValue := -1.0
	total := 0.0
	total2 := 0.0
	for i, value := range s.data {
		if i > 0 {
			delta := (value - lastValue) / lastValue
			total += delta
			total2 += delta * delta
		}
		lastValue = value
	}
	s.perPeriodReturn = total / n
	s.stddev = math.Sqrt(total2 / n - s.perPeriodReturn * s.perPeriodReturn)
}

func (s *statsAccumulator) NumItems() int {
	return len(s.data)
}

func (s *statsAccumulator) DeltaForPeriod(i int) float64 {
	s.freeze()
	if i == 0 {
		return 0.0
	}
	value := s.data[i]
	lastValue := s.data[i - 1]
	return (value - lastValue) / lastValue
}

func (s *statsAccumulator) PerPeriodReturn() float64 {
	s.freeze()
	return s.perPeriodReturn
}

func (s *statsAccumulator) ArithmeticMean() float64 {
	s.freeze()
	return s.perPeriodReturn
}

func (s *statsAccumulator) StdDev() float64 {
	s.freeze()
	return s.stddev
}
