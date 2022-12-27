package pattern

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"math"
	"testing"
)

func TestType_makeDetailOptions(t *testing.T) {
	tc := map[Type]detailOptions{
		Unsupported:       nil,
		RandomIO:          &RandOptions{},
		SequentialIO:      &SeqOptions{},
		Type(math.MaxInt): nil,
	}

	for tested, expected := range tc {
		assert.IsType(t, tested.makeDetailOptions(), expected)
	}
}

func TestType_Parse(t *testing.T) {
	var tested Type

	// fail test
	for _, in := range []string{"rnd", "sequential_io", "unsupported"} {
		tested = RandomIO
		assert.Error(t, tested.Parse(in))
		assert.Equal(t, Unsupported, tested)
	}

	// success test
	tc := map[Type][]string{
		RandomIO:     {"Random", "rand", " RANDOM ", "RAND ", " ranDom"},
		SequentialIO: {"Sequential", "seq", " SEQUENTIAL ", "SEQ ", " sequential"},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested = Unsupported
			assert.NoError(t, tested.Parse(in))
			assert.Equal(t, expected, tested)
		}
	}
}

func TestType_String(t *testing.T) {
	tc := map[Type]string{
		Unsupported:       "unsupported",
		RandomIO:          "random",
		SequentialIO:      "sequential",
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

	// distribution type fail test
	for _, in := range []string{"rnd", "sq", "unsupported"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, json.Unmarshal(buffer, tested), "unsupported IO pattern: "+in)
	}

	// type unmarshal succes stest
	tc := map[string]Type{
		"seq":    SequentialIO,
		"random": RandomIO,
	}

	for in, expected := range tc {
		tested.Type = Unsupported
		buffer := []byte(fmt.Sprintf(template, in))
		assert.NoError(t, json.Unmarshal(buffer, tested))
		assert.Equal(t, expected, tested.Type)
	}
}

func TestType_UnmarshalYAML(t *testing.T) {
	tested := &struct {
		Type Type
	}{}
	template := "type: %s"

	// distribution type fail test
	for _, in := range []string{"rnd", "sq", "unsupported"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, yaml.Unmarshal(buffer, tested), "unsupported IO pattern: "+in)
	}

	// type unmarshal succes stest
	tc := map[string]Type{
		"\"seq\"": SequentialIO,
		"random":  RandomIO,
	}

	for in, expected := range tc {
		tested.Type = Unsupported
		buffer := []byte(fmt.Sprintf(template, in))
		assert.NoError(t, yaml.Unmarshal(buffer, tested))
		assert.Equal(t, expected, tested.Type)
	}
}

func TestOptions_UnmarshalJSON(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestOptions_UnmarshalYAML(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestOptions_MakeGenerator(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}
