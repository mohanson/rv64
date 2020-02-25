package riscv

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 1.1
// The base integer instruction sets use a two’s-complement representation for signed integer values.

import "math"

// signExtend extends given bit (counting from 0) in v. This function allows
// converting signed numbers from an N-bit (N<32) representation to 64-bit
// representation
func signExtend(v uint64, bit int) uint64 {
	b := signBits[bit]
	if v&b.signBit != 0 {
		return v | b.ones
	}
	return v
}

var signBits = [32]struct {
	signBit uint64
	ones    uint64
}{}

func init() {
	b := uint64(1)
	ones := uint64(math.MaxUint64)
	for i := 0; i < len(signBits); i++ {
		signBits[i].signBit = b
		signBits[i].ones = ones
		b <<= 1
		ones <<= 1
	}
}
