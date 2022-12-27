package rand

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestZipf_Uint64(t *testing.T) {
	nRange := uint64(1000)
	buckets := make([]uint64, nRange)
	tested, _ := (&ZipfOptions{Theta: 1.2}).MakeRandomizer(0, nRange, 0)
	loop := 1000000

	tested.EnableHash(false)

	for i := 0; i < loop; i++ {
		testedV := tested.Uint64()
		assert.Less(t, testedV, nRange)
		buckets[testedV]++
	}

	last := buckets[0]
	for i, testedC := range buckets[1:] {
		// TODO is this the best way???
		// Top 3% values should be decrementing order because there are a lot the large appearance
		// change differentials.
		if i < 30 {
			assert.Less(t, testedC, last)
		} else {
			break
		}

		last = testedC
	}
}

func TestZipfOptions_UnmarshalJSON(t *testing.T) {
	tested := &ZipfOptions{}
	template := "{\"theta\": %v}"

	// type mismatch error
	assert.Error(t, json.Unmarshal([]byte(fmt.Sprintf(template, "\"string\"")), tested))
	assert.Zero(t, tested.Theta)

	// invalid value range
	assert.Error(t, json.Unmarshal([]byte(fmt.Sprintf(template, 0.5)), tested))
	assert.Zero(t, tested.Theta)

	// valid change
	expected := 1.5
	assert.NoError(t, json.Unmarshal([]byte(fmt.Sprintf(template, expected)), tested))
	assert.Equal(t, expected, tested.Theta)

	// default option check
	assert.NoError(t, json.Unmarshal([]byte("{}"), tested))
	assert.Equal(t, DefaultZipfTheta, tested.Theta)
}

func TestZipfOptions_UnmarshalYAML(t *testing.T) {
	tested := &ZipfOptions{}
	template := "theta: %v"

	// type mismatch error
	assert.Error(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, "\"string\"")), tested))
	assert.Zero(t, tested.Theta)

	// invalid value range
	assert.Error(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, 0.5)), tested))
	assert.Zero(t, tested.Theta)

	// valid change
	expected := 1.5
	assert.NoError(t, yaml.Unmarshal([]byte(fmt.Sprintf(template, expected)), tested))
	assert.Equal(t, expected, tested.Theta)

	// default option check
	assert.NoError(t, yaml.Unmarshal([]byte("empty: string"), tested))
	assert.Equal(t, DefaultZipfTheta, tested.Theta)
}

func TestZipfOptions_MakeRandomizer(t *testing.T) {
	tested := ZipfOptions{}

	// zipf creation fail
	tested.Theta = 0.5
	testedV, err := tested.MakeRandomizer(0, 100, 0.5)
	assert.Nil(t, testedV)
	assert.EqualError(t, err, fmt.Sprintf("theta value is not acceptable to create zipf random: %v", 0.5))

	tested.Theta = 1.2

	// initialization fail
	testedV, err = tested.MakeRandomizer(0, 100, -2.0)
	assert.Nil(t, testedV)
	assert.EqualError(t, err, fmt.Sprintf("unexpected center range: %v", -2.0))

	// successful creation
	testedV, err = tested.MakeRandomizer(0, 100, 0.5)
	assert.NotNil(t, tested)
	assert.NoError(t, err)
}
