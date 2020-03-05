package rv64

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
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0010_0000_0011_0011: // MULHSU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHSU", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0011_0000_0011_0011: // MULHU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULHU", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_0011: // DIV
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIV", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_0011: // DIVU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVU", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_0011: // REM
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REM", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_0011: // REMU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMU", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0000_0000_0011_1011: // MULW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("MULW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0100_0000_0011_1011: // DIVW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVW", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0101_0000_0011_1011: // DIVUW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("DIVUW", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0110_0000_0011_1011: // REMW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMW", rd, rs1, rs2)
	case i&0b_1111_1110_0000_0111_0000_0000_0111_1111 == 0b_0000_0010_0000_0000_0111_0000_0011_1011: // REMUW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("REMUW", rd, rs1, rs2)
	}
	return 0, nil
}
