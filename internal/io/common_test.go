package io

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
	"unsafe"
)

var (
	tcValidTypeStr map[uint64][]string
)

func init() {
	aReadFunc := AsyncRead
	sReadFunc := SyncRead
	writeFunc := Write

	aReadPtr := *(*uint64)(unsafe.Pointer(&aReadFunc))
	sReadPtr := *(*uint64)(unsafe.Pointer(&sReadFunc))
	writePtr := *(*uint64)(unsafe.Pointer(&writeFunc))

	tcValidTypeStr = map[uint64][]string{
		aReadPtr: {"async_read", "async read"},
		sReadPtr: {"sync_read", "sync read"},
		writePtr: {"write"},
	}

}

func tcPaddingStr(in string) []string {
	return []string{
		in, " " + in, in + " ", " " + in + " ",
		strings.ToUpper(in), " " + strings.ToUpper(in), strings.ToUpper(in) + " ", " " + strings.ToUpper(in) + " ",
		strings.ToTitle(in), " " + strings.ToTitle(in), strings.ToTitle(in) + " ", " " + strings.ToTitle(in) + " ",
	}
}

func TestParseType(t *testing.T) {
	var (
		generated Type
		err       error
	)

	for _, in := range tcPaddingStr("") {
		generated, err = parseType(in)
		assert.Nil(t, generated)
		assert.EqualError(t, err, "unexpected IO type name: "+in)
	}

	for expectedPtr, tcItems := range tcValidTypeStr {
		for _, tc := range tcItems {
			for _, in := range tcPaddingStr(tc) {
				generated, err = parseType(in)
				assert.NotNil(t, generated)
				assert.NoError(t, err)

				generatedPtr := *(*uint64)(unsafe.Pointer(&generated))
				assert.Equal(t, expectedPtr, generatedPtr)
			}

			generated, err = parseType(tc[:len(tc)-1])
			assert.Nil(t, generated)
			assert.EqualError(t, err, "unexpected IO type name: "+tc[:len(tc)-1])
		}
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	template := "{\"io_type\": \"%s\"}"
	tested := &struct {
		IOType Type `json:"io_type" yaml:"io_type"`
	}{}

	for _, in := range tcPaddingStr("") {
		assert.EqualError(t, json.Unmarshal([]byte(fmt.Sprintf(template, in)), tested), "unexpected IO type name: "+in)
		assert.Nil(t, tested.IOType)
	}

	for expectedPtr, tcItems := range tcValidTypeStr {
		for _, tc := range tcItems {
			for _, in := range tcPaddingStr(tc) {
				assert.NoError(t, json.Unmarshal([]byte(fmt.Sprintf(template, in)), tested))
				assert.NotNil(t, tested.IOType)

				generatedPtr := *(*uint64)(unsafe.Pointer(&tested.IOType))
				assert.Equal(t, expectedPtr, generatedPtr)
			}

			in := tc[:len(tc)-1]
			assert.EqualError(t, json.Unmarshal([]byte(fmt.Sprintf(template, in)), tested), "unexpected IO type name: "+in)
			assert.Nil(t, tested.IOType)
		}
	}
}

func TestType_UnmarshalYAML(t *testing.T) {
	template := "io_type: \"%s\""
	tested := &struct {
		IOType Type `json:"io_type" yaml:"io_type"`
	}{}

	for _, in := range tcPaddingStr("") {
		assert.EqualError(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, in)), tested), "unexpected IO type name: "+in)
		assert.Nil(t, tested.IOType)
	}

	for expectedPtr, tcItems := range tcValidTypeStr {
		for _, tc := range tcItems {
			for _, in := range tcPaddingStr(tc) {
				assert.NoError(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, in)), tested))
				assert.NotNil(t, tested.IOType)

				generatedPtr := *(*uint64)(unsafe.Pointer(&tested.IOType))
				assert.Equal(t, expectedPtr, generatedPtr)
			}

			in := tc[:len(tc)-1]
			assert.EqualError(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, in)), tested), "unexpected IO type name: "+in)
			assert.Nil(t, tested.IOType)
		}
	}
}
