package riscv

import (
	"encoding/binary"
	"log"
)

func ExecuterRV64I(c *CPU, i uint64) (int, error) {
	m := c.Memory
	switch {
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0000_0011: // LWU
		// I
		log.Println("LWU")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0000_0011: // LD
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LD", rd, rs1, imm)
		a := c.Register[rs1] + imm
		c.Register[rd] = binary.LittleEndian.Uint64(m[a : a+8])
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_0011: // SD
		rs1, rs2, imm := SType(i)
		DebuglnIType("SD", rs1, rs2, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		m[a] = byte(c.Register[rs2])
		m[a+1] = byte(c.Register[rs2] >> 8)
		m[a+2] = byte(c.Register[rs2] >> 16)
		m[a+3] = byte(c.Register[rs2] >> 24)
		m[a+4] = byte(c.Register[rs2] >> 32)
		m[a+5] = byte(c.Register[rs2] >> 40)
		m[a+6] = byte(c.Register[rs2] >> 48)
		m[a+7] = byte(c.Register[rs2] >> 56)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_0011: // SLLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] << imm
		c.PC += 4
		return 1, nil
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_0011: // SRLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] >> imm
		c.PC += 4
		return 1, nil
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_0011: // SRAI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAI", rd, rs1, imm)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> imm)
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_1011: // ADDIW
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDIW", rd, rs1, imm)
		c.Register[rd] = uint64(int32(c.Register[rs1]) + int32(imm))
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_1011: // SLLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLIW", rd, rs1, imm)
		c.Register[rd] = SignExtend(uint64(uint32(c.Register[rs1]<<imm)), 31)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_1011: // SRLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLIW", rd, rs1, imm)
		c.Register[rd] = SignExtend(uint64(uint32(c.Register[rs1]>>imm)), 31)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_1011: // SRAIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAIW", rd, rs1, imm)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> imm)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_1011: // ADDW
		// r
		log.Println("ADDW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0000_0000_0011_1011: // SUBW
		// r
		log.Println("SUBW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0011_1011: // SLLW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLLW", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] << InstructionPart(c.Register[rs2], 0, 5)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0011_1011: // SRLW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRLW", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] >> InstructionPart(c.Register[rs2], 0, 5)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0011_1011: // SRAW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRAW", rd, rs1, rs2)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> InstructionPart(c.Register[rs2], 0, 5))
		c.PC += 4
		return 1, nil
	}
	return 0, nil
}
