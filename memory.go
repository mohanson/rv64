package rv64

import "encoding/binary"

type MemoryCore interface {
	Get(uint64, uint64) ([]byte, error)
	Set(uint64, []byte) error
	Len() uint64
}

type Memory struct {
	MemoryCore
}

func (m *Memory) GetUint8(a uint64) (uint8, error) {
	mem, err := m.Get(a, 1)
	if err != nil {
		return 0, err
	}
	return mem[0], nil
}

func (m *Memory) SetUint8(a uint64, n uint8) error {
	return m.Set(a, []byte{n})
}

func (m *Memory) GetUint16(a uint64) (uint16, error) {
	mem, err := m.Get(a, 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(mem), nil
}

func (m *Memory) SetUint16(a uint64, n uint16) error {
	mem := make([]byte, 2)
	binary.LittleEndian.PutUint16(mem, n)
	return m.Set(a, mem)
}

func (m *Memory) GetUint32(a uint64) (uint32, error) {
	mem, err := m.Get(a, 4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(mem), nil
}

func (m *Memory) SetUint32(a uint64, n uint32) error {
	mem := make([]byte, 4)
	binary.LittleEndian.PutUint32(mem, n)
	return m.Set(a, mem)
}

func (m *Memory) GetUint64(a uint64) (uint64, error) {
	mem, err := m.Get(a, 8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(mem), nil
}

func (m *Memory) SetUint64(a uint64, n uint64) error {
	mem := make([]byte, 8)
	binary.LittleEndian.PutUint64(mem, n)
	return m.Set(a, mem)
}

type MemoryCoreLinear struct {
	data []byte
}

func (m *MemoryCoreLinear) Get(a uint64, n uint64) ([]byte, error) {
	l := uint64(len(m.data))
	e := a + n
	if e < a || e > l {
		return nil, ErrOutOfMemory
	}
	return m.data[a:e], nil
}

func (m *MemoryCoreLinear) Set(a uint64, b []byte) error {
	l := uint64(len(m.data))
	e := a + uint64(len(b))
	if e < a || e > l {
		return ErrOutOfMemory
	}
	copy(m.data[a:e], b)
	return nil
}

func (m *MemoryCoreLinear) Len() uint64 {
	return uint64(len(m.data))
}

func NewMemoryLinear(size uint64) *Memory {
	return &Memory{MemoryCore: &MemoryCoreLinear{data: make([]byte, size)}}
}
