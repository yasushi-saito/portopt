//
// Created by Yaz Saito on 06/15/12.
//

package portopt
import "testing"
import "log"
func TestStat_Basic(t *testing.T) {
	s := newStatsAccumulator("test")
	s.Add(1.0)
	s.Add(2.0)
	s.Add(4.0)
	s.Add(5.0)
	s.Add(16.0)
	log.Print("S=", s.PerPeriodReturn(), s.StdDev())
}
