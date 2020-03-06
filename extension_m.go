package rv64

import (
	"math"
	"math/big"
)

func ExecuterM(c *CPU, i uint64) (uint64, error) {
	switch {
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0000_0000_0011_0011: // MUL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MUL", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0001_0000_0011_0011: // MULH
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULH", rd, rs1, rs2)
		v := func() uint64 {
			ag1 := big.NewInt(int64(c.GetRegister(rs1)))
			ag2 := big.NewInt(int64(c.GetRegister(rs2)))
			tmp := big.NewInt(0)
			tmp.Mul(ag1, ag2)
			tmp.Rsh(tmp, 64)
			return uint64(tmp.Int64())
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0010_0000_0011_0011: // MULHSU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHSU", rd, rs1, rs2)
		v := func() uint64 {
			ag1 := big.NewInt(int64(c.GetRegister(rs1)))
			ag2 := big.NewInt(int64(c.GetRegister(rs2)))
			if ag2.Cmp(big.NewInt(0)) == -1 {
				tmp := big.NewInt(0)
				tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
				tmp.Add(tmp, big.NewInt(2))
				ag2 = tmp.Add(tmp, ag2)
			}
			tmp := big.NewInt(0)
			tmp.Mul(ag1, ag2)
			tmp.Rsh(tmp, 64)
			return uint64(tmp.Int64())
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0011_0000_0011_0011: // MULHU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHU", rd, rs1, rs2)
		v := func() uint64 {
			ag1 := big.NewInt(int64(c.GetRegister(rs1)))
			ag2 := big.NewInt(int64(c.GetRegister(rs2)))
			if ag1.Cmp(big.NewInt(0)) == -1 {
				tmp := big.NewInt(0)
				tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
				tmp.Add(tmp, big.NewInt(2))
				ag1 = tmp.Add(tmp, ag1)
			}
			if ag2.Cmp(big.NewInt(0)) == -1 {
				tmp := big.NewInt(0)
				tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
				tmp.Add(tmp, big.NewInt(2))
				ag2 = tmp.Add(tmp, ag2)
			}
			tmp := big.NewInt(0)
			tmp.Mul(ag1, ag2)
			tmp.Rsh(tmp, 64)
			return tmp.Uint64()
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_0011: // DIV
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIV", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))/int64(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_0011: // DIVU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVU", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_0011: // REM
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REM", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))%int64(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_0011: // REMU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMU", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, c.GetRegister(rs1)%c.GetRegister(rs2))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0000_0000_0011_1011: // MULW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_1011: // DIVW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_1011: // DIVUW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVUW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_1011: // REMW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_1011: // REMUW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMUW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	return 0, nil
}
