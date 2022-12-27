package rand

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"math"
	"testing"
)

func TestDistribution_makeDetailOptions(t *testing.T) {
	// supported option test
	tc := map[Distribution]detailOptions{
		Unsupported:               nil,
		UniformDist:               &UniformOptions{},
		ZipfDist:                  &ZipfOptions{},
		Distribution(math.MaxInt): nil,
	}

	for tested, expectedType := range tc {
		assert.IsType(t, tested.makeDetailOptions(), expectedType)
	}
}

func TestDistribution_Parse(t *testing.T) {
	var tested Distribution

	// fail test
	for _, in := range []string{"unif", "zip", "unsupport"} {
		tested = UniformDist
		assert.Error(t, tested.Parse(in))
		assert.Equal(t, Unsupported, tested)
	}

	// success test
	tc := map[Distribution][]string{
		UniformDist: {"uniform", " UniForm ", " UNIFORM", "uniForm "},
		ZipfDist:    {"zipf", " ZipF ", " ZIPF", "zipF "},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested = Unsupported
			assert.NoError(t, tested.Parse(in))
			assert.Equal(t, expected, tested)
		}
	}
}

func TestDistribution_String(t *testing.T) {
	tc := map[Distribution]string{
		Unsupported:               "unsupported",
		UniformDist:               "uniform",
		ZipfDist:                  "zipf",
		Distribution(math.MaxInt): "unsupported",
	}

	for tested, expected := range tc {
		assert.Equal(t, expected, tested.String())
	}
}

func TestDistribution_UnmarshalJSON(t *testing.T) {
	tested := &struct {
		Type Distribution
	}{}
	template := "{\"type\": \"%s\"}"

	// distribution type fail test
	for _, in := range []string{"unif", "zip", "unsupported"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, json.Unmarshal(buffer, tested), "unsupported random distribution")
	}

	// distribution type success test
	tc := map[string]Distribution{
		"ZipF":    ZipfDist,
		"Uniform": UniformDist,
	}

	for in, expected := range tc {
		tested.Type = Unsupported
		buffer := []byte(fmt.Sprintf(template, in))
		assert.NoError(t, json.Unmarshal(buffer, tested))
		assert.Equal(t, expected, tested.Type)
	}
}

func TestDistribution_UnmarshalYAML(t *testing.T) {
	tested := &struct {
		Distribution Distribution
	}{}
	template := "distribution: %v"

	// distribution type fail test
	for _, in := range []string{"unif", "zip", "unsupported"} {
		buffer := []byte(fmt.Sprintf(template, in))
		assert.EqualError(t, yaml.Unmarshal(buffer, tested), "unsupported random distribution")
	}

	// distribution type success test
	tc := map[Distribution][]string{
		UniformDist: {"uniform", "\"uniform\""},
		ZipfDist:    {"zipf", "\"zipf\""},
	}

	for expected, tcItems := range tc {
		for _, in := range tcItems {
			tested.Distribution = Unsupported
			buffer := []byte(fmt.Sprintf(template, in))
			assert.NoError(t, yaml.Unmarshal(buffer, tested))
			assert.Equal(t, expected, tested.Distribution)

		}
	}
}

func TestOptions_UnmarshalJSON(t *testing.T) {
	tested := &Options{}

	tc := map[string]interface{}{
		"distribution": "zip", // unexpected distribution type
		"center":       0.5,
		"seed":         DefaultSeed - 1,
		"disable_hash": !DefaultDisableHash,
		"theta":        0.5, // invalid theta value for zipf
	}
	in, _ := json.Marshal(tc)

	// unexpected distribution error test
	assert.Error(t, json.Unmarshal(in, tested))

	// detail option error test
	tc["distribution"] = "zipf"
	in, _ = json.Marshal(tc)
	assert.Error(t, json.Unmarshal(in, tested))

	// normal zipf unmarshal test
	tc["theta"] = 1.2
	in, _ = json.Marshal(tc)
	assert.NoError(t, json.Unmarshal(in, tested))
	assert.NotEqual(t, DefaultCenter, tested.Center)
	assert.NotEqual(t, DefaultSeed, tested.Seed)
	assert.NotEqual(t, DefaultDisableHash, tested.DisableHash)
	assert.IsType(t, &ZipfOptions{}, tested.detail)

	// default option check
	assert.NoError(t, json.Unmarshal([]byte("{}"), tested))
	assert.Equal(t, DefaultDist, tested.Distribution)
	assert.Equal(t, DefaultCenter, tested.Center)
	assert.Equal(t, DefaultSeed, tested.Seed)
	assert.Equal(t, DefaultDisableHash, tested.DisableHash)
	assert.IsType(t, &UniformOptions{}, tested.detail)
}

func TestOptions_UnmarshalYAML(t *testing.T) {
	tested := &Options{}

	tc := map[string]interface{}{
		"distribution": "zip", // unexpected distribution type
		"center":       0.5,
		"seed":         DefaultSeed - 1,
		"disable_hash": !DefaultDisableHash,
		"theta":        0.5, // invalid theta value for zipf
	}
	in, _ := yaml.Marshal(tc)

	// unexpected distribution error test
	assert.Error(t, yaml.Unmarshal(in, tested))

	// detail option error test
	tc["distribution"] = "zipf"
	in, _ = yaml.Marshal(tc)
	assert.Error(t, yaml.Unmarshal(in, tested))

	// normal zipf unmarshal test
	tc["theta"] = 1.2
	in, _ = yaml.Marshal(tc)
	assert.NoError(t, yaml.Unmarshal(in, tested))
	assert.NotEqual(t, DefaultCenter, tested.Center)
	assert.NotEqual(t, DefaultSeed, tested.Seed)
	assert.NotEqual(t, DefaultDisableHash, tested.DisableHash)
	assert.IsType(t, &ZipfOptions{}, tested.detail)

	// default option check
	assert.NoError(t, yaml.Unmarshal([]byte("{}"), tested))
	assert.Equal(t, DefaultDist, tested.Distribution)
	assert.Equal(t, DefaultCenter, tested.Center)
	assert.Equal(t, DefaultSeed, tested.Seed)
	assert.Equal(t, DefaultDisableHash, tested.DisableHash)
	assert.IsType(t, &UniformOptions{}, tested.detail)
}
