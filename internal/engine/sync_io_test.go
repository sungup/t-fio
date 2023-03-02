package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/test"
	"testing"
	"time"
)

func TestSyncIO_ReadAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested := &SyncIO{}

	// fail test
	tcCounter.Add(1)
	assert.EqualError(t, tested.ReadAt(nil, 0, tcFailedCB), "invalid argument")
	assert.Zero(t, tcCounter.Len())

	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestSyncIO_ReadAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested = &SyncIO{fp: tcFile}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		var testedBuffer [test.BufferSz]byte

		tcCounter.Add(1)
		assert.NoError(t, tested.ReadAt(testedBuffer[:], tcOffset, tcSuccessCB))
		assert.Zero(t, tcCounter.Len())
		assert.EqualValues(t, test.Buffer, testedBuffer[:])
	}
}

func TestSyncIO_WriteAt(t *testing.T) {
	tcCounter := test.AtomicCounter(0)
	writtenBuffer := make([]byte, test.BufferSz)
	expectedBuffer := make([]byte, test.BufferSz)
	test.FillBuffer(expectedBuffer, time.Now().UnixNano())

	tcFailedCB := makeFailedCallback(t, &tcCounter)
	tcSuccessCB := makeSuccessCallback(t, &tcCounter)

	tested := SyncIO{}

	// fail test
	tcCounter.Add(1)
	assert.Error(t, tested.WriteAt(nil, 0, tcFailedCB))
	assert.Zero(t, tcCounter.Len())

	// success test
	tcFile, tcCloser, err := test.OpenTCFile("TestSyncIO_WriteAt", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested = SyncIO{fp: tcFile}

	for tcOffset := int64(0); tcOffset < tcFileSz; tcOffset += test.BufferSz {
		// check write is success
		tcCounter.Add(1)
		assert.NoError(t, tested.WriteAt(expectedBuffer, tcOffset, tcSuccessCB))
		assert.Zero(t, tcCounter.Len())

		// check written data
		_, err = tcFile.ReadAt(writtenBuffer, tcOffset)
		assert.NoError(t, err)
		assert.Equal(t, expectedBuffer, writtenBuffer)
		assert.NotEqual(t, test.Buffer, writtenBuffer)
	}
}

func TestSyncIO_Close(t *testing.T) {
	tcFile, tcCloser, err := test.OpenTCFile("TestSyncIO_Close", tcFileSz)
	assert.NoError(t, err)
	defer tcCloser()

	tested := SyncIO{fp: tcFile}

	// close file without error
	assert.NoError(t, tested.Close())

	// raise error already closed file
	assert.Error(t, tested.Close())
}
