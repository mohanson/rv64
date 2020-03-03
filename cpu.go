package riscv

type CPU struct {
	csr      [4096]uint64
	memory   *Memory
	pc       uint64
	register [32]uint64
	status   int
	system   System
}

func (c *CPU) GetCSR(i int) uint64         { return c.csr[i] }
func (c *CPU) SetCSR(i int, u uint64)      { c.csr[i] = u }
func (c *CPU) GetMemory() *Memory          { return c.memory }
func (c *CPU) SetMemory(m *Memory)         { c.memory = m }
func (c *CPU) GetPC() uint64               { return c.pc }
func (c *CPU) SetPC(i uint64)              { c.pc = i }
func (c *CPU) GetStatus() int              { return c.status }
func (c *CPU) SetStatus(i int)             { c.status = i }
func (c *CPU) GetSystem() System           { return c.system }
func (c *CPU) SetSystem(s System)          { c.system = s }
func (c *CPU) SetRegister(i int, u uint64) { c.register[i] = Cond(i == Rzero, 0x00, u) }
func (c *CPU) GetRegister(i int) uint64    { return Cond(i == Rzero, 0x00, c.register[i]) }

func Cond(b bool, y uint64, f uint64) uint64 {
	if b {
		return y
	}
	return f
}
