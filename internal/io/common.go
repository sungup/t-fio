package io

import (
	"fmt"
	"github.com/sungup/t-fio/internal/engine"
	"github.com/sungup/t-fio/pkg/sys"
	"gopkg.in/yaml.v3"
	"strings"
)

type DoIO func(p []byte, offset int64, callback engine.Callback) (err error)

type Type func(_ sys.File, _ int64, _ []byte, _ func(bool)) error

func parseType(str string) (Type, error) {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "async_read", "async read":
		return AsyncRead, nil
	case "sync_read", "sync read":
		return SyncRead, nil
	case "write":
		return Write, nil
	default:
		return nil, fmt.Errorf("unexpected IO type name: " + str)
	}
}

func (t *Type) UnmarshalJSON(data []byte) (err error) {
	*t, err = parseType(strings.Trim(string(data), "\""))
	return err
}

func (t *Type) UnmarshalYAML(value *yaml.Node) (err error) {
	*t, err = parseType(value.Value)
	return err
}
