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
	pool sync.Map
)

func newPool(size int) (p *sync.Pool) {
	p = &sync.Pool{}
	p.New = func() interface{} { return &ByteBuf{data: directio.AlignedBlock(size), pool: p} }
	return p
}

func Alloc(size int) (b *ByteBuf) {
	v, _ := pool.LoadOrStore(size, newPool(size))

	return v.(*sync.Pool).Get().(*ByteBuf)
}

func ForceCleanByteBufPool() {
	pool.Range(func(k, _ any) bool {
		pool.Delete(k)
		return true
	})
	runtime.GC()
}
