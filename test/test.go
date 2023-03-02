package test

import (
	"math/rand"
	"os"
	"path"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ncw/directio"
)

const (
	BufferSz = directio.BlockSize
)

var (
	Buffer []byte
)

func init() {
	Buffer = directio.AlignedBlock(directio.BlockSize)

	FillBuffer(Buffer, time.Now().Add(time.Hour*-1).UnixNano())
}

func FillBuffer(buffer []byte, seed int64) {
	rnd := rand.New(rand.NewSource(seed)) // #nosec: G404

	for i := range buffer {
		buffer[i] = uint8(rnd.Intn(255))
	}
}

func fillFile(fp *os.File, size int64) (err error) {
	for offset := int64(0); offset < size; offset += directio.BlockSize {
		if _, err = fp.WriteAt(Buffer, offset); err != nil {
			break
		}
	}

	return err
}

func OpenTCFile(filename string, size int64) (fp *os.File, closer func(), err error) {
	tcFilePath := path.Join(os.TempDir(), filename+"-"+time.Now().Format("20060102150405"))

	if fp, err = os.Create(path.Clean(tcFilePath)); err == nil {
		closer = func() {
			_ = fp.Close()
			_ = os.Remove(tcFilePath)
		}

		err = fillFile(fp, size)
	}

	if err != nil {
		if closer != nil {
			closer()
		}

		fp, closer = nil, nil
	}

	return fp, closer, err
}

func ProjectDir() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Dir(path.Dir(filename))
}

type AtomicCounter int32

func (c *AtomicCounter) Add(delta int32) { atomic.AddInt32((*int32)(c), delta) }
func (c *AtomicCounter) Done()           { c.Add(-1) }
func (c *AtomicCounter) Len() int32      { return atomic.LoadInt32((*int32)(c)) }
