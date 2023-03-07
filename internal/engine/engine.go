package engine

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

type Callback func(n int, err error)
type DoIO func(p []byte, offset int64, callback Callback) (err error)

type IOType int

const (
	Unsupported = IOType(iota)
	Read
	Write
)

func (t *IOType) Parse(text string) (err error) {
	switch strings.TrimSpace(strings.ToLower(text)) {
	case "read", "rd":
		*t = Read
	case "write", "wr":
		*t = Write
	default:
		*t = Unsupported
		err = fmt.Errorf("unsupported IO type")
	}

	return
}

func (t *IOType) String() string {
	switch *t {
	case Read:
		return "read"
	case Write:
		return "write"
	default:
		return "unsupported"
	}
}

func (t *IOType) UnmarshalJSON(data []byte) error {
	return t.Parse(strings.Trim(string(data), "\""))
}

func (t *IOType) UnmarshalYAML(value *yaml.Node) error {
	return t.Parse(value.Value)
}

type Type int

const (
	SyncEngine = Type(iota)
	AsyncEngine
	IOURingEngine
)

func (t *Type) Parse(text string) (err error) {
	switch strings.TrimSpace(strings.ToLower(text)) {
	case "sync":
		*t = SyncEngine
	case "async":
		*t = AsyncEngine
	case "iouring":
		*t = IOURingEngine
	default:
		*t = SyncEngine
		err = fmt.Errorf("unsupported engine type")
	}

	return
}

func (t *Type) String() string {
	switch *t {
	case SyncEngine:
		return "sync"
	case AsyncEngine:
		return "async"
	case IOURingEngine:
		return "iouring"
	default:
		return "unsupported"
	}
}

func (t *Type) UnmarshalJSON(data []byte) error {
	return t.Parse(strings.Trim(string(data), "\""))
}

func (t *Type) UnmarshalYAML(value *yaml.Node) error {
	return t.Parse(value.Value)
}

type Engine interface {
	ReadAt(p []byte, offset int64, callback Callback) (err error)
	WriteAt(p []byte, offset int64, callback Callback) (err error)
	GetIOFunc(ioType IOType) (io DoIO, err error)
	io.Closer
}

func getIOFunc(engine Engine, ioType IOType) (io DoIO, err error) {
	switch ioType {
	case Read:
		io = engine.ReadAt
	case Write:
		io = engine.WriteAt
	default:
		err = fmt.Errorf("unsupported IO type")
	}

	return
}
