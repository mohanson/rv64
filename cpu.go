package rv64

type CPU struct {
	csr       [4096]uint64
	lraddr    uint64
	memory    *Memory
	pc        uint64
	register  [32]uint64
	status    uint64
	system    System
	fregister [32]float64
	Fcsr      *FCSR
}

func (c *CPU) GetCSR(i uint64) uint64               { return c.csr[i] }
func (c *CPU) SetCSR(i uint64, u uint64)            { c.csr[i] = u }
func (c *CPU) GetLoadReservation() uint64           { return c.lraddr }
func (c *CPU) SetLoadReservation(a uint64)          { c.lraddr = a }
func (c *CPU) GetMemory() *Memory                   { return c.memory }
func (c *CPU) SetMemory(m *Memory)                  { c.memory = m }
func (c *CPU) GetPC() uint64                        { return c.pc }
func (c *CPU) SetPC(i uint64)                       { c.pc = i }
func (c *CPU) GetStatus() uint64                    { return c.status }
func (c *CPU) SetStatus(i uint64)                   { c.status = i }
func (c *CPU) GetSystem() System                    { return c.system }
func (c *CPU) SetSystem(s System)                   { c.system = s }
func (c *CPU) SetRegister(i uint64, u uint64)       { c.register[i] = Cond(i == Rzero, 0x00, u) }
func (c *CPU) GetRegister(i uint64) uint64          { return Cond(i == Rzero, 0x00, c.register[i]) }
func (c *CPU) SetFloatRegister(i uint64, f float64) { c.fregister[i] = f }
func (c *CPU) GetFloatRegister(i uint64) float64    { return c.fregister[i] }

func Cond(b bool, y uint64, f uint64) uint64 {
	if b {
		return y
	}
	return f
}
