package riscv

import (
	"encoding/binary"
	"log"
)

func ExecuterRV32I(c *CPU, i uint64) (int, error) {
	m := c.Memory
	switch {
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_0111: // LUI
		rd, imm := UType(i)
		imm = SignExtend(imm, 31)
		DebuglnUType("LUI", rd, imm)
		c.Register[rd] = imm
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_00010_111: // AUIPC
		rd, imm := UType(i)
		DebuglnUType("AUIPC", rd, imm)
		c.Register[rd] = c.PC + imm
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_1111: // JAL
		rd, imm := JType(i)
		imm = SignExtend(imm, 19)
		DebuglnJType("JAL", rd, imm)
		c.Register[rd] = c.PC + 4
		c.PC += imm
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_0111: // JALR
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("JALR", rd, rs1, imm)
		c.Register[rd] = c.PC + 4
		c.PC = ((c.Register[rs1] + imm) >> 1) << 1
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_0011: // BEQ
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BEQ", rs1, rs2, imm)
		if c.Register[rs1] == c.Register[rs2] {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0110_0011: // BNE
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BNE", rs1, rs2, imm)
		if c.Register[rs1] != c.Register[rs2] {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0110_0011: // BLT
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BLT", rs1, rs2, imm)
		if int64(c.Register[rs1]) < int64(c.Register[rs2]) {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0110_0011: // BGE
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BGE", rs1, rs2, imm)
		if int64(c.Register[rs1]) >= int64(c.Register[rs2]) {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0110_0011: // BLTU
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BLTU", rs1, rs2, imm)
		if c.Register[rs1] < c.Register[rs2] {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0110_0011: // BGEU
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BGEU", rs1, rs2, imm)
		if c.Register[rs1] >= c.Register[rs2] {
			c.PC += imm
		} else {
			c.PC += 4
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0000_0011: // LB
		rd, rs1, imm := IType(i)
		DebuglnIType("LB", rd, rs1, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		v := SignExtend(uint64(m[a]), 7)
		c.Register[rd] = v
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0000_0011: // LH
		rd, rs1, imm := IType(i)
		DebuglnIType("LH", rd, rs1, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		v := SignExtend(binary.LittleEndian.Uint64(m[a:a+2]), 15)
		c.Register[rd] = v
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0000_0011: // LW
		rd, rs1, imm := IType(i)
		DebuglnIType("LW", rd, rs1, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		v := SignExtend(uint64(binary.LittleEndian.Uint32(m[a:a+4])), 63)
		c.Register[rd] = v
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0000_0011: // LBU
		rd, rs1, imm := IType(i)
		DebuglnIType("LBU", rd, rs1, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		v := uint64(m[a])
		c.Register[rd] = v
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0000_0011: // LHU
		rd, rs1, imm := IType(i)
		DebuglnIType("LHU", rd, rs1, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		v := binary.LittleEndian.Uint64(m[a : a+2])
		c.Register[rd] = v
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0010_0011: // SB
		rs1, rs2, imm := SType(i)
		DebuglnIType("SB", rs1, rs2, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		m[a] = byte(c.Register[rs2])
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0010_0011: // SH
		rs1, rs2, imm := SType(i)
		DebuglnIType("SH", rs1, rs2, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		m[a] = byte(c.Register[rs2])
		m[a+1] = byte(c.Register[rs2] >> 8)
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0010_0011: // SW
		rs1, rs2, imm := SType(i)
		DebuglnIType("SW", rs1, rs2, imm)
		a := c.Register[rs1] + SignExtend(imm, 11)
		m[a] = byte(c.Register[rs2])
		m[a+1] = byte(c.Register[rs2] >> 8)
		m[a+2] = byte(c.Register[rs2] >> 16)
		m[a+3] = byte(c.Register[rs2] >> 24)
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_0011: // ADDI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] + imm
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0001_0011: // SLTI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SLTI", rd, rs1, imm)
		if int64(c.Register[rs1]) < int64(imm) {
			c.Register[rd] = 1
		} else {
			c.Register[rd] = 0
		}
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0001_0011: // SLTIU
		rd, rs1, imm := IType(i)
		DebuglnIType("SLTIU", rd, rs1, imm)
		if c.Register[rs1] < imm {
			c.Register[rd] = 1
		} else {
			c.Register[rd] = 0
		}
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0001_0011: // XORI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("XORI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] ^ imm
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0001_0011: // ORI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ORI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] | imm
		c.PC += 4
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0001_0011: // ANDI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ANDI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] & imm
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_0011: // SLLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SLLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] << imm
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_0011: // SRLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SRLI", rd, rs1, imm)
		c.Register[rd] = c.Register[rs1] >> imm
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_0011: // SRAI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SRAI", rd, rs1, imm)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> imm)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_0011: // ADD
		rd, rs1, rs2 := RType(i)
		DebuglnRType("ADD", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] + c.Register[rs2]
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0000_0000_0011_0011: // SUB
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SUB", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] - c.Register[rs2]
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0011_0011: // SLL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLL", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] << InstructionPart(c.Register[rs2], 0, 4)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0011_0011: // SLT
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLT", rd, rs1, rs2)
		if int64(c.Register[rs1]) < int64(c.Register[rs2]) {
			c.Register[rd] = 1
		} else {
			c.Register[rd] = 0
		}
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0011_0011: // SLTU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLTU", rd, rs1, rs2)
		if c.Register[rs1] < c.Register[rs2] {
			c.Register[rd] = 1
		} else {
			c.Register[rd] = 0
		}
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0011_0011: // XOR
		rd, rs1, rs2 := RType(i)
		DebuglnRType("XOR", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] ^ c.Register[rs2]
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0011_0011: // SRL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRL", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] >> InstructionPart(c.Register[rs2], 0, 4)
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0011_0011: // SRA
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRA", rd, rs1, rs2)
		c.Register[rd] = uint64(int64(c.Register[rs1]) >> InstructionPart(c.Register[rs2], 0, 4))
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0011_0011: // OR
		rd, rs1, rs2 := RType(i)
		DebuglnRType("XOR", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] | c.Register[rs2]
		c.PC += 4
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0011_0011: // AND
		rd, rs1, rs2 := RType(i)
		DebuglnRType("XOR", rd, rs1, rs2)
		c.Register[rd] = c.Register[rs1] & c.Register[rs2]
		c.PC += 4
		return 1, nil
	case i&0b_1111_0000_0000_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0000_0000_0000_1111: // FENCE
		log.Println("FENCE")
	case i&0b_1111_1111_1111_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0001_0000_0000_1111: // FENCE.I
		log.Println("FENCE.I")
	case i&0b_1111_1111_1111_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0000_0000_0111_0011: // ECALL
		log.Println("ECALL")
	case i&0b_1111_1111_1111_1111_1111_1111_1111_1111 == 0b_0000_0000_0001_0000_0000_0000_0111_0011: // EBREAK
		log.Println("EBREAK")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0111_0011: // CSRRW
		log.Println("CSRRW")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0111_0011: // CSRRS
		log.Println("CSRRS")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0111_0011: // CSRRC
		log.Println("CSRRC")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0111_0011: // CSRRWI
		log.Println("CSRRWI")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0111_0011: // CSRRSI
		log.Println("CSRRSI")
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0111_0011: // CSRRCI
		log.Println("CSRRCI")
	}
	return 0, nil
}
