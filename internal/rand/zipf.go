package rand

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/rand"
)

const (
	errMsgZipfTheta = "theta value is not acceptable to create zipf random: %v"

	DefaultZipfTheta = 1.2
)

// Zipf is a random value generator following the Zipf's law
type Zipf struct {
	core
	zipf *rand.Zipf
}

func (z *Zipf) Uint64() uint64 {
	return z.core.hash(z.zipf.Uint64())
}

// ZipfOptions is an additional option container to generate the Zipf randomizer
type ZipfOptions struct {
	Theta float64
}

func (o *ZipfOptions) UnmarshalJSON(data []byte) (err error) {
	buffer := struct{ Theta float64 }{DefaultZipfTheta}
	if err = json.Unmarshal(data, &buffer); err == nil {
		if 1.0 < buffer.Theta {
			o.Theta = buffer.Theta
		} else {
			err = fmt.Errorf(errMsgZipfTheta, buffer.Theta)
		}
	}

	return err
}

func (o *ZipfOptions) UnmarshalYAML(value *yaml.Node) (err error) {
	buffer := struct{ Theta float64 }{DefaultZipfTheta}
	if err = value.Decode(&buffer); err == nil {
		if 1.0 < buffer.Theta {
			o.Theta = buffer.Theta
		} else {
			err = fmt.Errorf(errMsgZipfTheta, buffer.Theta)
		}
	}

	return err
}

func (o *ZipfOptions) MakeRandomizer(seed int64, nRange uint64, center float64) (Rand, error) {
	z := &Zipf{}

	if err := z.init(seed, nRange, center); err != nil {
		return nil, err
	}

	if z.zipf = rand.NewZipf(z.rand, o.Theta, 1.0, nRange); z.zipf == nil {
		return nil, fmt.Errorf(errMsgZipfTheta, o.Theta)
	}

	return z, nil
}
