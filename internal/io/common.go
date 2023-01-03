package io

import "os"

type Type func(_ *os.File, _ int64, _ []byte, _ func(bool)) error
