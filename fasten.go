package rv64

// Fasten is the interface that groups the basic Get, Set and Len methods.
type Fasten interface {
	Get(uint64) (byte, error)
	Set(uint64, byte) error
	Len() uint64
}
