package rv64

import (
	"math"
	"math/big"
)

type isaI struct{}

func (_ *isaI) lui(c *CPU, rd uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) aupic(c *CPU, rd uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetPC()+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) jal(c *CPU, rd uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetPC()+4)
	r := c.GetPC() + imm
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaI) jalr(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetPC()+4)
	r := (c.GetRegister(rs1) + imm) & 0xfffffffffffffffe
	if r%2 != 0x00 {
		return 0, ErrMisalignedInstructionFetch
	}
	c.SetPC(r)
	return 1, nil
}

func (_ *isaI) beq(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) bne(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) blt(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) bge(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) bltu(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) bgeu(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lb(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lh(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lw(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) ld(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lbu(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lhu(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) lwu(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
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

func (_ *isaI) sb(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sh(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sw(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2))); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sd(c *CPU, rs1 uint64, rs2 uint64, imm uint64) (uint64, error) {
	a := c.GetRegister(rs1) + imm
	if err := c.GetMemory().SetUint64(a, c.GetRegister(rs2)); err != nil {
		return 0, err
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) addi(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)+imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slti(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	if int64(c.GetRegister(rs1)) < int64(imm) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sltiu(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	if c.GetRegister(rs1) < imm {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) xori(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)^imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) ori(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)|imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) andi(c *CPU, rd uint64, rs1 uint64, imm uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)&imm)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slli(c *CPU, rd uint64, rs1 uint64, shamt uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)<<shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srli(c *CPU, rd uint64, rs1 uint64, shamt uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)>>shamt)
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srai(c *CPU, rd uint64, rs1 uint64, shamt uint64) (uint64, error) {
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>shamt))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) add(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)+c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sub(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)-c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (i *isaI) sll(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)<<(c.GetRegister(rs2)&0x3f))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) slt(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) sltu(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs1) < c.GetRegister(rs2) {
		c.SetRegister(rd, 1)
	} else {
		c.SetRegister(rd, 0)
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) xor(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)^c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) srl(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)>>(c.GetRegister(rs2)&0x3f))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}
func (_ *isaI) sra(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>(c.GetRegister(rs2)&0x3f)))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) or(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)|c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) and(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, c.GetRegister(rs1)&c.GetRegister(rs2))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaI) fenci()  {}
func (_ *isaI) ecall()  {}
func (_ *isaI) ebreak() {}

func (_ *isaI) addiw() {}
func (_ *isaI) slliw() {}
func (_ *isaI) srliw() {}
func (_ *isaI) sraiw() {}
func (_ *isaI) addw()  {}
func (_ *isaI) subw()  {}
func (_ *isaI) sllw()  {}
func (_ *isaI) srlw()  {}
func (_ *isaI) sraw()  {}

type isaZifencei struct{}

func (_ *isaZifencei) fencei() {}

type isaZicsr struct{}

func (_ *isaZicsr) csrrw()  {}
func (_ *isaZicsr) csrrs()  {}
func (_ *isaZicsr) csrrc()  {}
func (_ *isaZicsr) csrrwi() {}
func (_ *isaZicsr) csrrsi() {}
func (_ *isaZicsr) csrrci() {}

type isaM struct{}

func (_ *isaM) mul(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulh(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
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

func (_ *isaM) mulhsu(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
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

func (_ *isaM) mulhu(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
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

func (_ *isaM) div(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))/int64(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divu(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) rem(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))%int64(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) remu(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, c.GetRegister(rs1)%c.GetRegister(rs2))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) mulw(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divw(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) divuw(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, math.MaxUint64)
	} else {
		c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}

func (_ *isaM) remw(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
	if c.GetRegister(rs2) == 0 {
		c.SetRegister(rd, c.GetRegister(rs1))
	} else {
		c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))%int32(c.GetRegister(rs2))))
	}
	c.SetPC(c.GetPC() + 4)
	return 1, nil
}
func (_ *isaM) remuw(c *CPU, rd uint64, rs1 uint64, rs2 uint64) (uint64, error) {
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

func (_ *isaC) addi4spn() {}
func (_ *isaC) fld()      {}
func (_ *isaC) lw()       {}
func (_ *isaC) ld()       {}
func (_ *isaC) fsd()      {}
func (_ *isaC) sw()       {}
func (_ *isaC) sd()       {}
func (_ *isaC) nop()      {}
func (_ *isaC) addi()     {}
func (_ *isaC) addiw()    {}
func (_ *isaC) li()       {}
func (_ *isaC) addi16sp() {}
func (_ *isaC) lui()      {}
func (_ *isaC) srli64()   {}
func (_ *isaC) srai64()   {}
func (_ *isaC) andi()     {}
func (_ *isaC) sub()      {}
func (_ *isaC) xor()      {}
func (_ *isaC) or()       {}
func (_ *isaC) and()      {}
func (_ *isaC) subw()     {}
func (_ *isaC) addw()     {}
func (_ *isaC) j()        {}
func (_ *isaC) beqz()     {}
func (_ *isaC) bnez()     {}
func (_ *isaC) slli64()   {}
func (_ *isaC) fldsp()    {}
func (_ *isaC) lwsp()     {}
func (_ *isaC) ldsp()     {}
func (_ *isaC) jr()       {}
func (_ *isaC) mv()       {}
func (_ *isaC) ebreak()   {}
func (_ *isaC) jalr()     {}
func (_ *isaC) add()      {}
func (_ *isaC) fsdsp()    {}
func (_ *isaC) sqsp()     {}
func (_ *isaC) swsp()     {}
func (_ *isaC) sdsp()     {}

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
