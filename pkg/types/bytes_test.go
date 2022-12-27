package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_Parse(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestUnit_Size(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestUnit_String(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestByte_Parse(t *testing.T) {
	tested := Byte{}

	assert.NoError(t, tested.Parse("1"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1."))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1B"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1B"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.0B"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1 B"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1 B"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1GB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1GB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1 GB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1 GB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1GiB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1GiB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1 GiB"))
	t.Log(tested.String())
	assert.NoError(t, tested.Parse("1.1 GiB"))
	t.Log(tested.String())

	assert.Fail(t, "not yet implemented")

}

func TestByte_String(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

func TestByte_Uint64(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}
