package job

import (
	"github.com/sungup/t-fio/internal/hash"
	"github.com/sungup/t-fio/pkg/bytebuf"
	"math/rand"
	"time"
	"unsafe"
)

const (
	bucketSz   = uintptr(8)
	szUint64   = unsafe.Sizeof(uint64(0)) // #nosec G103
	bucketMv   = bucketSz * szUint64
	bucketLoop = bucketSz * szUint64
)

var (
	prime = []uint64{1, 2, 3, 5, 7, 11, 13, 17,
		19, 23, 29, 31, 37, 41, 43, 47}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// #nosec G103
func fillRandomBuf8(b, sz uintptr, seed uint64) {
	e := b + ((sz / szUint64) * szUint64)
	rest := sz % szUint64

	for ; b < e; b += szUint64 {
		p := (*uint64)(unsafe.Pointer(b))
		*p = seed
		seed = hash.Hash(seed)
	}

	if rest != 0 {
		for i := uintptr(0); i < rest; i++ {
			d := (*uint8)(unsafe.Pointer(b + i))
			*d = *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&seed)) + i))
		}
	}
}

// #nosec G103
func fillRandomBuf64(buffer []byte, seed uint64) {
	sz := uintptr(len(buffer))

	seeds := [bucketSz]uint64{}

	fill := (sz / szUint64 / bucketSz) * bucketSz * szUint64

	b := uintptr(unsafe.Pointer(&buffer[0]))
	e := b + fill

	for i := range seeds {
		seeds[i] = seed * prime[i]
	}

	for ; b < e; b += bucketMv {
		for i, o := 0, uintptr(0); o < bucketLoop; i, o = i+1, o+szUint64 {
			p := (*uint64)(unsafe.Pointer(b + o))
			*p = seeds[i]
			seeds[i] = hash.Hash(seeds[i])
		}
	}

	fillRandomBuf8(b, sz-fill, seeds[0])
}

func AllocReadBuffer(size int) *bytebuf.ByteBuf {
	return bytebuf.Alloc(size)
}

func AllocWriteBuffer(size int) *bytebuf.ByteBuf {
	buffer := bytebuf.Alloc(size)

	// #nosec G404 ignore weak-random-generator because this is not effects to generate uniform random number
	fillRandomBuf64(buffer.Buffer(), rand.Uint64())

	return buffer
}
