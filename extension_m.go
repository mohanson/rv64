package rv64

import (
	"math"
)

func ExecuterM(c *CPU, i uint64) (uint64, error) {
	switch {
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0000_0000_0011_0011: // MUL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MUL", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0001_0000_0011_0011: // MULH
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULH", rd, rs1, rs2)
		v := func() uint64 {
			n1, n2 := int64(c.GetRegister(rs1)), int64(c.GetRegister(rs2))
			var neg1, neg2 bool
			if n1 < 0 {
				neg1, n1 = true, -n1
			}
			if n2 < 0 {
				neg2, n2 = true, -n2
			}
			ah, al := uint64(n1)>>32, uint64(n1)&0xffffffff
			bh, bl := uint64(n2)>>32, uint64(n2)&0xffffffff
			a := ah * bh
			b := ah * bl
			c := al * bh
			d := al * bl
			v := a + b>>32 + c>>32 + (d>>32+b&0xffffffff+c&0xffffffff)>>32

			if neg1 != neg2 {
				v = -v
			}
			return v
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0010_0000_0011_0011: // MULHSU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHSU", rd, rs1, rs2)
		v := func() uint64 {
			n1, n2 := int64(c.GetRegister(rs1)), c.GetRegister(rs2)
			var neg bool
			if n1 < 0 {
				neg, n1 = true, -n1
			}

			ah, al := uint64(n1)>>32, uint64(n1)&0xffffffff
			bh, bl := n2>>32, n2&0xffffffff
			a := ah * bh
			b := ah * bl
			c := al * bh
			d := al * bl
			v := a + b>>32 + c>>32 + (d>>32+b&0xffffffff+c&0xffffffff)>>32

			if neg {
				v = -v
			}
			return v
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0011_0000_0011_0011: // MULHU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHU", rd, rs1, rs2)
		v := func() uint64 {
			ah, al := c.GetRegister(rs1)>>32, c.GetRegister(rs1)&0xffffffff
			bh, bl := c.GetRegister(rs2)>>32, c.GetRegister(rs2)&0xffffffff
			a := ah * bh
			b := ah * bl
			c := al * bh
			d := al * bl
			v := a + b>>32 + c>>32 + (d>>32+b&0xffffffff+c&0xffffffff)>>32
			return v
		}()
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_0011: // DIV
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIV", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))/int64(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_0011: // DIVU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVU", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_0011: // REM
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REM", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))%int64(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_0011: // REMU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMU", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, c.GetRegister(rs1)%c.GetRegister(rs2))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0000_0000_0011_1011: // MULW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_1011: // DIVW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_1011: // DIVUW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVUW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_1011: // REMW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMW", rd, rs1, rs2)
		if c.GetRegister(rs2) == 0 {
			c.SetRegister(rd, math.MaxUint64)
		} else {
			c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))))
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_1011: // REMUW
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
