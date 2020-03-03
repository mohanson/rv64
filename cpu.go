package riscv

type CPU struct {
	Memory   []byte
	pc       uint64
	register [32]uint64
	status   int
	System   System
}

func (c *CPU) BindMemory(m []byte) { c.Memory = m }
func (c *CPU) BindSystem(s System) { c.System = s }
func (c *CPU) GetPC() uint64       { return c.pc }
func (c *CPU) SetPC(i uint64)      { c.pc = i }
func (c *CPU) GetStatus() int      { return c.status }
func (c *CPU) SetStatus(i int)     { c.status = i }

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
