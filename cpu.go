package riscv

type CPU struct {
	pc       uint64
	register [32]uint64
	status   int
	Memory   []byte
	System   System
}

func (c *CPU) SetRegister(i int, u uint64) {
	if i == Rzero {
		return
	}
	c.register[i] = u
}

func (c *CPU) GetRegister(i int) uint64 {
	if i == Rzero {
		return 0x00
	}
	return c.register[i]
}

func (c *CPU) GetPC() uint64   { return c.pc }
func (c *CPU) SetPC(i uint64)  { c.pc = i }
func (c *CPU) GetStatus() int  { return c.status }
func (c *CPU) SetStatus(i int) { c.status = i }
