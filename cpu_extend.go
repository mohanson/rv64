package riscv

func Cond(b bool, y uint64, f uint64) uint64 {
	if b {
		return y
	}
	return f
}

// func (c *CPU) GetMemoryUint8(a uint64) (uint8, error) {
// 	mem, err := c.memory.Get(a, 1)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return mem[0], nil
// }

// func (c *CPU) SetMemoryUint8(a uint64, n uint8) error {
// 	return c.memory.Set(a, []byte{n})
// }

// func (c *CPU) GetMemoryUint16(a uint64) (uint16, error) {
// 	mem, err := c.memory.Get(a, 2)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return binary.LittleEndian.Uint16(mem), nil
// }

// func (c *CPU) SetMemoryUint16(a uint64, n uint16) error {
// 	mem := make([]byte, 2)
// 	binary.LittleEndian.PutUint16(mem, n)
// 	return c.memory.Set(a, mem)
// }

// func (c *CPU) GetMemoryUint32(a uint64) (uint32, error) {
// 	mem, err := c.memory.Get(a, 4)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return binary.LittleEndian.Uint32(mem), nil
// }

// func (c *CPU) SetMemoryUint32(a uint64, n uint32) error {
// 	mem := make([]byte, 4)
// 	binary.LittleEndian.PutUint32(mem, n)
// 	return c.memory.Set(a, mem)
// }

// func (c *CPU) GetMemoryUint64(a uint64) (uint64, error) {
// 	mem, err := c.memory.Get(a, 8)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return binary.LittleEndian.Uint64(mem), nil
// }

// func (c *CPU) SetMemoryUint64(a uint64, n uint64) error {
// 	mem := make([]byte, 8)
// 	binary.LittleEndian.PutUint64(mem, n)
// 	return c.memory.Set(a, mem)
// }
