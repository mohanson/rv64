package riscv

type Memory interface {
	Get(uint64, uint64) ([]byte, error)
	Set(uint64, []byte) error
	Len() uint64
}

type MemoryLinear struct {
	data []byte
}

func (m *MemoryLinear) Get(a uint64, n uint64) ([]byte, error) {
	l := uint64(len(m.data))
	e := a + n
	if e < a || e > l {
		return nil, ErrOutOfMemory
	}
	return m.data[a:e], nil
}

func (m *MemoryLinear) Set(a uint64, b []byte) error {
	l := uint64(len(m.data))
	e := a + uint64(len(b))
	if e < a || e > l {
		return ErrOutOfMemory
	}
	copy(m.data[a:e], b)
	return nil
}

func (m *MemoryLinear) Len() uint64 {
	return uint64(len(m.data))
}

func NewMemoryLinear(size uint64) Memory {
	return &MemoryLinear{data: make([]byte, size)}
}
