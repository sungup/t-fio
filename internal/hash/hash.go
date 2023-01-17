package hash

const (
	// GoldenRatio64 is a hash key value used in the FIO
	GoldenRatio64 = uint64(0x61C8864680B583EB)
)

func Hash(v uint64) uint64 {
	return v * GoldenRatio64
}
