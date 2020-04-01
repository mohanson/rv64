package rv64

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 1.1

import (
	"math"
)

// Sign extension is the operation, in computer arithmetic, of increasing the number of bits of a binary number while
// preserving the number's sign (positive/negative) and value. This is done by appending digits to the most significant
// side of the number, following a procedure dependent on the particular signed number representation used.
func SignExtend(v uint64, n uint64) uint64 {
	if v&(1<<n) != 0 {
		return v | (uint64(math.MaxUint64) << n)
	} else {
		return v & (uint64(math.MaxUint64) >> (63 - n))
	}
}
