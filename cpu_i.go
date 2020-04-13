package rv64

import (
	"fmt"
	"math"
	"math/big"
)

type isaI struct{}

func (_ *isaI) lui(c *CPU, i uint64) (uint64, error) {
	rd, imm := UType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "lui", c.LogI(rd), imm))
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) aupic(c *CPU, i uint64) (uint64, error) {
	rd, imm := UType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "auipc", c.LogI(rd), imm))
	c.SetRegister(rd, c.GetPC()+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) jal(c *CPU, i uint64) (uint64, error) {
	rd, imm := JType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "jal", c.LogI(rd), imm))
	c.SetRegister(rd, c.GetPC()+4)
	r := c.GetPC() + imm
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaI) jalr(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "jalr", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetPC()+4)
	r := c.GetRegister(rs1) + imm
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaI) beq(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "beq", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) == c.GetRegister(rs2) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) bne(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "bne", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) != c.GetRegister(rs2) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) blt(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "blt", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) bge(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "bge", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if int64(c.GetRegister(rs1)) >= int64(c.GetRegister(rs2)) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) bltu(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "bltu", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) < c.GetRegister(rs2) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) bgeu(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "bgeu", c.LogI(rs1), c.LogI(rs2), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) >= c.GetRegister(rs2) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 4)
	}
	return 1, nil
}

func (_ *isaI) lb(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lb", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint8(a)
	if err != nil {
		return 0, err
	}
	v := SignExtend(uint64(b), 7)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lh(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lh", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint16(a)
	if err != nil {
		return 0, err
	}
	v := SignExtend(uint64(b), 15)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lw", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	v := SignExtend(uint64(b), 31)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) ld(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "ld", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	v := b
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lbu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lbu", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint8(a)
	if err != nil {
		return 0, err
	}
	v := uint64(b)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lhu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lhu", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint16(a)
	if err != nil {
		return 0, err
	}
	v := uint64(b)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lwu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "lwu", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	v := uint64(b)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sb(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "sb", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sh(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "sh", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sw(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "sw", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sd(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "sd", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegister(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) addi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "addi", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slti(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "slti", c.LogI(rd), c.LogI(rs1), imm))
	if int64(c.GetRegister(rs1)) < int64(imm) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sltiu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "sltiu", c.LogI(rd), c.LogI(rs1), imm))
	if c.GetRegister(rs1) < imm {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) xori(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "xori", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)^imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) ori(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "ori", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)|imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) andi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "andi", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)&imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slli(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	shamt := imm & 0x3f
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "slli", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)<<shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srli(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "srli", c.LogI(rd), c.LogI(rs1), imm))
	shamt := imm & 0x3f
	c.SetRegister(rd, c.GetRegister(rs1)>>shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srai(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "srai", c.LogI(rd), c.LogI(rs1), imm))
	shamt := imm & 0x3f
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>shamt))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) add(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "add", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)+c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sub(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sub", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)-c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sll(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sll", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)<<(c.GetRegister(rs2)&0x3f))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slt(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "slt", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sltu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sltu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs1) < c.GetRegister(rs2) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) xor(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "xor", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)^c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srl(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "srl", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)>>(c.GetRegister(rs2)&0x3f))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}
func (_ *isaI) sra(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sra", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>(c.GetRegister(rs2)&0x3f)))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) or(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "or", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)|c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) and(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "and", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs1)&c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) fence(c *CPU, _ uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "fence"))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) ecall(c *CPU, _ uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "ecall"))
	return c.GetSystem().HandleCall(c)
}

func (_ *isaI) ebreak(c *CPU, _ uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "ebreak"))
	return 1, nil
}

func (_ *isaI) addiw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "addiw", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(imm)))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slliw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "slliw", c.LogI(rd), c.LogI(rs1), imm))
	if InstructionPart(imm, 5, 5) != 0x00 {
		return 0, ErrAbnormalInstruction
	}
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<imm), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srliw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "srliw", c.LogI(rd), c.LogI(rs1), imm))
	if InstructionPart(imm, 5, 5) != 0x00 {
		return 0, ErrAbnormalInstruction
	}
	shamt := InstructionPart(imm, 0, 4)
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>shamt), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sraiw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "sraiw", c.LogI(rd), c.LogI(rs1), imm))
	if InstructionPart(imm, 5, 5) != 0x00 {
		return 0, ErrAbnormalInstruction
	}
	shamt := InstructionPart(imm, 0, 4)
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>shamt))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) addw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "addw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) subw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "subw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))-int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sllw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sllw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	s := c.GetRegister(rs2) & 0x1f
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<s), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srlw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "srlw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	s := c.GetRegister(rs2) & 0x1f
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>s), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sraw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sraw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaZifencei struct{}

func (_ *isaZifencei) fencei(c *CPU, i uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "fence.i"))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaZicsr struct{}

func (_ *isaZicsr) csrrw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrw", c.LogI(rd), c.LogI(rs1), csr))
	if rd != Rzero {
		c.SetRegister(rd, c.GetCSR().Get(csr))
	}
	c.GetCSR().Set(csr, c.GetRegister(rs1))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaZicsr) csrrs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrs", c.LogI(rd), c.LogI(rs1), csr))
	c.SetRegister(rd, c.GetCSR().Get(csr))
	if rs1 != Rzero {
		c.GetCSR().Set(csr, c.GetCSR().Get(csr)|c.GetRegister(rs1))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaZicsr) csrrc(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrc", c.LogI(rd), c.LogI(rs1), csr))
	c.SetRegister(rd, c.GetCSR().Get(csr))
	if rs1 != Rzero {
		c.GetCSR().Set(csr, c.GetCSR().Get(csr)&(math.MaxUint64-c.GetRegister(rs1)))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaZicsr) csrrwi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrwi", c.LogI(rd), c.LogI(rs1), csr))
	if rd != Rzero {
		c.SetRegister(rd, c.GetCSR().Get(csr))
	}
	c.GetCSR().Set(csr, rs1)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaZicsr) csrrsi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrsi", c.LogI(rd), c.LogI(rs1), csr))
	c.SetRegister(rd, c.GetCSR().Get(csr))
	if csr != 0x00 {
		c.GetCSR().Set(csr, c.GetCSR().Get(csr)|rs1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaZicsr) csrrci(c *CPU, i uint64) (uint64, error) {
	rd, rs1, csr := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s csr: %#016x", c.GetPC(), "csrrci", c.LogI(rd), c.LogI(rs1), csr))
	c.SetRegister(rd, c.GetCSR().Get(csr))
	if csr != 0x00 {
		c.GetCSR().Set(csr, c.GetCSR().Get(csr)&(math.MaxUint64-rs1))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaM struct{}

func (_ *isaM) mul(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "mul", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulh(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "mulh", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	v := func() uint64 {
		ag1 := big.NewInt(int64(c.GetRegister(rs1)))
		ag2 := big.NewInt(int64(c.GetRegister(rs2)))
		tmp := big.NewInt(0)
		tmp.Mul(ag1, ag2)
		tmp.Rsh(tmp, 64)
		return uint64(tmp.Int64())
	}()
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulhsu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "mulhsu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	v := func() uint64 {
		ag1 := big.NewInt(int64(c.GetRegister(rs1)))
		ag2 := big.NewInt(int64(c.GetRegister(rs2)))
		if ag2.Cmp(big.NewInt(0)) == -1 {
			tmp := big.NewInt(0)
			tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
			tmp.Add(tmp, big.NewInt(2))
			ag2 = tmp.Add(tmp, ag2)
		}
		tmp := big.NewInt(0)
		tmp.Mul(ag1, ag2)
		tmp.Rsh(tmp, 64)
		return uint64(tmp.Int64())
	}()
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulhu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "mulhu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	v := func() uint64 {
		ag1 := big.NewInt(int64(c.GetRegister(rs1)))
		ag2 := big.NewInt(int64(c.GetRegister(rs2)))
		if ag1.Cmp(big.NewInt(0)) == -1 {
			tmp := big.NewInt(0)
			tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
			tmp.Add(tmp, big.NewInt(2))
			ag1 = tmp.Add(tmp, ag1)
		}
		if ag2.Cmp(big.NewInt(0)) == -1 {
			tmp := big.NewInt(0)
			tmp.Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))
			tmp.Add(tmp, big.NewInt(2))
			ag2 = tmp.Add(tmp, ag2)
		}
		tmp := big.NewInt(0)
		tmp.Mul(ag1, ag2)
		tmp.Rsh(tmp, 64)
		return tmp.Uint64()
	}()
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) div(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "div", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))/int64(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "divu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) rem(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "rem", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))%int64(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) remu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "remu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, c.GetRegister(rs1)%c.GetRegister(rs2))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "mulw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "divw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divuw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "divuw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) remw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "remw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))%int32(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}
func (_ *isaM) remuw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "remuw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))%uint32(c.GetRegister(rs2))), 31))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaA struct{}

func (_ *isaA) lrw()      {}
func (_ *isaA) scw()      {}
func (_ *isaA) amoswapw() {}
func (_ *isaA) amoaddw()  {}
func (_ *isaA) amoxorw()  {}
func (_ *isaA) amoandw()  {}
func (_ *isaA) amoorw()   {}
func (_ *isaA) amominw()  {}
func (_ *isaA) amomaxw()  {}
func (_ *isaA) amominuw() {}
func (_ *isaA) amomaxuw() {}
func (_ *isaA) lrd()      {}
func (_ *isaA) scd()      {}
func (_ *isaA) amoswapd() {}
func (_ *isaA) amoaddd()  {}
func (_ *isaA) amoxord()  {}
func (_ *isaA) amoandd()  {}
func (_ *isaA) amoord()   {}
func (_ *isaA) amomind()  {}
func (_ *isaA) amomaxd()  {}
func (_ *isaA) amominud() {}
func (_ *isaA) amomaxud() {}

type isaF struct{}

func (_ *isaF) flw()     {}
func (_ *isaF) fsw()     {}
func (_ *isaF) fmadds()  {}
func (_ *isaF) fmsubs()  {}
func (_ *isaF) fnmsubs() {}
func (_ *isaF) fnmadds() {}
func (_ *isaF) fadds()   {}
func (_ *isaF) fsubs()   {}
func (_ *isaF) fmuls()   {}
func (_ *isaF) fdivs()   {}
func (_ *isaF) fsqrts()  {}
func (_ *isaF) fsgnjs()  {}
func (_ *isaF) fsgnjns() {}
func (_ *isaF) fsgnjxs() {}
func (_ *isaF) fmins()   {}
func (_ *isaF) fmaxs()   {}
func (_ *isaF) fcvtws()  {}
func (_ *isaF) fcvtwus() {}
func (_ *isaF) fmvxw()   {}
func (_ *isaF) feqs()    {}
func (_ *isaF) flts()    {}
func (_ *isaF) fles()    {}
func (_ *isaF) fclasss() {}
func (_ *isaF) fcvtsw()  {}
func (_ *isaF) fcvtswu() {}
func (_ *isaF) fmvwx()   {}
func (_ *isaF) fcvtls()  {}
func (_ *isaF) fcvtlus() {}
func (_ *isaF) fcvtsl()  {}
func (_ *isaF) fcvtslu() {}

type isaD struct{}

func (_ *isaD) fld()     {}
func (_ *isaD) fsd()     {}
func (_ *isaD) fmaddd()  {}
func (_ *isaD) fmsubd()  {}
func (_ *isaD) fnmsubd() {}
func (_ *isaD) fnmaddd() {}
func (_ *isaD) faddd()   {}
func (_ *isaD) fsubd()   {}
func (_ *isaD) fmuld()   {}
func (_ *isaD) fdivd()   {}
func (_ *isaD) fsqrtd()  {}
func (_ *isaD) fsgnjd()  {}
func (_ *isaD) fsgnjnd() {}
func (_ *isaD) fsgnjxd() {}
func (_ *isaD) fmind()   {}
func (_ *isaD) fmaxd()   {}
func (_ *isaD) fcvtsd()  {}
func (_ *isaD) fcvtds()  {}
func (_ *isaD) feqd()    {}
func (_ *isaD) fltd()    {}
func (_ *isaD) fled()    {}
func (_ *isaD) fclassd() {}
func (_ *isaD) fcvtwd()  {}
func (_ *isaD) fcvtwud() {}
func (_ *isaD) fcvtdw()  {}
func (_ *isaD) fcvtdwu() {}
func (_ *isaD) fcvtld()  {}
func (_ *isaD) fcvtlud() {}
func (_ *isaD) fmvxd()   {}
func (_ *isaD) fcvtdl()  {}
func (_ *isaD) fcvtdlu() {}
func (_ *isaD) fmvdx()   {}

type isaC struct{}

func (_ *isaC) addi4spn(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 7, 10)<<6 | InstructionPart(i, 11, 12)<<4 | InstructionPart(i, 5, 5)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.addi4spn", c.LogI(rd), imm))
	if imm == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, c.GetRegister(Rsp)+imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fld() {}

func (_ *isaC) lw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = InstructionPart(i, 5, 5)<<6 | InstructionPart(i, 10, 12)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "c.lw", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	v := SignExtend(uint64(b), 31)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) ld(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "c.ld", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, b)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fsd() {}

func (_ *isaC) sw(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 5, 5)<<6 | InstructionPart(i, 10, 12)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "c.sw", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) sd(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "c.sd", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegister(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) nop(c *CPU, _ uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "c.nop"))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (z *isaC) addi(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	if rd != Rzero && imm != 0 {
		Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.addi", c.LogI(rd), imm))
		c.SetRegister(rd, c.GetRegister(rd)+imm)
		c.SetPC(c.GetPC() + 2)
		return 1, nil
	} else if rd == Rzero {
		return z.nop(c, i)
	} else {
		return 0, ErrHint
	}
}

func (_ *isaC) addiw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	if rd == 0x00 {
		return 0, ErrReservedInstruction
	}
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.addiw", c.LogI(rd), imm))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rd))+int32(imm)))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) li(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.li", c.LogI(rd), imm))
	if rd == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) addi16sp(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<9|InstructionPart(i, 3, 4)<<7|InstructionPart(i, 5, 5)<<6|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 6, 6)<<4, 9)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.addi16sp", c.LogI(rd), imm))
	if imm == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, c.GetRegister(Rsp)+imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) lui(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<17|InstructionPart(i, 2, 6)<<12, 17)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.lui", c.LogI(rd), imm))
	if rd == 0x00 || rd == 0x02 || imm == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) srli(c *CPU, i uint64) (uint64, error) {
	var (
		rd    = InstructionPart(i, 7, 9) + 8
		shamt = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.srli", c.LogI(rd), shamt))
	c.SetRegister(rd, c.GetRegister(rd)>>shamt)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) srai(c *CPU, i uint64) (uint64, error) {
	var (
		rd    = InstructionPart(i, 7, 9) + 8
		shamt = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.srai", c.LogI(rd), shamt))
	c.SetRegister(rd, uint64(int64(c.GetRegister(rd))>>shamt))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) andi(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.andi", c.LogI(rd), imm))
	c.SetRegister(rd, c.GetRegister(rd)&imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) sub(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.sub", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rd)-c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) xor(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.xor", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rd)^c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) or(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.or", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rd)|c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) and(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.and", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rd)&c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) subw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.subw", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rd))-int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) addw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.addw", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rd))+int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) j(c *CPU, i uint64) (uint64, error) {
	var imm = InstructionPart(i, 12, 12)<<11 |
		InstructionPart(i, 8, 8)<<10 |
		InstructionPart(i, 9, 10)<<8 |
		InstructionPart(i, 6, 6)<<7 |
		InstructionPart(i, 7, 7)<<6 |
		InstructionPart(i, 2, 2)<<5 |
		InstructionPart(i, 11, 11)<<4 |
		InstructionPart(i, 3, 5)<<1
	Debugln(fmt.Sprintf("%#08x % 10s imm: ____(%#016x)", c.GetPC(), "c.j", imm))
	r := c.GetPC() + imm
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaC) beqz(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = SignExtend(InstructionPart(i, 3, 4)<<1|InstructionPart(i, 10, 11)<<3|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 5, 6)<<6|InstructionPart(i, 12, 12)<<8, 8)
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s imm: ____(%#016x)", c.GetPC(), "c.beqz", c.LogI(rs1), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) == c.GetRegister(Rzero) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 2)
	}
	return 1, nil
}

func (_ *isaC) bnez(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = SignExtend(InstructionPart(i, 3, 4)<<1|InstructionPart(i, 10, 11)<<3|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 5, 6)<<6|InstructionPart(i, 12, 12)<<8, 8)
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s imm: ____(%#016x)", c.GetPC(), "c.bnez", c.LogI(rs1), imm))
	if imm%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	if c.GetRegister(rs1) != c.GetRegister(Rzero) {
		c.SetPC(c.GetPC() + imm)
	} else {
		c.SetPC(c.GetPC() + 2)
	}
	return 1, nil
}

func (_ *isaC) slli(c *CPU, i uint64) (uint64, error) {
	var (
		rd    = InstructionPart(i, 7, 9) + 8
		shamt = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ____(%#016x)", c.GetPC(), "c.slli", c.LogI(rd), shamt))
	if rd == 0 {
		return 0, ErrAbnormalInstruction
	}
	c.SetRegister(rd, c.GetRegister(rd)<<shamt)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fldsp() {}
func (_ *isaC) lwsp()  {}
func (_ *isaC) ldsp()  {}

func (_ *isaC) jr(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 11)
		rs2 = InstructionPart(i, 2, 6)
	)
	if rs1 == 0 {
		return 0, ErrReservedInstruction
	}
	if rs2 != 0 {
		return 0, ErrAbnormalInstruction
	}
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s", c.GetPC(), "c.jr", c.LogI(rs1)))
	r := c.GetRegister(rs1)
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaC) mv(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		rs2 = InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.mv", c.LogI(rd), c.LogI(rs2)))
	c.SetRegister(rd, c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) ebreak() {}

func (_ *isaC) jalr(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 11)
	)
	if rs1 == 0 {
		return 0, ErrReservedInstruction
	}
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s", c.GetPC(), "c.jalr", c.LogI(rs1)))
	c.SetRegister(Rra, c.GetPC()+2)
	r := c.GetRegister(rs1)
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaC) add()   {}
func (_ *isaC) fsdsp() {}
func (_ *isaC) sqsp()  {}
func (_ *isaC) swsp()  {}
func (_ *isaC) sdsp()  {}

var (
	aluI        = &isaI{}
	aluZifencei = &isaZifencei{}
	aluZicsr    = &isaZicsr{}
	aluM        = &isaM{}
	aluA        = &isaA{}
	aluF        = &isaF{}
	aluD        = &isaD{}
	aluC        = &isaC{}
)
