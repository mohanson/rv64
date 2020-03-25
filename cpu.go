package rv64

import (
	"math"
)

type CPU struct {
	fasten Fasten
	system System
	csr    CSR
	reg0   [32]uint64
	reg1   [32]uint64
	pc     uint64
	lraddr uint64
	status uint64
}

func (c *CPU) GetCSR() CSR    { return c.csr }
func (c *CPU) SetCSR(csr CSR) { c.csr = csr }

func (c *CPU) GetLoadReservation() uint64  { return c.lraddr }
func (c *CPU) SetLoadReservation(a uint64) { c.lraddr = a }

func (c *CPU) GetMemory() *Memory { return &Memory{Fasten: c.fasten} }
func (c *CPU) SetFasten(f Fasten) { c.fasten = f }

func (c *CPU) GetPC() uint64  { return c.pc }
func (c *CPU) SetPC(i uint64) { c.pc = i }

func (c *CPU) GetStatus() uint64  { return c.status }
func (c *CPU) SetStatus(i uint64) { c.status = i }

func (c *CPU) GetSystem() System  { return c.system }
func (c *CPU) SetSystem(s System) { c.system = s }

func (c *CPU) SetRegister(i uint64, u uint64) {
	if i == Rzero {
		return
	}
	c.reg0[i] = u
}
func (c *CPU) GetRegister(i uint64) uint64 {
	if i == Rzero {
		return 0x00
	}
	return c.reg0[i]
}

func (c *CPU) SetRegisterFloat(i uint64, f uint64) { c.reg1[i] = f }
func (c *CPU) GetRegisterFloat(i uint64) uint64    { return c.reg1[i] }

func (c *CPU) SetRegisterFloatAsFloat64(i uint64, f float64) { c.reg1[i] = math.Float64bits(f) }
func (c *CPU) GetRegisterFloatAsFLoat64(i uint64) float64    { return math.Float64frombits(c.reg1[i]) }

func (c *CPU) SetRegisterFloatAsFloat32(i uint64, f float32) {
	c.SetRegisterFloatAsFloat64(i, NaNBoxing(f))
}
func (c *CPU) GetRegisterFloatAsFLoat32(i uint64) float32 {
	return NaNGnixob(c.GetRegisterFloatAsFLoat64(i))
}

func (c *CPU) SetFloatFlag(flag uint64, b int) {
	if b == 0 {
		flag = ^flag
		c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)&flag)
	} else {
		c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)|flag)
	}
}
func (c *CPU) ClrFloatFlag() {
	c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)&0xffffffffffffffe0)
}

func NewCPU() *CPU {
	return &CPU{}
}
