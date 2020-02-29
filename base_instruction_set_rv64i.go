package riscv

import (
	"encoding/binary"
	"log"
)

func ExecuterRV64I(r *RegisterRV64I, m []byte, i uint64) int {
	switch {
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0000_0011: // LWU
		// I
		log.Println("LWU")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0000_0011: // LD
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LD", rd, rs1, imm)
		a := r.RG[rs1] + imm
		r.RG[rd] = binary.LittleEndian.Uint64(m[a : a+8])
		r.PC += 4
		return 1
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_0011: // SD
		// rd, rs1, imm := IType(i)
		// imm = SignExtend(imm, 11)
		// DebuglnIType("SD", rd, rs1, imm)
		// a := r.RG[rs1] + imm
		// r.RG[rd] = binary.LittleEndian.Uint64(m[a : a+8])
		// r.PC += 4
		// return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_0011: // SLLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLI", rd, rs1, imm)
		r.RG[rd] = r.RG[rs1] << imm
		r.PC += 4
		return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_0011: // SRLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLI", rd, rs1, imm)
		r.RG[rd] = r.RG[rs1] >> imm
		r.PC += 4
		return 1
	case i&0b_1111_1100_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_0011: // SRAI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAI", rd, rs1, imm)
		r.RG[rd] = uint64(int64(r.RG[rs1]) >> imm)
		r.PC += 4
		return 1
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_1011: // ADDIW
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDIW", rd, rs1, imm)
		r.RG[rd] = uint64(int32(r.RG[rs1]) + int32(imm))
		r.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_1011: // SLLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLIW", rd, rs1, imm)
		r.RG[rd] = SignExtend(uint64(uint32(r.RG[rs1]<<imm)), 31)
		r.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_1011: // SRLIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLIW", rd, rs1, imm)
		r.RG[rd] = SignExtend(uint64(uint32(r.RG[rs1]>>imm)), 31)
		r.PC += 4
		return 1
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_1011: // SRAIW
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAIW", rd, rs1, imm)
		r.RG[rd] = uint64(int64(r.RG[rs1]) >> imm)
		r.PC += 4
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
