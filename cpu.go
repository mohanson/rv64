package riscv

type CPU struct {
	Register [32]uint64
	Memory   []byte
	System   System
	PC       uint64
	Stop     bool
}

func (c *CPU) SetRegister(i int, u uint64) {
	if i == Rzero {
		return
	}
	c.Register[i] = u
}

func (c *CPU) GetRegister(i int) uint64 {
	if i == Rzero {
		return 0x00
	}
	return c.Register[i]
}
