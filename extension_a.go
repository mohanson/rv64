package rv64

var (
	lr uint64
)

func ExecuterA(c *CPU, i uint64) (uint64, error) {
	switch {
	case i&0b_1111_1001_1111_0000_0111_0000_0111_1111 == 0b_0001_0000_0000_0000_0010_0000_0010_1111: // LR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("LR.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		v, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, SignExtend(uint64(v), 31))
		lr = a
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0001_1000_0000_0000_0010_0000_0010_1111: // SC.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SC.W", rd, rs1, rs2)

		rs1_val := c.GetRegister(rs1)
		rs2_val := uint32(c.GetRegister(rs2))
		mem_addr := SignExtend(rs1_val, 31)

		if mem_addr != lr {
			c.SetRegister(rd, 1)
			lr = 0
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		}

		c.GetMemory().SetUint32(mem_addr, rs2_val)
		c.SetRegister(rd, 0)
		lr = 0
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_1000_0000_0000_0010_0000_0010_1111: // AMOSWAP.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOSWAP.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		r := uint32(c.GetRegister(rs2))
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0010_1111: // AMOADD.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOADD.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		r := b + uint32(c.GetRegister(rs2))
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0010_0000_0000_0000_0010_0000_0010_1111: // AMOXOR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOXOR.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		r := b ^ uint32(c.GetRegister(rs2))
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0110_0000_0000_0000_0010_0000_0010_1111: // AMOAND.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOAND.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		r := b & uint32(c.GetRegister(rs2))
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0010_0000_0010_1111: // AMOOR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOOR.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		r := b | uint32(c.GetRegister(rs2))
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1000_0000_0000_0000_0010_0000_0010_1111: // AMOMIN.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMIN.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		var r uint32
		if int32(b) < int32(uint32(c.GetRegister(rs2))) {
			r = b
		} else {
			r = uint32(c.GetRegister(rs2))
		}
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1010_0000_0000_0000_0010_0000_0010_1111: // AMOMAX.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAX.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		var r uint32
		if int32(b) > int32(uint32(c.GetRegister(rs2))) {
			r = b
		} else {
			r = uint32(c.GetRegister(rs2))
		}
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1100_0000_0000_0000_0010_0000_0010_1111: // AMOMINU.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMINU.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		var r uint32
		if b < uint32(c.GetRegister(rs2)) {
			r = b
		} else {
			r = uint32(c.GetRegister(rs2))
		}
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1110_0000_0000_0000_0010_0000_0010_1111: // AMOMAXU.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAXU.W", rd, rs1, rs2)
		a := SignExtend(c.GetRegister(rs1), 31)
		b, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		var r uint32
		if b > uint32(c.GetRegister(rs2)) {
			r = b
		} else {
			r = uint32(c.GetRegister(rs2))
		}
		c.GetMemory().SetUint32(a, r)
		c.SetRegister(rd, SignExtend(uint64(b), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1001_1111_0000_0111_0000_0111_1111 == 0b_0001_0000_0000_0000_0011_0000_0010_1111: // LR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("LR.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0001_1000_0000_0000_0011_0000_0010_1111: // SC.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SC.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		c.GetMemory().SetUint64(a, c.GetRegister(rs2))
		c.SetRegister(rd, 0x00)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_1000_0000_0000_0011_0000_0010_1111: // AMOSWAP.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOSWAP.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.GetMemory().SetUint64(a, c.GetRegister(rs2))
		c.SetRegister(rs2, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_1111: // AMOADD.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOADD.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.GetMemory().SetUint64(a, v+c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0010_0000_0000_0000_0011_0000_0010_1111: // AMOXOR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOXOR.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.GetMemory().SetUint64(a, v^c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0110_0000_0000_0000_0011_0000_0010_1111: // AMOAND.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOAND.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.GetMemory().SetUint64(a, v&c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0011_0000_0010_1111: // AMOOR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOOR.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.GetMemory().SetUint64(a, v|c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1000_0000_0000_0000_0011_0000_0010_1111: // AMOMIN.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMIN.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		var w uint64 = 0
		if int64(v) < int64(c.GetRegister(rs2)) {
			w = v
		} else {
			w = c.GetRegister(rs2)
		}
		c.GetMemory().SetUint64(a, w)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1010_0000_0000_0000_0011_0000_0010_1111: // AMOMAX.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAX.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		var w uint64 = 0
		if int64(v) > int64(c.GetRegister(rs2)) {
			w = v
		} else {
			w = c.GetRegister(rs2)
		}
		c.GetMemory().SetUint64(a, w)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1100_0000_0000_0000_0011_0000_0010_1111: // AMOMINU.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMINU.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		var w uint64 = 0
		if v < c.GetRegister(rs2) {
			w = v
		} else {
			w = c.GetRegister(rs2)
		}
		c.GetMemory().SetUint64(a, w)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1110_0000_0000_0000_0011_0000_0010_1111: // AMOMAXU.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAXU.D", rd, rs1, rs2)
		a := c.GetRegister(rs1)
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		var w uint64 = 0
		if v > c.GetRegister(rs2) {
			w = v
		} else {
			w = c.GetRegister(rs2)
		}
		c.GetMemory().SetUint64(a, w)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	return 0, nil
}
