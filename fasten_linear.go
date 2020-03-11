package rv64

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
