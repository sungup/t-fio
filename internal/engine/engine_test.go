package engine

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"math"
	"testing"
)

type tcEngine struct{}

func (t *tcEngine) ReadAt(_ []byte, _ int64, _ Callback) error  { return fmt.Errorf("read at") }
func (t *tcEngine) WriteAt(_ []byte, _ int64, _ Callback) error { return fmt.Errorf("write at") }
func (t *tcEngine) GetIOFunc(_ IOType) (DoIO, error)            { return nil, nil }
func (t *tcEngine) Close() error                                { return nil }

func TestIOType_Parse(t *testing.T) {
	var tested IOType

	// fail test
	for _, in := range []string{"rx", "wx", "unsupport"} {
		tested = Write
		assert.Error(t, tested.Parse(in))
		assert.Equal(t, Unsupported, tested)
	}

	// success test
	tc := map[IOType][]string{
		Read:  {"read", " Read ", " READ", "reaD "},
		Write: {"write", " Write ", " WRITE", "writE "},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested = Unsupported
			assert.NoError(t, tested.Parse(in))
			assert.Equal(t, expected, tested)
		}
	}
}

func TestIOType_String(t *testing.T) {
	tc := map[IOType]string{
		Read:                "read",
		Write:               "write",
		IOType(math.MaxInt): "unsupported",
	}

	for tested, expected := range tc {
		assert.Equal(t, expected, tested.String())
	}
}

func TestIOType_UnmarshalJSON(t *testing.T) {
	tested := &struct {
		Type IOType
	}{}
	template := "{\"type\": \"%s\"}"

	// fail test
	for _, in := range []string{"unsupported", "down"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, json.Unmarshal(buffer, tested), "unsupported IO type")
	}

	// success test
	tc := map[string]IOType{
		"read":  Read,
		"write": Write,
	}

	for in, expected := range tc {
		tested.Type = Unsupported
		buffer := []byte(fmt.Sprintf(template, in))
		assert.NoError(t, json.Unmarshal(buffer, tested))
		assert.Equal(t, expected, tested.Type)
	}
}

func TestIOType_UnmarshalYAML(t *testing.T) {
	tested := &struct {
		Type IOType
	}{}
	template := "type: %v"

	// fail test
	for _, in := range []string{"unsupported", "down"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, yaml.Unmarshal(buffer, tested), "unsupported IO type")
	}

	// success test
	tc := map[IOType][]string{
		Read:  {"read", "\"read\""},
		Write: {"write", "\"write\""},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested.Type = Unsupported
			buffer := []byte(fmt.Sprintf(template, in))
			assert.NoError(t, yaml.Unmarshal(buffer, tested))
			assert.Equal(t, expected, tested.Type)
		}
	}
}

func TestType_Parse(t *testing.T) {
	var tested Type

	// fail test
	for _, in := range []string{"dir", "file", "unsupport"} {
		tested = Type(0xFF)
		assert.Error(t, tested.Parse(in))
		assert.Equal(t, SyncEngine, tested)
	}

	// success test
	tc := map[Type][]string{
		SyncEngine:    {"sync", " Sync ", " SYNC", "synC "},
		AsyncEngine:   {"async", " ASync ", " ASYNC", "aSync "},
		IOURingEngine: {"iouring", " IOURing ", "IOURING", "ioUring "},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested = Type(0xFF)
			assert.NoError(t, tested.Parse(in))
			assert.Equal(t, expected, tested)
		}
	}
}

func TestType_String(t *testing.T) {
	tc := map[Type]string{
		SyncEngine:        "sync",
		AsyncEngine:       "async",
		IOURingEngine:     "iouring",
		Type(math.MaxInt): "unsupported",
	}

	for tested, expected := range tc {
		assert.Equal(t, expected, tested.String())
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	tested := &struct {
		Type Type
	}{}
	template := "{\"type\": \"%s\"}"

	// fail test
	for _, in := range []string{"unsupported", "down"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, json.Unmarshal(buffer, tested), "unsupported engine type")
	}

	// success test
	tc := map[string]Type{
		"sync":    SyncEngine,
		"async":   AsyncEngine,
		"iouring": IOURingEngine,
	}

	for in, expected := range tc {
		tested.Type = Type(0xFF)
		buffer := []byte(fmt.Sprintf(template, in))
		assert.NoError(t, json.Unmarshal(buffer, tested))
		assert.Equal(t, expected, tested.Type)
	}
}

func TestType_UnamrshalYAML(t *testing.T) {
	tested := &struct {
		Type Type
	}{}
	template := "type: %v"

	// fail test
	for _, in := range []string{"unsupported", "down"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, yaml.Unmarshal(buffer, tested), "unsupported engine type")
	}

	// success test
	tc := map[Type][]string{
		SyncEngine:    {"sync", "\"sync\""},
		AsyncEngine:   {"async", "\"async\""},
		IOURingEngine: {"iouring", "\"iouring\""},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested.Type = Type(0xFF)
			buffer := []byte(fmt.Sprintf(template, in))
			assert.NoError(t, yaml.Unmarshal(buffer, tested))
			assert.Equal(t, expected, tested.Type)
		}
	}
}

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
