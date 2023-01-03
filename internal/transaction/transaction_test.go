package transaction

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/t-fio/internal/io"
	"math"
	"math/big"
	"os"
	"testing"
)

func tcMakeIOList(tcFunc io.Type, jobId int64) []*io.IO {
	ios := make([]*io.IO, 0)

	for offset := int64(0); offset < (4096 << 3); offset += 4096 {
		ios = append(ios, io.NewIO(tcFunc, jobId, offset, make([]byte, 4096)))
	}

	return ios
}

func TestTransaction_ProcessAll(t *testing.T) {
	const errMessage = "TestTransaction_ProcessAll error message"

	vRand, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	tested := Transaction{jobId: vRand.Int64()}
	testedCounter := 0

	tcErrorIO := func(_ *os.File, _ int64, _ []byte, cb func(bool)) error {
		testedCounter++
		cb(false)
		return fmt.Errorf(errMessage)
	}

	tested.ios = tcMakeIOList(tcErrorIO, tested.jobId)
	testedCounter = 0
	assert.EqualError(t, tested.ProcessAll(), errMessage)
	assert.Equal(t, testedCounter, 1)

	assert.Fail(t, "not yet implemented")
}
