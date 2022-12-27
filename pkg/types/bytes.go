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

type unit uint64

const (
	B = unit(1 << (10 * iota))
	KiB
	MiB
	GiB
	TiB
	PiB

	unitStr = "BKMGTP"
)

func (u *unit) Parse(text string) (err error) {
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
				*u = unit(1 << (10 * i))
				return nil
			}
		}
	}

	return fmt.Errorf("unexpected size format: \"%v\"", text)
}

func (u *unit) Size() uint64 {
	return uint64(*u)
}

func (u *unit) String() string {
	if *u == B {
		return "B"
	}

	for i, shift := 1, 10; i < len(unitStr); i, shift = i+1, shift+10 {
		if uint64(*u) == uint64(1)<<shift {
			return string(unitStr[i]) + "iB"
		}
	}

	return "NaN"
}

type Byte struct {
	bytes  uint64
	number float64
	unit   unit
}

func (b *Byte) Parse(text string) (err error) {
	var (
		n float64
		u unit
	)

	form := byteFormat.FindStringSubmatch(text)

	if n, err = strconv.ParseFloat(form[1], 64); err == nil {
		if err = u.Parse(form[len(form)-1]); err == nil {
			b.number = n
			b.unit = u

			b.bytes = uint64(b.number * float64(b.unit.Size()))
		}
	}

	return err
}

func (b *Byte) String() string {
	return strconv.FormatFloat(b.number, 'f', 3, 64) + " " + b.unit.String()
}

func (b *Byte) Uint64() uint64 {
	return b.bytes
}
