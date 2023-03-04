package job

import (
	"github.com/ncw/directio"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func TestFillRandomBuf64(t *testing.T) {
	const sz = 1 << 26 // 64MB random buffer
	tested := directio.AlignedBlock(sz)

	samples := uint64(0)
	fillRandomBuf64(tested, uint64(time.Now().UnixNano()))

	for _, v := range tested {
		samples += uint64(v)
	}

	// sum of generate buffer values cannot be zero
	assert.NotZero(t, samples)
}

func TestAllocReadBuffer(t *testing.T) {
	defer func() { bytebuf.ForceCleanByteBufPool() }()
	const mallocRange = 1 << 24
	for i := 0; i < 1000; i++ {
		sz := 100 + rand.Intn(mallocRange-100)
		tested := AllocReadBuffer(sz)
		assert.NotNil(t, tested)

		// sampling 100 numbers and make sum of samples
		sample := int64(0)
		for s := 0; s < 100; s++ {
			sample += int64(tested.Buffer()[rand.Intn(sz)])
		}

		// sum of all samples should be zero because golang always allocate the zeroing buffer
		assert.Zero(t, sample)
	}
}

func TestAllocWriteBuffer(t *testing.T) {
	defer func() { bytebuf.ForceCleanByteBufPool() }()
	const mallocRange = 1 << 24
	for i := 0; i < 1000; i++ {
		sz := 100 + rand.Intn(mallocRange-100)
		tested := AllocWriteBuffer(sz)
		assert.NotNil(t, tested)

		// sampling 100 numbers and make sum of samples
		sample := int64(0)
		for s := 0; s < 100; s++ {
			sample += int64(tested.Buffer()[rand.Intn(sz)])
		}
		// may be sum of all samples cannot be zero because fillRandomBuf64 fill random values in buffer
		assert.NotZero(t, sample)
	}
}

// Benchmark function to check the performance fillRandomBuf64
func BenchmarkFillRandomBuf64(b *testing.B) {
	const sz = 999
	buffer := directio.AlignedBlock(sz)

	b.Run("fillRandomBuf8", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fillRandomBuf8(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)), rand.Uint64())
		}
	})
	b.Run("fillRandomBuf64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fillRandomBuf64(buffer, rand.Uint64())
		}
	})
}

// Benchmark function to check the performance of random value generator to call fillRandomBuf64
func BenchmarkAllocWriteBuffer(b *testing.B) {
	compareAllocWriteBuffer := func(size int) []byte {
		buffer := directio.AlignedBlock(size)
		fillRandomBuf64(buffer, uint64(time.Now().UnixNano()))
		return buffer
	}

	// compare target to check the localRandomizer
	b.Run("AllocWriteBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = AllocWriteBuffer(999)
		}
	})

	// comparing target only use fillRandomBuf using time.Now()
	b.Run("CompareAllocWriteBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = compareAllocWriteBuffer(999)
		}
	})

	// compare target to check the directio.AlignedBlock(size)
	b.Run("AllocReadBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = AllocReadBuffer(999)
		}
	})
}
