package portopt
import "time"

const priceSamplingIntervalSecs = (3600 * 24 * 30)
const priceSamplingInterval = time.Duration(time.Second * priceSamplingIntervalSecs)
var quantizedNow time.Time

func init() {
	now := time.Now().Unix()
	quantized := (now / priceSamplingIntervalSecs - 1) * priceSamplingIntervalSecs
	quantizedNow = time.Unix(quantized, 0)
}