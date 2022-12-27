package pattern

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/rand"
	"gopkg.in/yaml.v3"
	"testing"
)

type tcRandomizer struct {
	i uint64
}

func (t *tcRandomizer) Uint64() uint64 {
	v := t.i
	t.i++
	return v
}

func (t *tcRandomizer) EnableHash(_ bool) {}

func TestRandomizer_PageNo(t *testing.T) {
	tested := Randomizer{rnd: &tcRandomizer{}}
	loop := int64(1024)

	for expected := int64(0); expected < loop; expected++ {
		assert.Equal(t, expected, tested.PageNo())
	}
}

func TestRandOptions_MakeIOPattern(t *testing.T) {
	tested := RandOptions{}

	// test unexpected settings
	_ = yaml.Unmarshal([]byte("distribution: zipf"), &tested)
	tested.Center = -2.0
	generated, err := tested.MakeIOPattern(1024)
	assert.Error(t, err)
	assert.Nil(t, generated)

	// test zipf default settings
	tested.Center = -1.0
	generated, err = tested.MakeIOPattern(1024)
	assert.NoError(t, err)
	assert.IsType(t, &rand.Zipf{}, generated.(*Randomizer).rnd)

	// test default value
	_ = yaml.Unmarshal([]byte("distribution: uniform"), &tested)
	generated, err = tested.MakeIOPattern(1024)
	assert.NoError(t, err)
	assert.IsType(t, &rand.Uniform{}, generated.(*Randomizer).rnd)
}
