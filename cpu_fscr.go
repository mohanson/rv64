package rv64

func (c *CPU) SetFloatFlag(flag uint64, b int) {
	if b == 0 {
		c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)&(^flag))
	} else {
		c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)|flag)
	}
}

func (c *CPU) ClrFloatFlag() {
	c.csr.Set(CSRfcsr, c.csr.Get(CSRfcsr)&0xffffffffffffffe0)
}
