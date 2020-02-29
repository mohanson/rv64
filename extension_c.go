package riscv

import (
	"encoding/binary"
	"log"
)

func ExecuterC(c *CPU, i uint64) int {
	m := c.Memory
	switch {
	case i&0b_1111_1111_1111_1111 == 0b_0000_0000_0000_0000: // Illegal instruction
		log.Println("Illegal instruction")
	case i&0b_1110_0000_0000_0011 == 0b_0000_0000_0000_0000: // C.ADDI4SPN
		log.Println("C.ADDI4SPN")
	case i&0b_1110_0000_0000_0011 == 0b_0010_0000_0000_0000: // C.FLD
		log.Println("C.FLD")
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0000: // C.LW
		log.Println("C.LW")
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0000: // C.LD
		log.Println("C.LD")
	case i&0b_1110_0000_0000_0011 == 0b_1000_0000_0000_0000: // Reserved
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0000: // C.FSD
		log.Println("C.FSD")
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0000: // C.SW
		log.Println("C.SW")
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0000: // C.SD
		var (
			rs1 = int(InstructionPart(i, 7, 9)) + 8
			rs2 = int(InstructionPart(i, 2, 4)) + 8
			imm = InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 10, 12)<<3
		)
		DebuglnSType("C.SD", rs1, rs2, imm)
		binary.LittleEndian.PutUint64(m[int(c.Register[rs1]+imm):int(c.Register[rs1]+imm)+8], c.Register[rs2])
		c.PC += 2
		return 1
	case i&0b_1111_1111_1111_1111 == 0b_0000_0000_0000_0001: // C.NOP
		log.Println("C.NOP")
	case i&0b_1110_0000_0000_0011 == 0b_0000_0000_0000_0001: // C.ADDI
		var (
			rd  = int(InstructionPart(i, 7, 11))
			imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6)<<0, 5)
		)
		DebuglnIType("C.ADDI", rd, rd, imm)
		c.Register[rd] = c.Register[rd] + imm
		c.PC += 2
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_0010_0000_0000_0001: // C.ADDIW
		log.Println("C.ADDIW")
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0001: // C.LI
		var (
			rd  = int(InstructionPart(i, 7, 11))
			imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
		)
		if rd == 0 {
			log.Panicln("")
		}
		DebuglnIType("C.LI", rd, rd, imm)
		c.Register[rd] = imm
		c.PC += 2
		return 1
	case i&0b_1110_1111_1000_0011 == 0b_0110_0001_0000_0001: // C.ADDI16SP
		log.Println("C.ADDI16SP")
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0001: // C.LUI
		log.Println("C.LUI")
	case i&0b_1111_1100_0111_1111 == 0b_1000_0000_0000_0001: // C.SRLI64
		log.Println("C.SRLI64")
	case i&0b_1111_1100_0111_1111 == 0b_1000_0100_0000_0001: // C.SRAI64
		log.Println("C.SRAI64")
	case i&0b_1110_1100_0000_0011 == 0b_1000_1000_0000_0001: // C.ANDI
		var (
			rd  = int(InstructionPart(i, 7, 9)) + 8
			imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6)<<0, 5)
		)
		DebuglnIType("C.ANDI", rd, rd, imm)
		c.Register[rd] = c.Register[rd] & imm
		c.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0000_0001: // C.SUB
		var (
			rd  = int(InstructionPart(i, 7, 9)) + 8
			rs2 = int(InstructionPart(i, 2, 4)) + 8
		)
		DebuglnRType("C.SUB", rd, rd, rs2)
		c.Register[rd] = c.Register[rd] - c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0010_0001: // C.XOR
		var (
			rd  = int(InstructionPart(i, 7, 9)) + 8
			rs2 = int(InstructionPart(i, 2, 4)) + 8
		)
		DebuglnRType("C.XOR", rd, rd, rs2)
		c.Register[rd] = c.Register[rd] ^ c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0100_0001: // C.OR
		var (
			rd  = int(InstructionPart(i, 7, 9)) + 8
			rs2 = int(InstructionPart(i, 2, 4)) + 8
		)
		DebuglnRType("C.OR", rd, rd, rs2)
		c.Register[rd] = c.Register[rd] | c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1000_1100_0110_0001: // C.AND
		var (
			rd  = int(InstructionPart(i, 7, 9)) + 8
			rs2 = int(InstructionPart(i, 2, 4)) + 8
		)
		DebuglnRType("C.AND", rd, rd, rs2)
		c.Register[rd] = c.Register[rd] & c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0000_0001: // C.SUBW
		log.Println("C.SUBW")
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0010_0001: // C.ADDW
		log.Println("C.ADDW")
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0100_0001: // Reserved
	case i&0b_1111_1100_0110_0011 == 0b_1001_1100_0110_0001: // Reserved
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0001: // C.J
		var (
			imm = SignExtend(
				InstructionPart(i, 12, 12)<<11|
					InstructionPart(i, 8, 8)<<10|
					InstructionPart(i, 9, 10)<<8|
					InstructionPart(i, 6, 6)<<7|
					InstructionPart(i, 7, 7)<<6|
					InstructionPart(i, 2, 2)<<5|
					InstructionPart(i, 11, 11)<<4|
					InstructionPart(i, 3, 5)<<1,
				11)
		)
		DebuglnJType("C.J", Rzero, imm)
		c.PC += imm
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0001: // C.BEQZ
		var (
			rs1 = int(InstructionPart(i, 7, 9)) + 8
			imm = SignExtend(InstructionPart(i, 12, 12)<<8|InstructionPart(i, 5, 6)<<6|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 10, 11)<<3|InstructionPart(i, 3, 4)<<1, 8)
		)
		DebuglnBType("C.BNEZ", rs1, Rzero, imm)
		if c.Register[rs1] == c.Register[Rzero] {
			c.PC = c.PC + imm
		} else {
			c.PC += 2
		}
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0001: // C.BNEZ
		var (
			rs1 = int(InstructionPart(i, 7, 9)) + 8
			imm = SignExtend(InstructionPart(i, 12, 12)<<8|InstructionPart(i, 5, 6)<<6|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 10, 11)<<3|InstructionPart(i, 3, 4)<<1, 8)
		)
		DebuglnBType("C.BNEZ", rs1, Rzero, imm)
		if c.Register[rs1] != c.Register[Rzero] {
			c.PC = c.PC + imm
		} else {
			c.PC += 2
		}
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_0000_0000_0000_0010: // C.SLLI
		var (
			rd  = int(InstructionPart(i, 7, 11))
			imm = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)<<0
		)
		if rd == 0 || imm == 0 {
			log.Panicln("")
		}
		DebuglnIType("C.SLLI", rd, rd, imm)
		c.Register[rd] = c.Register[rd] << imm
		c.PC += 2
		return 1
	case i&0b_1111_0000_0000_0011 == 0b_0010_0000_0000_0010: // C.FLDSP
		log.Println("C.FLDSP")
	case i&0b_1110_0000_0000_0011 == 0b_0100_0000_0000_0010: // C.LWSP
		log.Println("C.LWSP")
	case i&0b_1110_0000_0000_0011 == 0b_0110_0000_0000_0010: // C.LDSP
		log.Println("C.LDSP")
	case i&0b_1111_0000_0111_1111 == 0b_1000_0000_0000_0010: // C.JR
		var (
			rs1 = int(InstructionPart(i, 7, 11))
		)
		DebuglnIType("C.JR", Rzero, rs1, 0)
		c.PC = c.Register[rs1]
		return 1
	case i&0b_1111_0000_0000_0011 == 0b_1000_0000_0000_0010: // C.MV
		var (
			rd  = int(InstructionPart(i, 7, 11))
			rs2 = int(InstructionPart(i, 2, 6))
		)
		if rd == 0 || rs2 == 0 {
			log.Panicln("")
		}
		DebuglnRType("C.MV", rd, Rzero, rs2)
		c.Register[rd] = c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1111_1111_1111_1111 == 0b_1001_0000_0000_0010: // C.EBREAK
		log.Println("C.EBREAK")
	case i&0b_1111_0000_0111_1111 == 0b_1001_0000_0000_0010: // C.JALR
		var (
			rs1 = int(InstructionPart(i, 7, 11))
		)
		DebuglnIType("C.JALR", Rra, rs1, 0)
		c.PC = c.Register[rs1]
		return 1
	case i&0b_1111_0000_0000_0011 == 0b_1001_0000_0000_0010: // C.ADD
		var (
			rd  = int(InstructionPart(i, 7, 11))
			rs2 = int(InstructionPart(i, 2, 6))
		)
		if rd == 0 || rs2 == 0 {
			log.Panicln("")
		}
		DebuglnRType("C.ADD", rd, rd, rs2)
		c.Register[rd] += c.Register[rs2]
		c.PC += 2
		return 1
	case i&0b_1110_0000_0000_0011 == 0b_1010_0000_0000_0010: // C.FSDSP
		log.Println("C.FSDSP")
	case i&0b_1110_0000_0000_0011 == 0b_1100_0000_0000_0010: // C.SWSP
		log.Println("C.SWSP")
	case i&0b_1110_0000_0000_0011 == 0b_1110_0000_0000_0010: // C.SDSP
		log.Println("C.SDSP")
	}

	return 0
}
