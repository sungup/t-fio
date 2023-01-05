package job

import (
	"github.com/sungup/t-fio/internal/io"
	"github.com/sungup/t-fio/internal/pattern"
	"gopkg.in/yaml.v3"
	"time"
)

type Options struct {
	pattern.Options

	Weight            float64       `json:"weight" yaml:"weight"`                         // values between 0~1
	IOType            io.Type       `json:"io_type" yaml:"io_type"`                       // transaction's
	IOSize            int           `json:"io_size" yaml:"io_size"`                       // unit io size
	Delay             time.Duration `json:"delay" yaml:"delay"`                           // stating delay
	TransactionLength int           `json:"transaction_length" yaml:"transaction_length"` // transaction length
}

func (o *Options) UnmarshalJSON(data []byte) error {
	return nil
}

func (o *Options) UnmarshalYAML(value *yaml.Node) error {
	return nil
}
