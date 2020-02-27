package riscv

import "log"

func ExecuterC(r *RegisterRV64I, i uint64) int {
	switch {
	case i&0b_1111_1111_1111_1111 == 0b_0000_0000_0000_0000: // Illegal instruction
	case i&0b_1110_0000_0000_0011 == 0b_0000_0000_0000_0000: // C.ADDI4SPN
	case i&0b_1110_0000_0000_0011 == 0b_0010_0000_0000_0000: // C.FLD
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0000: // C.LW
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0000: // C.LD
	case i&0b_1110_0000_0000_0011 == 0b_1000_0000_0000_0000: // Reserved
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0000: // C.FSD
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0000: // C.SW
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0000: // C.SD

	case i&0b_1111_1111_1111_1111 == 0b_0000_0000_0000_0001: // C.NOP
	case i&0b_1110_0000_0000_0011 == 0b_0000_0000_0000_0001: // C.ADDI
	case i&0b_1110_0000_0000_0011 == 0b_0010_0000_0000_0001: // C.ADDIW
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0001: // C.LI
		var (
			rd  = int(InstructionPart(i, 7, 11))
			imm = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
		)
		if rd == 0 {
			log.Panicln("")
		}
		DebuglnIType("C.LI", rd, rd, SignExtend(imm, 5))
		r.RG[rd] = r.RG[rd] + SignExtend(imm, 5)
		r.PC += 2
		return 1
	case i&0b_1110_1111_1000_0011 == 0b_0110_0001_0000_0001: // C.ADDI16SP
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0001: // C.LUI
	case i&0b_1111_1100_0111_1111 == 0b_1000_0000_0000_0001: // C.SRLI64
	case i&0b_1111_1100_0111_1111 == 0b_1000_0100_0000_0001: // C.SRAI64
	case i&0b_1110_1100_0000_0011 == 0b_1000_1000_0000_0001: // C.ANDI
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0000_0001: // C.SUB
		var (
			rd  = int(InstructionPart(i, 7, 9))
			rs2 = int(InstructionPart(i, 2, 4))
		)
		DebuglnRType("C.SUB", rd+8, rd+8, rs2+8)
		r.RG[rd+8] = r.RG[rd+8] - r.RG[rs2+8]
		r.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0010_0001: // C.XOR
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0100_0001: // C.OR
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0110_0001: // C.AND
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0000_0001: // C.SUBW
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0010_0001: // C.ADDW
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0100_0001: // Reserved
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0110_0001: // Reserved
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0001: // C.J
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0001: // C.BEQZ
		var (
			rs1 = int(InstructionPart(i, 7, 9))
			imm = InstructionPart(i, 12, 12)<<8 | InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 2, 2)<<5 | InstructionPart(i, 10, 11)<<3 | InstructionPart(i, 3, 4)<<1
		)
		DebuglnBType("C.BNEZ", rs1+8, Rzero, imm)
		if r.RG[rs1+8] == r.RG[Rzero] {
			r.PC = r.PC + SignExtend(imm, 8)
		} else {
			r.PC += 2
		}
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0001: // C.BNEZ
		var (
			rs1 = int(InstructionPart(i, 7, 9))
			imm = InstructionPart(i, 12, 12)<<8 | InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 2, 2)<<5 | InstructionPart(i, 10, 11)<<3 | InstructionPart(i, 3, 4)<<1
		)
		DebuglnBType("C.BNEZ", rs1+8, Rzero, imm)
		if r.RG[rs1+8] != r.RG[Rzero] {
			r.PC = r.PC + SignExtend(imm, 8)
		} else {
			r.PC += 2
		}
		return 1
	case i&0b_1111_0000_0111_1111 == 0b_0000_0000_0000_0010: // C.SLLI64
	case i&0b_1111_0000_0000_0011 == 0b_0010_0000_0000_0010: // C.FLDSP
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0010: // C.LWSP
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0010: // C.LDSP
	case i&0b_1111_0000_0111_1111 == 0b_1000_0000_0000_0010: // C.JR
	case i&0b_1111_0000_0000_0011 == 0b_1000_0000_0000_0010: // C.MV
		var (
			rd  = int(InstructionPart(i, 7, 11))
			rs2 = int(InstructionPart(i, 2, 6))
		)
		if rd == 0 || rs2 == 0 {
			log.Panicln("")
		}
		DebuglnRType("C.MV", rd, Rzero, rs2)
		r.RG[rd] = r.RG[rs2]
		r.PC += 2
		return 1
	case i&0b_1111_1111_1111_1111 == 0b_1001_0000_0000_0010: // C.EBREAK
	case i&0b_1111_0000_0111_1111 == 0b_1001_0000_0000_0010: // C.JALR
	case i&0b_1111_0000_0000_0011 == 0b_1001_0000_0000_0010: // C.ADD
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0010: // C.FSDSP
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0010: // C.SWSP
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0010: // C.SDSP
	}

	return 0
}
