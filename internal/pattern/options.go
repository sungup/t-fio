package pattern

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type detailOptions interface {
	MakeIOPattern(nRange int64) (IOPattern, error)
}

type Type int

const (
	Unsupported = Type(iota)
	RandomIO
	SequentialIO
)

var (
	ioTypeInfo = map[Type]struct {
		makeDetailOptions func() detailOptions
		longName          string
		shortName         string
	}{
		RandomIO: {
			func() detailOptions { return &RandOptions{} },
			"random",
			"rand",
		},
		SequentialIO: {
			func() detailOptions { return &SeqOptions{} },
			"sequential",
			"seq",
		},
	}
)

func (t *Type) makeDetailOptions() detailOptions {
	if info, ok := ioTypeInfo[*t]; ok {
		return info.makeDetailOptions()
	} else {
		return nil
	}
}

func (t *Type) Parse(text string) error {
	text = strings.TrimSpace(strings.ToLower(text))

	for k, info := range ioTypeInfo {
		if info.shortName == text || info.longName == text {
			*t = k
			return nil
		}
	}

	*t = Unsupported
	return fmt.Errorf("unsupported IO pattern: %v", text)
}

func (t *Type) String() string {
	if info, ok := ioTypeInfo[*t]; ok {
		return info.longName
	} else {
		return "unsupported"
	}
}

func (t *Type) UnmarshalJSON(data []byte) error {
	return t.Parse(strings.Trim(string(data), "\""))
}

func (t *Type) UnmarshalYAML(value *yaml.Node) error {
	return t.Parse(value.Value)
}

type Options struct {
	Type     Type  // IO Pattern Type
	Offset   int64 // IO block start address
	PageSize int64 // unit io size
	IORange  int64 // IO range

	detail detailOptions
}

const (
	DefaultOffset   = int64(0)
	DefaultPageSize = int64(4096) // Default IO page size is 4KB IO
)

func (o *Options) UnmarshalJSON(data []byte) (err error) {
	buffer := struct {
		Type Type

		// TODO need to change Byte format data structure
		Offset   int64
		PageSize int64
		IORange  int64
	}{
		Type:     RandomIO,
		Offset:   DefaultOffset,
		PageSize: DefaultPageSize,
	}

	if err = json.Unmarshal(data, &buffer); err == nil {
		detail := buffer.Type.makeDetailOptions()

		if err = json.Unmarshal(data, detail); err == nil {
			o.Type = buffer.Type
			o.Offset = buffer.Offset
			o.PageSize = buffer.PageSize
			o.IORange = buffer.IORange
		}
	}

	return err
}

func (o *Options) UnmarshalYAML(value *yaml.Node) (err error) {
	buffer := struct {
		Type Type

		// TODO need to change Byte format data structure
		Offset   int64
		PageSize int64
		IORange  int64
	}{
		Type:     RandomIO,
		Offset:   DefaultOffset,
		PageSize: DefaultPageSize,
	}

	if err = value.Decode(&buffer); err == nil {
		detail := buffer.Type.makeDetailOptions()

		if err = value.Decode(detail); err == nil {
			o.Type = buffer.Type
			o.Offset = buffer.Offset
			o.PageSize = buffer.PageSize
			o.IORange = buffer.IORange
		}
	}

	return err
}

func (o *Options) MakeGenerator() (generator *Generator, err error) {
	var pattern IOPattern

	if pattern, err = o.detail.MakeIOPattern(o.IORange / o.PageSize); err == nil {
		generator = &Generator{
			pattern:    pattern,
			pageOffset: o.Offset,
			pageSz:     o.PageSize,
		}
	}

	return generator, err
}
