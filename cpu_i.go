package rv64

import (
	"fmt"
	"math"
	"math/big"
)

type isaI struct{}

func (_ *isaI) lui(c *CPU, i uint64) (uint64, error) {
	rd, imm := UType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "lui", c.LogI(rd), imm))
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) aupic(c *CPU, i uint64) (uint64, error) {
	rd, imm := UType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "auipc", c.LogI(rd), imm))
	c.SetRegister(rd, c.GetPC()+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) jal(c *CPU, i uint64) (uint64, error) {
	rd, imm := JType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "jal", c.LogI(rd), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "jalr", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetPC()+4)
	r := c.GetRegister(rs1) + imm
	c.SetPC(r & 0xfffffffffffffffe)
	return 1, nil
}

func (_ *isaI) beq(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := BType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "beq", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "bne", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "blt", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "bge", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "bltu", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "bgeu", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lb", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lh", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lw", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "ld", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, b)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) lbu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lbu", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lhu", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "lwu", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "sb", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sh(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "sh", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sw(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "sw", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sd(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "sd", c.LogI(rs1), c.LogI(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegister(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) addi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "addi", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slti(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "slti", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "sltiu", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "xori", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)^imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) ori(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "ori", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)|imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) andi(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "andi", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)&imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slli(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	shamt := imm & 0x3f
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "slli", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, c.GetRegister(rs1)<<shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srli(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "srli", c.LogI(rd), c.LogI(rs1), imm))
	shamt := imm & 0x3f
	c.SetRegister(rd, c.GetRegister(rs1)>>shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srai(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "srai", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "addiw", c.LogI(rd), c.LogI(rs1), imm))
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(imm)))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slliw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "slliw", c.LogI(rd), c.LogI(rs1), imm))
	if InstructionPart(imm, 5, 5) != 0x00 {
		return 0, ErrAbnormalInstruction
	}
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<imm), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srliw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "srliw", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "sraiw", c.LogI(rd), c.LogI(rs1), imm))
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

func (_ *isaA) lrw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "lr.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetLoadReservation(a)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) scw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sc.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	if a == c.GetLoadReservation() {
		c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
		c.SetRegister(rd, 0)
	} else {
		c.SetRegister(rd, 1)
	}
	c.SetLoadReservation(0)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoswapw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoswap.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoaddw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoadd.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint32(a, v+uint32(c.GetRegister(rs2)))
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoxorw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoxor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint32(a, v^uint32(c.GetRegister(rs2)))
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoandw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoand.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint32(a, v&uint32(c.GetRegister(rs2)))
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoorw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint32(a, v|uint32(c.GetRegister(rs2)))
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amominw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomin.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	var r uint32
	if int32(v) < int32(uint32(c.GetRegister(rs2))) {
		r = v
	} else {
		r = uint32(c.GetRegister(rs2))
	}
	c.GetMemory().SetUint32(a, r)
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amomaxw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomax.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	var r uint32
	if int32(v) > int32(uint32(c.GetRegister(rs2))) {
		r = v
	} else {
		r = uint32(c.GetRegister(rs2))
	}
	c.GetMemory().SetUint32(a, r)
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amominuw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amominu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	var r uint32
	if v < uint32(c.GetRegister(rs2)) {
		r = v
	} else {
		r = uint32(c.GetRegister(rs2))
	}
	c.GetMemory().SetUint32(a, r)
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amomaxuw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomaxu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := SignExtend(c.GetRegister(rs1), 31)
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	var r uint32
	if v > uint32(c.GetRegister(rs2)) {
		r = v
	} else {
		r = uint32(c.GetRegister(rs2))
	}
	c.GetMemory().SetUint32(a, r)
	c.SetRegister(rd, SignExtend(uint64(v), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) lrd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "lr.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, v)
	c.SetLoadReservation(a)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) scd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sc.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	if a == c.GetLoadReservation() {
		c.GetMemory().SetUint64(a, c.GetRegister(rs2))
		c.SetRegister(rd, 0)
	} else {
		c.SetRegister(rd, 1)
	}
	c.SetLoadReservation(0)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoswapd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoswap.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint64(a, c.GetRegister(rs2))
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoaddd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoadd.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint64(a, v+c.GetRegister(rs2))
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoxord(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoxor.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint64(a, v^c.GetRegister(rs2))
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoandd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoand.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint64(a, v&c.GetRegister(rs2))
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amoord(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoor.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.GetMemory().SetUint64(a, v|c.GetRegister(rs2))
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amomind(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomin.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	var r uint64 = 0
	if int64(v) < int64(c.GetRegister(rs2)) {
		r = v
	} else {
		r = c.GetRegister(rs2)
	}
	c.GetMemory().SetUint64(a, r)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amomaxd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomax.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	var r uint64 = 0
	if int64(v) > int64(c.GetRegister(rs2)) {
		r = v
	} else {
		r = c.GetRegister(rs2)
	}
	c.GetMemory().SetUint64(a, r)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amominud(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amominu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	var r uint64 = 0
	if v < c.GetRegister(rs2) {
		r = v
	} else {
		r = c.GetRegister(rs2)
	}
	c.GetMemory().SetUint64(a, r)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaA) amomaxud(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomaxu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
	a := c.GetRegister(rs1)
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	var r uint64 = 0
	if v > c.GetRegister(rs2) {
		r = v
	} else {
		r = c.GetRegister(rs2)
	}
	c.GetMemory().SetUint64(a, r)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaF struct{}

func (_ *isaF) flw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "flw", c.LogF(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	v, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	c.SetRegisterFloatAsFloat32(rd, math.Float32frombits(v))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsw(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "fsw", c.LogI(rs1), c.LogF(rs2), imm))
	a := c.GetRegister(rs1) + imm
	err := c.GetMemory().SetUint32(a, uint32(c.GetRegisterFloat(rs2)))
	if err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmadds(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fmadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	d := c.GetRegisterFloatAsFloat32(rs3)
	r := a*b + d
	c.SetRegisterFloatAsFloat32(rd, r)
	if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmsubs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fmsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	d := c.GetRegisterFloatAsFloat32(rs3)
	r := a*b - d
	c.SetRegisterFloatAsFloat32(rd, r)
	if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fnmsubs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fnmsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	d := c.GetRegisterFloatAsFloat32(rs3)
	r := a*b - d
	c.SetRegisterFloatAsFloat32(rd, -r)
	if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fnmadds(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fnmadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	d := c.GetRegisterFloatAsFloat32(rs3)
	r := a*b + d
	c.SetRegisterFloatAsFloat32(rd, -r)
	if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fadds(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	d := a + b
	c.SetRegisterFloatAsFloat32(rd, d)
	if d-a != b || d-b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsubs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	if (math.Signbit(float64(a)) == math.Signbit(float64(b))) && math.IsInf(float64(a), 0) && math.IsInf(float64(b), 0) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	d := a - b
	c.SetRegisterFloatAsFloat32(rd, d)
	if a-d != b || b+d != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmuls(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmul.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	d := a * b
	c.SetRegisterFloatAsFloat32(rd, d)
	if d/a != b || d/b != a || float64(a)*float64(b) != float64(d) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fdivs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fdiv.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	if b == 0 {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
		c.SetFloatFlag(FFlagsDZ, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	d := a / b
	c.SetRegisterFloatAsFloat32(rd, d)
	if a/d != b || b*d != a || float64(b)*float64(d) != float64(a) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsqrts(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsqrt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	c.ClrFloatFlag()
	if a < 0 {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	d := float32(math.Sqrt(float64(a)))
	c.SetRegisterFloatAsFloat32(rd, d)
	if a/d != d || d*d != a || float64(d)*float64(d) != float64(a) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsgnjs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnj.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	if math.Signbit(float64(b)) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
	} else {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsgnjns(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjn.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	if math.Signbit(float64(b)) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
	} else {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fsgnjxs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjx.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	if math.Signbit(float64(a)) != math.Signbit(float64(b)) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
	} else {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmins(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmin.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	if math.IsNaN(float64(a)) && math.IsNaN(float64(b)) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if math.IsNaN(float64(a)) {
		c.SetRegisterFloatAsFloat32(rd, b)
		if IsSNaN32(a) {
			c.SetFloatFlag(FFlagsNV, 1)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if math.IsNaN(float64(b)) {
		c.SetRegisterFloatAsFloat32(rd, a)
		if IsSNaN32(b) {
			c.SetFloatFlag(FFlagsNV, 1)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if (math.Signbit(float64(a)) && !math.Signbit(float64(b))) || a < b {
		c.SetRegisterFloatAsFloat32(rd, a)
	} else {
		c.SetRegisterFloatAsFloat32(rd, b)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmaxs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmax.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	c.ClrFloatFlag()
	if math.IsNaN(float64(a)) && math.IsNaN(float64(b)) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if math.IsNaN(float64(a)) {
		c.SetRegisterFloatAsFloat32(rd, b)
		if IsSNaN32(a) {
			c.SetFloatFlag(FFlagsNV, 1)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if math.IsNaN(float64(b)) {
		c.SetRegisterFloatAsFloat32(rd, a)
		if IsSNaN32(b) {
			c.SetFloatFlag(FFlagsNV, 1)
		}
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if (!math.Signbit(float64(a)) && math.Signbit(float64(b))) || a > b {
		c.SetRegisterFloatAsFloat32(rd, a)
	} else {
		c.SetRegisterFloatAsFloat32(rd, b)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtws(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.w.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	d := c.GetRegisterFloatAsFloat32(rs1)
	if math.IsNaN(float64(d)) {
		c.SetRegister(rd, 0x7fffffff)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d > float32(math.MaxInt32) {
		c.SetRegister(rd, SignExtend(0x7fffffff, 31))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d < float32(math.MinInt32) {
		c.SetRegister(rd, SignExtend(0x80000000, 31))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	c.SetRegister(rd, SignExtend(uint64(int32(d)), 31))
	if math.Ceil(float64(d)) != float64(d) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtwus(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.wu.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	d := c.GetRegisterFloatAsFloat32(rs1)
	if math.IsNaN(float64(d)) {
		c.SetRegister(rd, 0xffffffffffffffff)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d > float32(math.MaxUint32) {
		c.SetRegister(rd, SignExtend(0xffffffff, 31))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d <= float32(-1) {
		c.SetRegister(rd, SignExtend(0x00000000, 31))
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	c.SetRegister(rd, SignExtend(uint64(uint32(d)), 31))
	if math.Ceil(float64(d)) != float64(d) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmvxw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.x.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegisterFloat(rs1))), 31))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) feqs(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "feq.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	var cond bool
	if IsSNaN32(a) || IsSNaN32(b) {
		c.SetFloatFlag(FFlagsNV, 1)
	} else {
		cond = a == b
	}
	if cond {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) flts(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "flt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	var cond bool
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		c.SetFloatFlag(FFlagsNV, 1)
	} else {
		cond = a < b
	}
	if cond {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fles(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fle.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	b := c.GetRegisterFloatAsFloat32(rs2)
	var cond bool
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		c.SetFloatFlag(FFlagsNV, 1)
	} else {
		cond = a <= b
	}
	if cond {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fclasss(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fclass.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	a := c.GetRegisterFloatAsFloat32(rs1)
	c.SetRegister(rd, FClassS(a))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtsw(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegisterFloatAsFloat32(rd, float32(int32(c.GetRegister(rs1))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtswu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.wu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegisterFloatAsFloat32(rd, float32(uint32(c.GetRegister(rs1))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fmvwx(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.w.x", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(uint32(c.GetRegister(rs1))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtls(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.l.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	d := c.GetRegisterFloatAsFloat32(rs1)
	if math.IsNaN(float64(d)) {
		c.SetRegister(rd, 0x7fffffffffffffff)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d > float32(math.MaxInt64) {
		c.SetRegister(rd, 0x7fffffffffffffff)
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d < float32(math.MinInt64) {
		c.SetRegister(rd, 0x8000000000000000)
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	c.SetRegister(rd, uint64(int64(d)))
	if math.Ceil(float64(d)) != float64(d) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtlus(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.lu.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	d := c.GetRegisterFloatAsFloat32(rs1)
	if math.IsNaN(float64(d)) {
		c.SetRegister(rd, 0xffffffffffffffff)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d > float32(math.MaxUint64) {
		c.SetRegister(rd, 0xffffffffffffffff)
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	if d <= float32(-1) {
		c.SetRegister(rd, 0x0000000000000000)
		c.SetFloatFlag(FFlagsNV, 1)
		c.SetPC(c.GetPC() + 4)
		return 1, nil
	}
	c.SetRegister(rd, uint64(d))
	if math.Ceil(float64(d)) != float64(d) {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtsl(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.l", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegisterFloatAsFloat32(rd, float32(int64(c.GetRegister(rs1))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaF) fcvtslu(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.lu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	c.SetRegisterFloatAsFloat32(rd, float32(uint64(c.GetRegister(rs1))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

type isaD struct{}

func (_ *isaD) fld(c *CPU, i uint64) (uint64, error) {
	rd, rs1, imm := IType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "fld", c.LogF(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegisterFloat(rd, v)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaD) fsd(c *CPU, i uint64) (uint64, error) {
	rs1, rs2, imm := SType(i)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "fsd", c.LogI(rs1), c.LogF(rs2), imm))
	a := c.GetRegister(rs1) + imm
	err := c.GetMemory().SetUint64(a, c.GetRegisterFloat(rs2))
	if err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaD) fmaddd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fmadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat64(rs1)
	b := c.GetRegisterFloatAsFloat64(rs2)
	d := c.GetRegisterFloatAsFloat64(rs3)
	r := a*b + d
	c.SetRegisterFloatAsFloat64(rd, r)
	if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaD) fmsubd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fmsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat64(rs1)
	b := c.GetRegisterFloatAsFloat64(rs2)
	d := c.GetRegisterFloatAsFloat64(rs3)
	r := a*b - d
	c.SetRegisterFloatAsFloat64(rd, r)
	if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaD) fnmsubd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fnmsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat64(rs1)
	b := c.GetRegisterFloatAsFloat64(rs2)
	d := c.GetRegisterFloatAsFloat64(rs3)
	r := a*b - d
	c.SetRegisterFloatAsFloat64(rd, -r)
	if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaD) fnmaddd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2, rs3 := R4Type(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s rs3: %s", c.GetPC(), "fnmadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
	c.ClrFloatFlag()
	a := c.GetRegisterFloatAsFloat64(rs1)
	b := c.GetRegisterFloatAsFloat64(rs2)
	d := c.GetRegisterFloatAsFloat64(rs3)
	r := a*b + d
	c.SetRegisterFloatAsFloat64(rd, -r)
	if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
		c.SetFloatFlag(FFlagsNX, 1)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

// func (_ *isaD) faddd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fsubd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fmuld(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fdivd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fsqrtd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fsgnjd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fsgnjnd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fsgnjxd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fmind(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fmaxd(c *CPU, i uint64) (uint64, error) {}

func (_ *isaD) fcvtsd(c *CPU, i uint64) (uint64, error) {
	rd, rs1, rs2 := RType(i)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
	d := c.GetRegisterFloatAsFloat64(rs1)
	if math.IsNaN(d) {
		c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
	} else {
		c.SetRegisterFloatAsFloat32(rd, float32(d))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

// func (_ *isaD) fcvtds(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) feqd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fltd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fled(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fclassd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtwd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtwud(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtdw(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtdwu(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtld(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtlud(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fmvxd(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtdl(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fcvtdlu(c *CPU, i uint64) (uint64, error) {}

// func (_ *isaD) fmvdx(c *CPU, i uint64) (uint64, error) {}

type isaC struct{}

func (_ *isaC) addi4spn(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 7, 10)<<6 | InstructionPart(i, 11, 12)<<4 | InstructionPart(i, 5, 5)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.addi4spn", c.LogI(rd), imm))
	if imm == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, c.GetRegister(Rsp)+imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fld(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "c.fld", c.LogF(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	v, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegisterFloat(rd, v)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) lw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 2, 4) + 8
		rs1 = InstructionPart(i, 7, 9) + 8
		imm = InstructionPart(i, 5, 5)<<6 | InstructionPart(i, 10, 12)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "c.lw", c.LogI(rd), c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ----(%#016x)", c.GetPC(), "c.ld", c.LogI(rd), c.LogI(rs1), imm))
	a := c.GetRegister(rs1) + imm
	b, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, b)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fsd(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 5, 6)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.fsd", c.LogI(rs1), c.LogF(rs2), imm))
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegisterFloat(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) sw(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 9) + 8
		rs2 = InstructionPart(i, 2, 4) + 8
		imm = InstructionPart(i, 5, 5)<<6 | InstructionPart(i, 10, 12)<<3 | InstructionPart(i, 6, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.sw", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.sd", c.LogI(rs1), c.LogI(rs2), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.addi", c.LogI(rd), imm))
	if rd == Rzero {
		return z.nop(c, i)
	}
	if imm == 0x00 {
		return 0, ErrHint
	}
	c.SetRegister(rd, c.GetRegister(rd)+imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) addiw(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.addiw", c.LogI(rd), imm))
	if rd == Rzero {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(rd, uint64(int32(c.GetRegister(rd))+int32(imm)))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) li(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.li", c.LogI(rd), imm))
	if rd == Rzero {
		return 0, ErrHint
	}
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) addi16sp(c *CPU, i uint64) (uint64, error) {
	var (
		imm = SignExtend(InstructionPart(i, 12, 12)<<9|InstructionPart(i, 3, 4)<<7|InstructionPart(i, 5, 5)<<6|InstructionPart(i, 2, 2)<<5|InstructionPart(i, 6, 6)<<4, 9)
	)
	Debugln(fmt.Sprintf("%#08x % 10s imm: ----(%#016x)", c.GetPC(), "c.addi16sp", imm))
	if imm == 0x00 {
		return 0, ErrReservedInstruction
	}
	c.SetRegister(Rsp, c.GetRegister(Rsp)+imm)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) lui(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = SignExtend(InstructionPart(i, 12, 12)<<17|InstructionPart(i, 2, 6)<<12, 17)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.lui", c.LogI(rd), imm))
	if imm == 0x00 {
		return 0, ErrHint
	}
	if rd == Rzero {
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
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.srli", c.LogI(rd), shamt))
	c.SetRegister(rd, c.GetRegister(rd)>>shamt)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) srai(c *CPU, i uint64) (uint64, error) {
	var (
		rd    = InstructionPart(i, 7, 9) + 8
		shamt = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.srai", c.LogI(rd), shamt))
	c.SetRegister(rd, uint64(int64(c.GetRegister(rd))>>shamt))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) andi(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 9) + 8
		imm = SignExtend(InstructionPart(i, 12, 12)<<5|InstructionPart(i, 2, 6), 5)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.andi", c.LogI(rd), imm))
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
	var imm = SignExtend(InstructionPart(i, 12, 12)<<11|
		InstructionPart(i, 8, 8)<<10|
		InstructionPart(i, 9, 10)<<8|
		InstructionPart(i, 6, 6)<<7|
		InstructionPart(i, 7, 7)<<6|
		InstructionPart(i, 2, 2)<<5|
		InstructionPart(i, 11, 11)<<4|
		InstructionPart(i, 3, 5)<<1, 11)
	Debugln(fmt.Sprintf("%#08x % 10s imm: ----(%#016x)", c.GetPC(), "c.j", imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s imm: ----(%#016x)", c.GetPC(), "c.beqz", c.LogI(rs1), imm))
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
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s imm: ----(%#016x)", c.GetPC(), "c.bnez", c.LogI(rs1), imm))
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
		rd    = InstructionPart(i, 7, 11)
		shamt = InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.slli", c.LogI(rd), shamt))
	if rd == 0 {
		return 0, ErrHint
	}
	c.SetRegister(rd, c.GetRegister(rd)<<shamt)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fldsp(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = InstructionPart(i, 2, 4)<<6 | InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 5, 6)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.fldsp", c.LogF(rd), imm))
	v, err := c.GetMemory().GetUint64(c.GetRegister(Rsp) + imm)
	if err != nil {
		return 0, err
	}
	c.SetRegisterFloat(rd, v)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) lwsp(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = InstructionPart(i, 2, 3)<<6 | InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 4, 6)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.lwsp", c.LogI(rd), imm))
	if rd == Rzero {
		return 0, ErrReservedInstruction
	}
	a := c.GetRegister(Rsp) + imm
	b, err := c.GetMemory().GetUint32(a)
	if err != nil {
		return 0, err
	}
	v := SignExtend(uint64(b), 31)
	c.SetRegister(rd, v)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) ldsp(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		imm = InstructionPart(i, 2, 4)<<6 | InstructionPart(i, 12, 12)<<5 | InstructionPart(i, 5, 6)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s imm: ----(%#016x)", c.GetPC(), "c.ldsp", c.LogI(rd), imm))
	if rd == Rzero {
		return 0, ErrReservedInstruction
	}
	a := c.GetRegister(Rsp) + imm
	b, err := c.GetMemory().GetUint64(a)
	if err != nil {
		return 0, err
	}
	c.SetRegister(rd, b)
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) jr(c *CPU, i uint64) (uint64, error) {
	var rs1 = InstructionPart(i, 7, 11)
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s", c.GetPC(), "c.jr", c.LogI(rs1)))
	if rs1 == 0 {
		return 0, ErrReservedInstruction
	}
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
	if rd == Rzero {
		return 0, ErrHint
	}
	c.SetRegister(rd, c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) ebreak(c *CPU, i uint64) (uint64, error) {
	Debugln(fmt.Sprintf("%#08x % 10s", c.GetPC(), "c.ebreak"))
	return 1, nil
}

func (_ *isaC) jalr(c *CPU, i uint64) (uint64, error) {
	var (
		rs1 = InstructionPart(i, 7, 11)
	)
	if rs1 == 0 {
		return 0, ErrReservedInstruction
	}
	Debugln(fmt.Sprintf("%#08x % 10s rs1: %s", c.GetPC(), "c.jalr", c.LogI(rs1)))
	c.SetRegister(Rra, c.GetPC()+2)
	c.SetPC(c.GetRegister(rs1) & 0xfffffffffffffffe)
	return 1, nil
}

func (_ *isaC) add(c *CPU, i uint64) (uint64, error) {
	var (
		rd  = InstructionPart(i, 7, 11)
		rs2 = InstructionPart(i, 2, 6)
	)
	Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs2: %s", c.GetPC(), "c.mv", c.LogI(rd), c.LogI(rs2)))
	if rd == Rzero {
		return 0, ErrHint
	}
	c.SetRegister(rd, c.GetRegister(rd)+c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) fsdsp(c *CPU, i uint64) (uint64, error) {
	var (
		rs2 = InstructionPart(i, 2, 6)
		imm = InstructionPart(i, 7, 9)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.sdsp", c.LogF(rs2), imm))
	a := c.GetRegister(Rsp) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegisterFloat(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) swsp(c *CPU, i uint64) (uint64, error) {
	var (
		rs2 = InstructionPart(i, 2, 6)
		imm = InstructionPart(i, 7, 8)<<6 | InstructionPart(i, 9, 12)<<2
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.swsp", c.LogI(rs2), imm))
	a := c.GetRegister(Rsp) + imm
	if err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

func (_ *isaC) sdsp(c *CPU, i uint64) (uint64, error) {
	var (
		rs2 = InstructionPart(i, 2, 6)
		imm = InstructionPart(i, 7, 9)<<6 | InstructionPart(i, 10, 12)<<3
	)
	Debugln(fmt.Sprintf("%#08x % 10s rs2: %s imm: ----(%#016x)", c.GetPC(), "c.sdsp", c.LogI(rs2), imm))
	a := c.GetRegister(Rsp) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegister(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 2)
	return 1, nil
}

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
