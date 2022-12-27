package rand

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestUniform_Uint64(t *testing.T) {
	expectedRange := uint64(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(16384))
	tested, _ := (&UniformOptions{}).MakeRandomizer(0, expectedRange, 0)

	for l := 0; l < 100000; l++ {
		assert.Less(t, tested.Uint64(), expectedRange)
	}
}

func TestUniformOptions_MakeDistributor(t *testing.T) {
	for s := int64(0); s < 100; s++ {
		for r := uint64(128); r < 65536; r <<= 1 {
			tested, err := (&UniformOptions{}).MakeRandomizer(s, r, 0)
			assert.NoError(t, err)
			assert.NotNil(t, tested)
		}
	}
}
