package measure

import "time"

func LatencyMeasureStart() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}
