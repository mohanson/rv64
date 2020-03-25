package rv64

// Control and Status Registers.
//
// Number Privilege  Name     Description
// 0x001  Read/write fflags   Floating-Point Accrued Exceptions.
// 0x002  Read/write frm      Floating-Point Dynamic Rounding Mode.
// 0x003  Read/write fcsr     Floating-Point Control and Status Register (frm + fflags).
// 0xC00  Read-only  cycle    Cycle counter for RDCYCLE instruction.
// 0xC01  Read-only  time     Timer for RDTIME instruction.
// 0xC02  Read-only  instret  Instructions-retired counter for RDINSTRET instruction.
// 0xC80  Read-only  cycleh   Upper 32 bits of cycle, RV32I only.
// 0xC81  Read-only  timeh    Upper 32 bits of time, RV32I only.
// 0xC82  Read-only  instreth Upper 32 bits of instret, RV32I only.

type CSR interface {
	Get(uint64) uint64
	Set(uint64, uint64)
}

type CSRDaze struct {
	m [0x1000]uint64
}

func (c *CSRDaze) Get(i uint64) uint64 {
	switch {
	case i == CSRfflags:
		return c.m[CSRfcsr] & 0x1f
	case i == CSRfrm:
		return c.m[CSRfcsr] & 0xe0 >> 5
	case i == i:
		return c.m[i]
	}
	return c.m[i]
}

func (c *CSRDaze) Set(i uint64, u uint64) {
	switch {
	case i == CSRfcsr:
		c.m[i] = u & 0xff
	case i == CSRfflags:
		c.m[i] = u & 0x1f
		c.m[CSRfcsr] = (c.m[CSRfcsr] >> 5 << 5) | (u & 0x1f)
	case i == CSRfrm:
		c.m[i] = u & 0x07
		c.m[CSRfcsr] = c.m[CSRfcsr]&0xffffffffffffff1f | ((u & 0x07) << 5)
	case i == i:
		c.m[i] = u
	}
}
