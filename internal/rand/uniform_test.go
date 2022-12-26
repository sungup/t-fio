package rand

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestUniform_Uint64(t *testing.T) {
	expectedRange := uint64(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(16384))
	tested, _ := NewUniform(0, expectedRange)

	for l := 0; l < 100000; l++ {
		assert.Less(t, tested.Uint64(), expectedRange)
	}
}

func TestNewUniform(t *testing.T) {
	for s := int64(0); s < 100; s++ {
		for r := uint64(128); r < 65536; r <<= 1 {
			tested, err := NewUniform(s, r)
			assert.NoError(t, err)
			assert.NotNil(t, tested)
		}
	}
}
