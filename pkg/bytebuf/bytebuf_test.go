package bytebuf

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
)

func TestByteBuf_Close(t *testing.T) {
	// do nothing because there is no way test Put is success
	assert.NotPanics(t, func() {
		p := &sync.Pool{}
		tested := ByteBuf{
			data: make([]byte, 100),
			pool: p,
		}

		tested.Close()
	})
}

func TestByteBuf_Buffer(t *testing.T) {
	defer func() { ForceCleanByteBufPool() }()
	const loop = 100
	for i := 1; i <= loop; i++ {
		tested := Alloc(i)
		generated := tested.Buffer()
		assert.NotNil(t, generated)
		assert.Len(t, generated, i)
	}
}

func TestNewPool(t *testing.T) {
	const loop = 1024
	for expectedSz := 1; expectedSz <= loop; expectedSz++ {
		tested := newPool(expectedSz)
		generated := tested.Get().(*ByteBuf)
		assert.Len(t, generated.data, expectedSz)
		assert.Equal(t, generated.pool, tested)
		tested.Put(generated)
	}
}

func BenchmarkAlloc(b *testing.B) {
	defer func() { ForceCleanByteBufPool() }()
	const loop = 1000000
	const size = 128 << 10

	b.Run("benchmark_alloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < loop; l++ {
				generated := Alloc(size)
				generated.Close()
				generated.data[0] = byte(128)
			}
		}
	})

	b.Run("benchmark_noalloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for l := 0; l < loop; l++ {
				generated := &ByteBuf{
					data: make([]byte, size),
					pool: nil,
				}
				generated.data[0] = byte(128)
			}
		}
	})
}

func TestAlloc(t *testing.T) {
	defer func() { ForceCleanByteBufPool() }()
	const loop = 100
	for i := 1; i <= loop; i++ {
		generated := Alloc(i)
		assert.NotNil(t, generated)
		cnt := 0

		pool.Range(func(k, _ any) bool {
			assert.LessOrEqual(t, k, i)
			cnt++
			return true
		})

		assert.Equal(t, cnt, i)
	}
}

func TestForceCleanByteBufPool(t *testing.T) {
	const loop = 100
	var (
		memStat   = &runtime.MemStats{}
		generated *ByteBuf
		sum       = 0

		before, after uint64
	)

	for i := 1; i <= loop; i++ {
		generated = Alloc(i)
		generated.Close()
		sum += i
	}

	runtime.ReadMemStats(memStat)
	before = memStat.HeapObjects

	ForceCleanByteBufPool()

	runtime.ReadMemStats(memStat)
	after = memStat.HeapObjects

	cnt := 0
	pool.Range(func(_, _ any) bool {
		cnt++
		return true
	})
	assert.Zero(t, cnt)
	assert.Greater(t, before, after)
}
