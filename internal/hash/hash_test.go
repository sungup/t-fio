package hash

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	sample := map[uint64]int{}
	v, _ := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	data := v.Uint64()

	for i := 0; i < 10000000; i++ {
		data = Hash(data)
		sample[data]++
	}

	for _, selected := range sample {
		assert.Less(t, selected, 2)
	}
}
