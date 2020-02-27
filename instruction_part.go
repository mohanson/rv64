package riscv

import "math"

func InstructionPart(data []byte, f int, e int) uint64 {
	var s uint64 = 0
	for i := len(data) - 1; i >= 0; i-- {
		s += uint64(data[i]) << (8 * i)
	}
	s &= uint64(math.MaxUint64) << f
	s &= uint64(math.MaxUint64) >> (63 - e)
	return s >> f
}
