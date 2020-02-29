package riscv

type CPU struct {
	Register [32]uint64
	Memory   []byte
	PC       uint64
}
