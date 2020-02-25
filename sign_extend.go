package riscv

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 1.1
// The base integer instruction sets use a two’s-complement representation for signed integer values.

import (
	"math"
)

func SignExtend(v uint64, n int) uint64 {
	if v&(1<<n) != 0 {
		return v | (uint64(math.MaxUint64) << n)
	}
	return v
}
