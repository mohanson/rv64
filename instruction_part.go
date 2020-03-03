package riscv

import "math"

func InstructionPart(i uint64, f uint64, e uint64) uint64 {
	s := i
	s &= uint64(math.MaxUint64) << f
	s &= uint64(math.MaxUint64) >> (63 - e)
	return s >> f
}
