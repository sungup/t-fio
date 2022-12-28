package types

import (
	"fmt"
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
	const max = len(unitStr) - 1
	for i, unit := max, uint64(1)<<max; 0 < i; i, unit = i-1, unit>>10 {
		if unit < uint64(*b) {
			number := float64(*b) / float64(unit)
			return strconv.FormatFloat(number, 'f', 3, 64) + " " + string(unitStr[i-1]) + "iB"
		}
	}

	return "B"
}
