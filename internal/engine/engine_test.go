package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type tcEngine struct{}

func (t *tcEngine) ReadAt(_ []byte, _ int64, _ Callback) error  { return fmt.Errorf("read at") }
func (t *tcEngine) WriteAt(_ []byte, _ int64, _ Callback) error { return fmt.Errorf("write at") }
func (t *tcEngine) GetIOFunc(_ IOType) (DoIO, error)            { return nil, nil }
func (t *tcEngine) Close() error                                { return nil }

func Test_getIOFunc(t *testing.T) {
	var (
		input = &tcEngine{}

		generated DoIO
		err       error
	)

	// engine panic
	assert.Panics(t, func() { _, _ = getIOFunc(nil, Read) })
	assert.Panics(t, func() { _, _ = getIOFunc(nil, Write) })

	assert.NotPanics(t, func() {
		// tcEngine.ReadAt always return error with string "read at" to identify function
		generated, err = getIOFunc(input, Read)
		assert.NotNil(t, generated)
		assert.NoError(t, err)
		assert.EqualError(t, generated(nil, 0, nil), "read at")

		// tcEngine.WriteAt always return error with string "write at" to identify function
		generated, err = getIOFunc(input, Write)
		assert.NotNil(t, generated)
		assert.NoError(t, err)
		assert.EqualError(t, generated(nil, 0, nil), "write at")

		generated, err = getIOFunc(input, UnsupportedType)
		assert.Nil(t, generated)
		assert.EqualError(t, err, "unsupported IO type")

		generated, err = getIOFunc(nil, UnsupportedType)
		assert.Nil(t, generated)
		assert.EqualError(t, err, "unsupported IO type")
	})
}
