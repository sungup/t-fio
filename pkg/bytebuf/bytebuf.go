package bytebuf

import (
	"github.com/ncw/directio"
	"runtime"
	"sync"
)

type ByteBuf struct {
	data []byte
	pool *sync.Pool
}

func (b *ByteBuf) Close() {
	b.pool.Put(b)
}

func (b *ByteBuf) Buffer() []byte {
	return b.data
}

var (
	pool map[int]*sync.Pool
)

func init() {
	pool = make(map[int]*sync.Pool)
}

func Alloc(size int) (b *ByteBuf) {
	p, ok := pool[size]
	if !ok {
		p = &sync.Pool{}
		p.New = func() interface{} { return &ByteBuf{data: directio.AlignedBlock(size), pool: p} }
		pool[size] = p
	}

	return p.Get().(*ByteBuf)
}

func ForceCleanByteBufPool() {
	pool = make(map[int]*sync.Pool)
	runtime.GC()
}
