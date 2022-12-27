package rand

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/rand"
	"strings"
)

type detailOptions interface {
	MakeRandomizer(seed int64, nRange uint64, center float64) (Rand, error)
}

type Distribution int

const (
	Unsupported = Distribution(iota)
	UniformDist
	ZipfDist
)

var (
	distributionInfo = map[Distribution]struct {
		makeDetailOptions func() detailOptions
		name              string
	}{
		UniformDist: {
			func() detailOptions { return &UniformOptions{} },
			"uniform",
		},
		ZipfDist: {
			func() detailOptions { return &ZipfOptions{} },
			"zipf",
		},
	}
)

func (t *Distribution) makeDetailOptions() detailOptions {
	if info, ok := distributionInfo[*t]; ok {
		return info.makeDetailOptions()
	} else {
		return nil
	}
}

func (t *Distribution) Parse(text string) error {
	text = strings.TrimSpace(strings.ToLower(text))

	for k, info := range distributionInfo {
		if info.name == text {
			*t = k
			return nil
		}
	}

	*t = Unsupported
	return fmt.Errorf("unsupported random distribution")
}

func (t *Distribution) String() string {
	if info, ok := distributionInfo[*t]; ok {
		return info.name
	} else {
		return "unsupported"
	}
}

func (t *Distribution) UnmarshalJSON(data []byte) error {
	return t.Parse(strings.Trim(string(data), "\""))
}

func (t *Distribution) UnmarshalYAML(value *yaml.Node) error {
	return t.Parse(value.Value)
}

type Options struct {
	Distribution Distribution
	Center       float64
	Seed         int64
	DisableHash  bool

	detail detailOptions
}

const (
	DefaultDist        = UniformDist
	DefaultCenter      = float64(-1)
	DefaultDisableHash = false
)

var (
	DefaultSeed = rand.Int63()
)

func (o *Options) UnmarshalJSON(data []byte) (err error) {
	buffer := struct {
		Distribution Distribution
		Center       float64
		Seed         int64
		DisableHash  bool `json:"disable_hash" yaml:"disable_hash"`
	}{
		DefaultDist,
		DefaultCenter,
		DefaultSeed,
		DefaultDisableHash,
	}

	if err = json.Unmarshal(data, &buffer); err == nil {
		detail := buffer.Distribution.makeDetailOptions()

		if err = json.Unmarshal(data, detail); err == nil {
			o.Distribution = buffer.Distribution
			o.Center = buffer.Center
			o.Seed = buffer.Seed
			o.DisableHash = buffer.DisableHash
			o.detail = detail
		}
	}

	return err
}

func (o *Options) UnmarshalYAML(value *yaml.Node) (err error) {
	buffer := struct {
		Distribution Distribution
		Center       float64
		Seed         int64
		DisableHash  bool `json:"disable_hash" yaml:"disable_hash"`
	}{
		DefaultDist,
		DefaultCenter,
		DefaultSeed,
		DefaultDisableHash,
	}

	if err = value.Decode(&buffer); err == nil {
		detail := buffer.Distribution.makeDetailOptions()

		if err = value.Decode(detail); err == nil {
			o.Distribution = buffer.Distribution
			o.Center = buffer.Center
			o.Seed = buffer.Seed
			o.DisableHash = buffer.DisableHash
			o.detail = detail
		}
	}

	return err
}

func (o *Options) MakeRandomizer(nRange uint64) (Rand, error) {
	return o.detail.MakeRandomizer(o.Seed, nRange, o.Center)
}
