package io

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Type func(_ *os.File, _ int64, _ []byte, _ func(bool)) error

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
