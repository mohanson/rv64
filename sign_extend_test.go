package rv64

import (
	"testing"
)

func TestSignExtend(t *testing.T) {
	data := [][]uint64{
		[]uint64{0x0000000000000005, 0x02, 0xfffffffffffffffd},
		[]uint64{0x0000000000000005, 0x03, 0x0000000000000005},
		[]uint64{0x8000000000000005, 0x40, 0x8000000000000005},
		[]uint64{0x80000000000000a5, 0x07, 0xffffffffffffffa5},
	}
	for _, l := range data {
		t.Log(SignExtend(l[0], l[1]), l[2])
		if SignExtend(l[0], l[1]) != l[2] {
			t.Fail()
		}
	}
}
