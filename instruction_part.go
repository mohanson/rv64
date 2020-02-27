package riscv

import "math"

func InstructionPart(i uint64, f int, e int) uint64 {
	s := i
	s &= uint64(math.MaxUint64) << f
	s &= uint64(math.MaxUint64) >> (63 - e)
	return s >> f
}
