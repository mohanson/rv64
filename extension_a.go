package rv64

func ExecuterA(c *CPU, i uint64) (uint64, error) {
	switch {
	case i&0b_1111_1001_1111_0000_0111_0000_0111_1111 == 0b_0001_0000_0000_0000_0010_0000_0010_1111: // LR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("LR.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0001_1000_0000_0000_0010_0000_0010_1111: // SC.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SC.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_1000_0000_0000_0010_0000_0010_1111: // AMOSWAP.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOSWAP.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0010_1111: // AMOADD.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOADD.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0010_0000_0000_0000_0010_0000_0010_1111: // AMOXOR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOXOR.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0110_0000_0000_0000_0010_0000_0010_1111: // AMOAND.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOAND.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0010_0000_0010_1111: // AMOOR.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOOR.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1000_0000_0000_0000_0010_0000_0010_1111: // AMOMIN.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMIN.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1010_0000_0000_0000_0010_0000_0010_1111: // AMOMAX.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAX.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1100_0000_0000_0000_0010_0000_0010_1111: // AMOMINU.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMINU.W", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1110_0000_0000_0000_0010_0000_0010_1111: // AMOMAXU.W
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAXU.W", rd, rs1, rs2)
	case i&0b_1111_1001_1111_0000_0111_0000_0111_1111 == 0b_0001_0000_0000_0000_0011_0000_0010_1111: // LR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("LR.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0001_1000_0000_0000_0011_0000_0010_1111: // SC.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SC.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_1000_0000_0000_0011_0000_0010_1111: // AMOSWAP.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOSWAP.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_1111: // AMOADD.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOADD.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0010_0000_0000_0000_0011_0000_0010_1111: // AMOXOR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOXOR.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0110_0000_0000_0000_0011_0000_0010_1111: // AMOAND.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOAND.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0011_0000_0010_1111: // AMOOR.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOOR.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1000_0000_0000_0000_0011_0000_0010_1111: // AMOMIN.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMIN.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1010_0000_0000_0000_0011_0000_0010_1111: // AMOMAX.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAX.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1100_0000_0000_0000_0011_0000_0010_1111: // AMOMINU.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMINU.D", rd, rs1, rs2)
	case i&0b_1111_1000_0000_0000_0111_0000_0111_1111 == 0b_1110_0000_0000_0000_0011_0000_0010_1111: // AMOMAXU.D
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AMOMAXU.D", rd, rs1, rs2)
	}
	return 0, nil
}
