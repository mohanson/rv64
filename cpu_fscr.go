package rv64

func (c *CPU) SetFloatFlag(flag uint64, b int) {
	if b == 0 {
		c.csr[CSRfflags] = c.csr[CSRfflags] & (^flag)
		c.csr[CSRfcsr] = c.csr[CSRfcsr] & (^flag)
	} else {
		c.csr[CSRfflags] = c.csr[CSRfflags] | flag
		c.csr[CSRfcsr] = c.csr[CSRfcsr] | flag
	}
}

func (c *CPU) ClrFloatFlag() {
	c.csr[CSRfflags] = 0x00
	c.csr[CSRfcsr] = c.csr[CSRfcsr] & 0xffffffffffffffe0
}
