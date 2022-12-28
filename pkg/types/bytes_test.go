package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

var (
	tcNumbers = []float64{25.5, 10.0, 15.0, 1.0, 0.5, 0.25}

	tcValidUnitFormat     map[string]uint64
	tcInvalidUnitFormat   []string
	tcValidNumberFormat   map[string]uint64
	tcInvalidNumberFormat []string
)

func makePadding(text string) []string {
	return []string{text, text + " ", " " + text, " " + text + " "}
}

func init() {
	tcValidUnitFormat = map[string]uint64{"": 1, "B": 1}
	tcInvalidUnitFormat = []string{"invalid", "b"}
	tcValidNumberFormat = make(map[string]uint64)
	tcInvalidNumberFormat = []string{"", "B"}

	for _, n := range tcNumbers {
		strNumber := fmt.Sprintf("%f", n)
		for _, key := range makePadding(strNumber) {
			tcValidNumberFormat[key] = uint64(n)

			for _, tail := range makePadding("B") {
				tcValidNumberFormat[key+tail] = uint64(n)
			}
		}

		tcInvalidNumberFormat = append(tcInvalidNumberFormat, strings.ToLower(strNumber+"B"))
	}

	for i, v := range unitStr[1:] {
		unit := string(v)
		size := uint64(1) << ((i + 1) * 10)

		for _, tail := range []string{"iB", "B"} {
			tcValidUnitFormat[unit+tail] = size
			tcValidUnitFormat[unit+tail+" "] = size
			tcValidUnitFormat[" "+unit+tail] = size
			tcValidUnitFormat[" "+unit+tail+" "] = size
		}
		tcInvalidUnitFormat = append(tcInvalidUnitFormat, strings.ToLower(string(v)+"iB"), strings.ToLower(string(v)+"B"))

		for _, n := range tcNumbers {
			strNumber := fmt.Sprintf("%f", n)
			for _, number := range makePadding(strNumber) {
				for _, form := range []string{"iB", "B"} {
					for _, tail := range makePadding(unit + form) {
						tcValidNumberFormat[number+tail] = uint64(n * float64(size))
						tcInvalidNumberFormat = append(tcInvalidUnitFormat, strings.ToLower(number+tail))
					}
				}
			}
		}
	}
}

func TestParseUnit(t *testing.T) {
	for in, expected := range tcValidUnitFormat {
		generated, err := ParseUnit(in)
		assert.NoError(t, err)
		assert.Equal(t, expected, generated)
	}

	for _, in := range tcInvalidUnitFormat {
		generated, err := ParseUnit(in)
		assert.Error(t, err)
		assert.Zero(t, generated)
	}
}

func TestBytes_Parse(t *testing.T) {
	for in, expected := range tcValidNumberFormat {
		var tested Bytes
		assert.NoError(t, tested.Parse(in), in)
		assert.Equal(t, expected, uint64(tested))
	}

	for _, in := range tcInvalidNumberFormat {
		var tested Bytes
		assert.Error(t, tested.Parse(in))
		assert.Zero(t, uint64(tested))
	}
}

func TestBytes_String(t *testing.T) {
	for sample, value := range tcValidNumberFormat {
		tested := Bytes(value)
		in, err := strconv.ParseFloat(strings.Trim(sample, " iB"+unitStr), 64)
		t.Logf("result: %v | in: %v | err: %v", tested.String(), in, err)
	}
	assert.Fail(t, "not yet implemented")
}
