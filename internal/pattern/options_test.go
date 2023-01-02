package pattern

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/rand"
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
	tested := &Options{}

	tc := map[string]interface{}{
		"type":         "mixed", // unexpected distribution type
		"offset":       DefaultOffset + 1024,
		"page_size":    DefaultPageSize / 8,
		"io_range":     DefaultIORange / 10,
		"distribution": "zip",
		"center":       0.5,
		"seed":         rand.DefaultSeed - 1,
		"theta":        1.2,
		"start_from":   0.25,
	}
	in, _ := json.Marshal(tc)

	// unexpected
	assert.Error(t, json.Unmarshal(in, tested))

	// detail option error test
	tc["type"] = "rand"
	in, _ = json.Marshal(tc)
	assert.Error(t, json.Unmarshal(in, tested))

	// normal zipf random test
	tc["distribution"] = "zipf"
	in, _ = json.Marshal(tc)
	assert.NoError(t, json.Unmarshal(in, tested))
	assert.Equal(t, RandomIO, tested.Type)
	assert.NotEqual(t, DefaultOffset, tested.Offset)
	assert.NotEqual(t, DefaultPageSize, tested.PageSize)
	assert.NotEqual(t, DefaultIORange, tested.IORange)
	assert.IsType(t, &RandOptions{}, tested.detail)
	assert.IsType(t, rand.ZipfDist, tested.detail.(*RandOptions).Distribution)

	// default option check
	assert.NoError(t, json.Unmarshal([]byte("{}"), tested))
	assert.Equal(t, RandomIO, tested.Type)
	assert.Equal(t, DefaultOffset, tested.Offset)
	assert.Equal(t, DefaultPageSize, tested.PageSize)
	assert.Equal(t, DefaultIORange, tested.IORange)
	assert.IsType(t, &RandOptions{}, tested.detail)
	assert.IsType(t, rand.DefaultDist, tested.detail.(*RandOptions).Distribution)
}

func TestOptions_UnmarshalYAML(t *testing.T) {
	tested := &Options{}

	tc := map[string]interface{}{
		"type":         "mixed", // unexpected distribution type
		"offset":       DefaultOffset + 1024,
		"page_size":    DefaultPageSize / 8,
		"io_range":     DefaultIORange / 10,
		"distribution": "zip",
		"center":       0.5,
		"seed":         rand.DefaultSeed - 1,
		"theta":        1.2,
		"start_from":   0.25,
	}
	in, _ := yaml.Marshal(tc)

	// unexpected
	assert.Error(t, yaml.Unmarshal(in, tested))

	// detail option error test
	tc["type"] = "rand"
	in, _ = yaml.Marshal(tc)
	assert.Error(t, yaml.Unmarshal(in, tested))

	// normal zipf random test
	tc["distribution"] = "zipf"
	in, _ = yaml.Marshal(tc)
	assert.NoError(t, yaml.Unmarshal(in, tested))
	assert.Equal(t, RandomIO, tested.Type)
	assert.NotEqual(t, DefaultOffset, tested.Offset)
	assert.NotEqual(t, DefaultPageSize, tested.PageSize)
	assert.NotEqual(t, DefaultIORange, tested.IORange)
	assert.IsType(t, &RandOptions{}, tested.detail)
	assert.IsType(t, rand.ZipfDist, tested.detail.(*RandOptions).Distribution)

	// default option check
	assert.NoError(t, yaml.Unmarshal([]byte("{}"), tested))
	assert.Equal(t, RandomIO, tested.Type)
	assert.Equal(t, DefaultOffset, tested.Offset)
	assert.Equal(t, DefaultPageSize, tested.PageSize)
	assert.Equal(t, DefaultIORange, tested.IORange)
	assert.IsType(t, &RandOptions{}, tested.detail)
	assert.IsType(t, rand.DefaultDist, tested.detail.(*RandOptions).Distribution)
}

func TestOptions_MakeGenerator(t *testing.T) {
	detailRandOpts := new(RandOptions)
	_ = yaml.Unmarshal([]byte("distribution: zipf"), detailRandOpts)

	tested := &Options{
		Type:     SequentialIO,
		Offset:   DefaultOffset,
		PageSize: DefaultPageSize,
		IORange:  DefaultIORange,
		detail:   detailRandOpts,
	}

	// unexpected error test from not allowed randomizer error
	detailRandOpts.Center = -2.0
	generated, err := tested.MakeGenerator()
	assert.Nil(t, generated)
	assert.Error(t, err)
	detailRandOpts.Center = -1.0

	// success cate test
	generated, err = tested.MakeGenerator()
	assert.NotNil(t, generated)
	assert.IsType(t, &Randomizer{}, generated.pattern)
	assert.NoError(t, err)
}
