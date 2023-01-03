package measure

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLatencyMeasureStart(t *testing.T) {
	minDuration := time.Millisecond * 500

	tested := LatencyMeasureStart()
	time.Sleep(minDuration)
	assert.Greater(t, tested(), minDuration)
}
