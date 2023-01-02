package types

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
	"strings"
)

var (
	unitRegex  = regexp.MustCompile("^([KMGTP]iB|[KMGTP]B|B|)$")
	byteFormat = regexp.MustCompile("^([0-9]+(|\\.[0-9]*))[[:space:]]*([KMGTP]iB|[KMGTP]B|B|)$")
)

const (
	unitStr = "BKMGTP"
)

func ParseUnit(text string) (size uint64, err error) {
	text = strings.TrimSpace(text)
	if unitRegex.MatchString(text) {
		var ch byte

		if len(text) > 0 {
			ch = text[0]
		} else {
			ch = 'B'
		}

		for i, v := range []byte(unitStr) {
			if v == ch {
				return uint64(1) << (10 * i), nil
			}
		}
	}

	return 0, fmt.Errorf("unexpected size format: \"%v\"", text)
}

type Bytes uint64

func (b *Bytes) Int() int64 {
	return int64(*b)
}

func (b *Bytes) Parse(text string) (err error) {
	var (
		number float64
		unit   uint64
		form   []string
	)

	if form = byteFormat.FindStringSubmatch(strings.TrimSpace(text)); len(form) == 0 {
		return fmt.Errorf("invalid byte string format")
	}

	if number, err = strconv.ParseFloat(form[1], 64); err == nil {
		if unit, err = ParseUnit(form[len(form)-1]); err == nil {
			*b = Bytes(uint64(number * float64(unit)))
		}
	}

	return err
}

func (b *Bytes) String() string {
	v, sz := uint64(*b), uint64(1024)

	if v < sz {
		return strconv.FormatUint(v, 10) + " B"
	} else {
		var i int
		for ; sz <= v; sz <<= 10 {
			i++
		}

		return strconv.FormatFloat(float64(v)/float64(sz>>10), 'f', 3, 64) + " " + string(unitStr[i]) + "iB"
	}
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	return b.Parse(strings.Trim(string(data), "\""))
}

func (b *Bytes) MarshalJSON() ([]byte, error) {
	return []byte("\"" + b.String() + "\""), nil
}

func (b *Bytes) UnmarshalYAML(value *yaml.Node) error {
	return b.Parse(value.Value)
}

// MarshalYAML is different from the MarshalJSON, receiver should be call by value format.
// GoLang's yaml module cannot call the pointer type receiver to marshal data.
func (b Bytes) MarshalYAML() (interface{}, error) {
	return b.String(), nil
}
