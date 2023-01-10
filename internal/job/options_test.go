package job

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/pattern"
	"github.com/sungup/t-fio/internal/rand"
	"gopkg.in/yaml.v3"
	"testing"
)

var (
	tcOptions = map[string]interface{}{
		"type":         "mixed",
		"offset":       pattern.DefaultOffset + 1024,
		"page_size":    pattern.DefaultPageSize / 8,
		"io_range":     pattern.DefaultIORange / 10,
		"distribution": "zipf",
		"center":       0.5,
		"seed":         rand.DefaultSeed - 1,
		"theta":        1.2,
		"start_from":   0.25,
	}
)

func TestOptions_UnmarshalYAML(t *testing.T) {
	assert.Fail(t, "not yet implemented")

	tested := &Options{}
	assert.NoError(t, yaml.Unmarshal([]byte(""), tested))
	//buffer, _ := yaml.Marshal(tested)
	//fmt.Println("\n" + string(buffer))

}

func TestOptions_UnmarshalJSON(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}
