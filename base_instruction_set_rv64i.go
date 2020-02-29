package riscv

import (
	"encoding/binary"
	"log"
)

func ExecuterRV64I(c *CPU, i uint64) int {
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
		return 1
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_0011: // SD
		// rd, rs1, imm := IType(i)
		// imm = SignExtend(imm, 11)
		// DebuglnIType("SD", rd, rs1, imm)
		// a := c.Register[rs1] + imm
		// c.Register[rd] = binary.LittleEndian.Uint64(m[a : a+8])
		// c.PC += 4
		// return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_0011: // SLLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] << imm
		c.PC += 4
		return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_0011: // SRLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] >> imm
		c.PC += 4
		return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_0011: // SRAI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAI", rd, rs1, imm)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> imm)
		c.PC += 4
		return 1
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_1011: // ADDIW
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDIW", rd, rs1, imm)
		c.Register[rd] = uint64(int32(c.Register[rs1]) + int32(imm))
		c.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_1011: // SLLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLIW", rd, rs1, imm)
		c.Register[rd] = SignExtend(uint64(uint32(c.Register[rs1]<<imm)), 31)
		c.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_1011: // SRLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLIW", rd, rs1, imm)
		c.Register[rd] = SignExtend(uint64(uint32(c.Register[rs1]>>imm)), 31)
		c.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_1011: // SRAIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAIW", rd, rs1, imm)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> imm)
		c.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_1011: // ADDW
		// r
		log.Println("ADDW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0000_0000_0011_1011: // SUBW
		// r
		log.Println("SUBW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0011_1011: // SLLW
		// r
		log.Println("SLLW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0011_1011: // SRLW
		// r
		log.Println("SRLW")
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0011_1011: // SRAW
		log.Println("SRAW")
	}
	return 0
}
