package riscv

import (
	"fmt"
	"log"
)

func ExecuterRV64I(c *CPU, i uint64) (int, error) {
	switch {
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_0111: // LUI
		rd, imm := UType(i)
		imm = SignExtend(imm, 31)
		DebuglnUType("LUI", rd, imm)
		c.SetRegister(rd, imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_00010_111: // AUIPC
		rd, imm := UType(i)
		imm = SignExtend(imm, 31)
		DebuglnUType("AUIPC", rd, imm)
		c.SetRegister(rd, c.GetPC()+imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0000_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_1111: // JAL
		rd, imm := JType(i)
		imm = SignExtend(imm, 19)
		DebuglnJType("JAL", rd, imm)
		c.SetRegister(rd, c.GetPC()+4)
		c.SetPC(c.GetPC() + imm)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_0111: // JALR
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("JALR", rd, rs1, imm)
		c.SetRegister(rd, c.GetPC()+4)
		c.SetPC(((c.GetRegister(rs1) + imm) >> 1) << 1)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0110_0011: // BEQ
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BEQ", rs1, rs2, imm)
		if c.GetRegister(rs1) == c.GetRegister(rs2) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0110_0011: // BNE
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BNE", rs1, rs2, imm)
		if c.GetRegister(rs1) != c.GetRegister(rs2) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0110_0011: // BLT
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BLT", rs1, rs2, imm)
		if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0110_0011: // BGE
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BGE", rs1, rs2, imm)
		if int64(c.GetRegister(rs1)) >= int64(c.GetRegister(rs2)) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0110_0011: // BLTU
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BLTU", rs1, rs2, imm)
		if c.GetRegister(rs1) < c.GetRegister(rs2) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0110_0011: // BGEU
		rs1, rs2, imm := BType(i)
		imm = SignExtend(imm, 12)
		DebuglnBType("BGEU", rs1, rs2, imm)
		if c.GetRegister(rs1) >= c.GetRegister(rs2) {
			c.SetPC(c.GetPC() + imm)
		} else {
			c.SetPC(c.GetPC() + 4)
		}
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0000_0011: // LB
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LB", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		mem, err := c.GetMemory().Get(a, 1)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, SignExtend(uint64(mem[0]), 7))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0000_0011: // LH
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LH", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint16(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, SignExtend(uint64(v), 15))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0000_0011: // LW
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LW", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, SignExtend(uint64(v), 63))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0000_0011: // LBU
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LBU", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint8(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, uint64(v))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0000_0011: // LHU
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LH", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint16(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, uint64(v))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0010_0011: // SB
		rs1, rs2, imm := SType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SB", rs1, rs2, imm)
		a := c.GetRegister(rs1) + imm
		err := c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2)))
		if err != nil {
			return 0, err
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0010_0011: // SH
		rs1, rs2, imm := SType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SH", rs1, rs2, imm)
		a := c.GetRegister(rs1) + imm
		err := c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2)))
		if err != nil {
			return 0, err
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0010_0011: // SW
		rs1, rs2, imm := SType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SW", rs1, rs2, imm)
		a := c.GetRegister(rs1) + imm
		err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
		if err != nil {
			return 0, err
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_0011: // ADDI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)+imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0001_0011: // SLTI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SLTI", rd, rs1, imm)
		if int64(c.GetRegister(rs1)) < int64(imm) {
			c.SetRegister(rd, 1)
		} else {
			c.SetRegister(rd, 0)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0001_0011: // SLTIU
		rd, rs1, imm := IType(i)
		DebuglnIType("SLTIU", rd, rs1, imm)
		if c.GetRegister(rs1) < imm {
			c.SetRegister(rd, 1)
		} else {
			c.SetRegister(rd, 0)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0001_0011: // XORI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("XORI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)^imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0001_0011: // ORI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ORI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)|imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0001_0011: // ANDI
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ANDI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)&imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_0011: // SLLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SLLI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)<<imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_0011: // SRLI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRLI", rd, rs1, imm)
		c.SetRegister(rd, c.GetRegister(rs1)>>imm)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_0011: // SRAI
		rd, rs1, imm := IType(i)
		imm = InstructionPart(imm, 0, 5)
		DebuglnIType("SRAI", rd, rs1, imm)
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>imm))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_0011: // ADD
		rd, rs1, rs2 := RType(i)
		DebuglnRType("ADD", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)+c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0000_0000_0011_0011: // SUB
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SUB", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)-c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0011_0011: // SLL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLL", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)<<InstructionPart(c.GetRegister(rs2), 0, 5))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0010_0000_0011_0011: // SLT
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLT", rd, rs1, rs2)
		if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
			c.SetRegister(rd, 1)
		} else {
			c.SetRegister(rd, 0)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0011_0011: // SLTU
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLTU", rd, rs1, rs2)
		if c.GetRegister(rs1) < c.GetRegister(rs2) {
			c.SetRegister(rd, 1)
		} else {
			c.SetRegister(rd, 0)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0100_0000_0011_0011: // XOR
		rd, rs1, rs2 := RType(i)
		DebuglnRType("XOR", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)^c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0011_0011: // SRL
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRL", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)>>InstructionPart(c.GetRegister(rs2), 0, 5))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0011_0011: // SRA
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRA", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 5)))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0011_0011: // OR
		rd, rs1, rs2 := RType(i)
		DebuglnRType("OR", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)|c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0111_0000_0011_0011: // AND
		rd, rs1, rs2 := RType(i)
		DebuglnRType("AND", rd, rs1, rs2)
		c.SetRegister(rd, c.GetRegister(rs1)&c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_0000_0000_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0000_0000_0000_1111: // FENCE
		Debugln(fmt.Sprintf("Instr: % 10s |", "FENCE"))
		return 1, nil
	case i&0b_1111_1111_1111_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0001_0000_0000_1111: // FENCE.I
		Debugln(fmt.Sprintf("Instr: % 10s |", "FENCE.I"))
		return 1, nil
	case i&0b_1111_1111_1111_1111_1111_1111_1111_1111 == 0b_0000_0000_0000_0000_0000_0000_0111_0011: // ECALL
		rd, rs1, imm := IType(i)
		DebuglnIType("ECALL", rd, rs1, imm)
		return c.GetSystem().HandleCall(c)
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

	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0110_0000_0000_0011: // LWU
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LWU", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint32(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, uint64(v))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0000_0011: // LD
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("LD", rd, rs1, imm)
		a := c.GetRegister(rs1) + imm
		v, err := c.GetMemory().GetUint64(a)
		if err != nil {
			return 0, err
		}
		c.SetRegister(rd, v)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0011_0000_0010_0011: // SD
		rs1, rs2, imm := SType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("SD", rs1, rs2, imm)
		a := c.GetRegister(rs1) + imm
		c.GetMemory().SetUint64(a, c.GetRegister(rs2))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_0000_0000_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0001_1011: // ADDIW
		rd, rs1, imm := IType(i)
		imm = SignExtend(imm, 11)
		DebuglnIType("ADDIW", rd, rs1, imm)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(imm)))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0001_1011: // SLLIW
		rd, rs1, imm := IType(i)
		if InstructionPart(imm, 5, 5) != 0x00 {
			return 0, ErrAbnormalInstruction
		}
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SLLIW", rd, rs1, imm)
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<imm), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0001_1011: // SRLIW
		rd, rs1, imm := IType(i)
		if InstructionPart(imm, 5, 5) != 0x00 {
			return 0, ErrAbnormalInstruction
		}
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SRLIW", rd, rs1, imm)
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>imm), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0001_1011: // SRAIW
		rd, rs1, imm := IType(i)
		if InstructionPart(imm, 5, 5) != 0x00 {
			return 0, ErrAbnormalInstruction
		}
		imm = InstructionPart(imm, 0, 4)
		DebuglnIType("SRAIW", rd, rs1, imm)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>imm))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0000_0000_0011_1011: // ADDW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("ADDW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0000_0000_0011_1011: // SUBW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SUBW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(c.GetRegister(rs2))))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0001_0000_0011_1011: // SLLW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SLLW", rd, rs1, rs2)
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<InstructionPart(c.GetRegister(rs2), 0, 4)), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0000_0000_0000_0000_0101_0000_0011_1011: // SRLW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRLW", rd, rs1, rs2)
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)), 31))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	case i&0b_1111_1110_0000_0000_0111_0000_0111_1111 == 0b_0100_0000_0000_0000_0101_0000_0011_1011: // SRAW
		rd, rs1, rs2 := RType(i)
		DebuglnRType("SRAW", rd, rs1, rs2)
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	return 0, nil
}
