package rv64

import (
	"encoding/binary"
)

// Fasten is the interface that groups the basic Get, Set and Len methods.
type Fasten interface {
	Get(uint64) (byte, error)
	Set(uint64, byte) error
	Len() uint64
}

// Linear is a very simple memory implementation that maps data completely into a byte array
type Linear struct {
	data []byte
}

func (l *Linear) Get(a uint64) (byte, error) {
	if a >= l.Len() {
		return 0x00, ErrOutOfMemory
	}
	return l.data[a], nil
}

func (l *Linear) Set(a uint64, v byte) error {
	if a >= l.Len() {
		return ErrOutOfMemory
	}
	l.data[a] = v
	return nil
}

func (l *Linear) Len() uint64 {
	return uint64(len(l.data))
}

// Memory is a package for Fasten.
type Memory struct {
	Fasten
}

func (m *Memory) GetByte(a uint64, l uint64) ([]byte, error) {
	r := make([]byte, l)
	for i := 0; uint64(i) < l; i++ {
		b, err := m.Fasten.Get(a + uint64(i))
		if err != nil {
			return r, err
		}
		r[i] = b
	}
	return r, nil
}

func (m *Memory) SetByte(a uint64, b []byte) error {
	for i := 0; i < len(b); i++ {
		err := m.Set(a+uint64(i), b[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Memory) GetUint8(a uint64) (uint8, error) {
	mem, err := m.GetByte(a, 1)
	if err != nil {
		return 0, err
	}
	return mem[0], nil
}

func (m *Memory) SetUint8(a uint64, n uint8) error {
	return m.SetByte(a, []byte{n})
}

func (m *Memory) GetUint16(a uint64) (uint16, error) {
	mem, err := m.GetByte(a, 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(mem), nil
}

func (m *Memory) SetUint16(a uint64, n uint16) error {
	mem := make([]byte, 2)
	binary.LittleEndian.PutUint16(mem, n)
	return m.SetByte(a, mem)
}

func (m *Memory) GetUint32(a uint64) (uint32, error) {
	mem, err := m.GetByte(a, 4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(mem), nil
}

func (m *Memory) SetUint32(a uint64, n uint32) error {
	mem := make([]byte, 4)
	binary.LittleEndian.PutUint32(mem, n)
	return m.SetByte(a, mem)
}

func (m *Memory) GetUint64(a uint64) (uint64, error) {
	mem, err := m.GetByte(a, 8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(mem), nil
}

func (m *Memory) SetUint64(a uint64, n uint64) error {
	mem := make([]byte, 8)
	binary.LittleEndian.PutUint64(mem, n)
	return m.SetByte(a, mem)
}

func NewMemoryLinear(size uint64) *Memory {
	return &Memory{Fasten: &Linear{data: make([]byte, size)}}
}
