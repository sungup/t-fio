package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
	"time"
)

func TestAsyncIO_ReadAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested := &AsyncIO{}

	// fail test
	tcCounter.Add(1)
	assert.NoError(t, tested.ReadAt(nil, 0, tcFailedCB))
	assert.NotZero(t, tcCounter.Len())

	// wait until all thread completed
	tcCounter.Wait()

	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestAsyncIO_ReadAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested = &AsyncIO{fp: tcFile}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		var testedBuffer [test.BufferSz]byte

		tcCounter.Add(1)
		assert.NoError(t, tested.ReadAt(testedBuffer[:], tcOffset, tcSuccessCB))
		assert.NotZero(t, tcCounter.Len())
	}

	// wait until all thread completed
	tcCounter.Wait()
}

func TestAsyncIO_WriteAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	writtenBuffer := make([]byte, test.BufferSz)
	expectedBuffer := make([]byte, test.BufferSz)
	test.FillBuffer(expectedBuffer, time.Now().UnixNano())

	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested := AsyncIO{}

	// fail test
	tcCounter.Add(1)
	assert.NoError(t, tested.WriteAt(nil, 0, tcFailedCB))
	assert.NotZero(t, tcCounter.Len())

	// wait until all thread completed
	tcCounter.Wait()

	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestAsyncIO_WriteAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested = AsyncIO{fp: tcFile}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		// check write is success
		tcCounter.Add(1)
		assert.NoError(t, tested.WriteAt(expectedBuffer, tcOffset, tcSuccessCB))
		assert.NotZero(t, tcCounter.Len())
	}

	// wait until all thread completed
	tcCounter.Wait()

	// check written data after all io completed
	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		_, err = tcFile.ReadAt(writtenBuffer, tcOffset)
		assert.NoError(t, err)
		assert.Equal(t, expectedBuffer, writtenBuffer)
		assert.NotEqual(t, test.Buffer, writtenBuffer)
	}
}

func TestAsyncIO_GetIOFunc(t *testing.T) {
	var (
		generated DoIO
		err       error

		tested = AsyncIO{}
	)

	generated, err = tested.GetIOFunc(Read)
	assert.NotNil(t, generated)
	assert.NoError(t, err)

	generated, err = tested.GetIOFunc(Write)
	assert.NotNil(t, generated)
	assert.NoError(t, err)

	generated, err = tested.GetIOFunc(UnsupportedType)
	assert.Nil(t, generated)
	assert.Error(t, err)
}

func TestAsyncIO_Close(t *testing.T) {
	tcFile, tcCloser, err := test.OpenTCFile("TestAsyncIO_Close", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested := AsyncIO{fp: tcFile}

	// close file without error
	assert.NoError(t, tested.Close())

	// raise error already closed file
	assert.Error(t, tested.Close())
}
