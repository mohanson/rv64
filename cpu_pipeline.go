package rv64

import (
	"fmt"
	"math"
	"math/big"
)

// The classic RISC pipeline comprises:
//   1. Instruction fetch
//   2. Instruction decode and register fetch
//   3. Execute
//   4. Memory access
//   5. Register write back
//
// | 1   | 2   | 3   | 4   | 5   | 6   | 7   | 8   | 9   |
// | --- | --- | --- | --- | --- | --- | --- | --- | --- |
// | IF  | ID  | EX  | MEM | WB  |     |     |     |     |
// |     | IF  | ID  | EX  | MEM | WB  |     |     |     |
// |     |     | IF  | ID  | EX  | MEM | WB  |     |     |
// |     |     |     | IF  | ID  | EX  | MEM | WB  |     |
// |     |     |     |     | IF  | ID  | EX  | MEM | WB  |

func (c *CPU) PipelineInstructionFetch() ([]byte, error) {
	a, err := c.GetMemory().GetByte(c.GetPC(), 2)
	if err != nil {
		return nil, err
	}
	b := InstructionLengthEncoding(a)
	r, err := c.GetMemory().GetByte(c.GetPC(), uint64(b))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *CPU) PipelineExecute(data []byte) (uint64, error) {
	switch len(data) {
	case 2:
		Panicln("Unreachable")
	case 4:
		var s uint64 = 0
		for i := len(data) - 1; i >= 0; i-- {
			s += uint64(data[i]) << (8 * i)
		}
		switch InstructionPart(s, 0, 6) {
		case 0b0110111: // ----------------------------------------------------------------------- LUI
			rd, imm := UType(s)
			Debugln(fmt.Sprintf("% 10s rd : %s imm: %#016x", "lui", c.LogI(rd), imm))
			c.SetRegister(rd, imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b0010111: // ----------------------------------------------------------------------- AUIPC
			rd, imm := UType(s)
			Debugln(fmt.Sprintf("% 10s rd : %s imm: %#016x", "auipc", c.LogI(rd), imm))
			c.SetRegister(rd, c.GetPC()+imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b1101111: // ----------------------------------------------------------------------- JAL
			rd, imm := JType(s)
			Debugln(fmt.Sprintf("% 10s rd : %s imm: %#016x", "jal", c.LogI(rd), imm))
			c.SetRegister(rd, c.GetPC()+4)
			r := c.GetPC() + imm
			if r%4 != 0x00 {
				return 0, ErrMisalignedInstructionFetch
			}
			c.SetPC(r)
			return 1, nil
		case 0b1100111: // ----------------------------------------------------------------------- JALR
			rd, rs1, imm := IType(s)
			Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "jalr", c.LogI(rd), c.LogI(rs1), imm))
			c.SetRegister(rd, c.GetPC()+4)
			r := (c.GetRegister(rs1) + imm) & 0xfffffffffffffffe
			if r%4 != 0x00 {
				return 0, ErrMisalignedInstructionFetch
			}
			c.SetPC(r)
			return 1, nil
		case 0b1100011:
			rs1, rs2, imm := BType(s)
			if imm%2 != 0x00 {
				return 0, ErrMisalignedInstructionFetch
			}
			var cond bool
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ BEQ
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "beq", c.LogI(rs1), c.LogI(rs2), imm))
				cond = c.GetRegister(rs1) == c.GetRegister(rs2)
			case 0b001: // ------------------------------------------------------------------------ BNE
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "bne", c.LogI(rs1), c.LogI(rs2), imm))
				cond = c.GetRegister(rs1) != c.GetRegister(rs2)
			case 0b100: // ------------------------------------------------------------------------ BLT
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "blt", c.LogI(rs1), c.LogI(rs2), imm))
				cond = int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2))
			case 0b101: // ------------------------------------------------------------------------ BGE
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "bge", c.LogI(rs1), c.LogI(rs2), imm))
				cond = int64(c.GetRegister(rs1)) >= int64(c.GetRegister(rs2))
			case 0b110: // ------------------------------------------------------------------------ BLTU
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "bltu", c.LogI(rs1), c.LogI(rs2), imm))
				cond = c.GetRegister(rs1) < c.GetRegister(rs2)
			case 0b111: // ------------------------------------------------------------------------ BGEU
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "bgeu", c.LogI(rs1), c.LogI(rs2), imm))
				cond = c.GetRegister(rs1) >= c.GetRegister(rs2)
			}
			if cond {
				c.SetPC(c.GetPC() + imm)
			} else {
				c.SetPC(c.GetPC() + 4)
			}
			return 1, nil
		case 0b0000011:
			rd, rs1, imm := IType(s)
			a := c.GetRegister(rs1) + imm
			var v uint64
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ LB
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lb", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint8(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 7)
			case 0b001: // ------------------------------------------------------------------------ LH
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lh", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint16(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 15)
			case 0b010: // ------------------------------------------------------------------------ LW
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lw", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 31)
			case 0b011: // ------------------------------------------------------------------------ LD
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "ld", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint64(a)
				if err != nil {
					return 0, err
				}
				v = b
			case 0b100: // ------------------------------------------------------------------------ LBU
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lbu", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint8(a)
				if err != nil {
					return 0, err
				}
				v = uint64(b)
			case 0b101: // ------------------------------------------------------------------------ LHU
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lhu", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint16(a)
				if err != nil {
					return 0, err
				}
				v = uint64(b)
			case 0b110: // ------------------------------------------------------------------------ LWU
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "lwu", c.LogI(rd), c.LogI(rs1), imm))
				b, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				v = uint64(b)
			}
			c.SetRegister(rd, v)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b0100011:
			rs1, rs2, imm := SType(s)
			a := c.GetRegister(rs1) + imm
			var err error
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ SB
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "sb", c.LogI(rs1), c.LogI(rs2), imm))
				err = c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2)))
			case 0b001: // ------------------------------------------------------------------------ SH
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "sh", c.LogI(rs1), c.LogI(rs2), imm))
				err = c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2)))
			case 0b010: // ------------------------------------------------------------------------ SW
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "sw", c.LogI(rs1), c.LogI(rs2), imm))
				err = c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
			case 0b011: // ------------------------------------------------------------------------ SD
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "sd", c.LogI(rs1), c.LogI(rs2), imm))
				err = c.GetMemory().SetUint64(a, c.GetRegister(rs2))
			}
			if err != nil {
				return 0, err
			}
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b0010011:
			rd, rs1, imm := IType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ ADDI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "addi", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, c.GetRegister(rs1)+imm)
			case 0b010: // ------------------------------------------------------------------------ SLTI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "slti", c.LogI(rd), c.LogI(rs1), imm))
				if int64(c.GetRegister(rs1)) < int64(imm) {
					c.SetRegister(rd, 1)
				} else {
					c.SetRegister(rd, 0)
				}
			case 0b011: // ------------------------------------------------------------------------ SLTIU
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "sltiu", c.LogI(rd), c.LogI(rs1), imm))
				if c.GetRegister(rs1) < imm {
					c.SetRegister(rd, 1)
				} else {
					c.SetRegister(rd, 0)
				}
			case 0b100: // ------------------------------------------------------------------------ XORI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "xori", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, c.GetRegister(rs1)^imm)
			case 0b110: // ------------------------------------------------------------------------ ORI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "ori", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, c.GetRegister(rs1)|imm)
			case 0b111: // ------------------------------------------------------------------------ ANDI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "andi", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, c.GetRegister(rs1)&imm)
			case 0b001: // ------------------------------------------------------------------------ SLLI
				shamt := InstructionPart(imm, 0, 5)
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "slli", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, c.GetRegister(rs1)<<shamt)
			case 0b101:
				shamt := InstructionPart(imm, 0, 5)
				switch InstructionPart(s, 26, 31) {
				case 0b000000: // ----------------------------------------------------------------- SRLI
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "srli", c.LogI(rd), c.LogI(rs1), imm))
					c.SetRegister(rd, c.GetRegister(rs1)>>shamt)
				case 0b010000: // ----------------------------------------------------------------- SRAI
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "srai", c.LogI(rd), c.LogI(rs1), imm))
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>shamt))
				}
			}
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b0110011:
			rd, rs1, rs2 := RType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- ADD
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "add", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)+c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MUL
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "mul", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SUB
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "sub", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)-c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b001:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SLL
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "sll", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)<<(c.GetRegister(rs2)&0x3f))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULH
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "mulh", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			case 0b010:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SLT
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "slt", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULHSU
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "mulhsu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			case 0b011:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SLTU
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "sltu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs1) < c.GetRegister(rs2) {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULHU
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "mulhu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			case 0b100:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- XOR
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "xor", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)^c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIV
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "div", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, math.MaxUint64)
					} else {
						c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))/int64(c.GetRegister(rs2))))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b101:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SRL
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "srl", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)>>(c.GetRegister(rs2)&0x3f))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIVU
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "divu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, math.MaxUint64)
					} else {
						c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRA
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "sra", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>(c.GetRegister(rs2)&0x3f)))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b110:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- OR
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "or", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)|c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- REM
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "rem", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, c.GetRegister(rs1))
					} else {
						c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))%int64(c.GetRegister(rs2))))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b111:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- AND
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "and", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, c.GetRegister(rs1)&c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- REMU
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "remu", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, c.GetRegister(rs1))
					} else {
						c.SetRegister(rd, c.GetRegister(rs1)%c.GetRegister(rs2))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			}
		case 0b0001111:
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ---------------------------------------------------------------------- FENCE
				Debugln(fmt.Sprintf("% 10s", "fence"))
			case 0b001: // ---------------------------------------------------------------------- FENCE.I
				Debugln(fmt.Sprintf("% 10s", "fence.i"))
			}
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b1110011:
			rd, rs1, csr := IType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000:
				switch InstructionPart(s, 20, 31) {
				case 0b000000000000: // ----------------------------------------------------------- ECALL
					Debugln(fmt.Sprintf("% 10s", "ecall"))
					return c.GetSystem().HandleCall(c)
				case 0b000000000001: // ----------------------------------------------------------- EBREAK
					Debugln(fmt.Sprintf("% 10s", "ebreak"))
					return 1, nil
				}
			case 0b001: // ------------------------------------------------------------------------ CSRRW
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrw", c.LogI(rd), c.LogI(rs1), csr))
				if rd != Rzero {
					c.SetRegister(rd, c.GetCSR().Get(csr))
				}
				c.GetCSR().Set(csr, c.GetRegister(rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b010: // ------------------------------------------------------------------------ CSRRS
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrs", c.LogI(rd), c.LogI(rs1), csr))
				c.SetRegister(rd, c.GetCSR().Get(csr))
				if rs1 != Rzero {
					c.GetCSR().Set(csr, c.GetCSR().Get(csr)|c.GetRegister(rs1))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ CSRRC
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrc", c.LogI(rd), c.LogI(rs1), csr))
				c.SetRegister(rd, c.GetCSR().Get(csr))
				if rs1 != Rzero {
					c.GetCSR().Set(csr, c.GetCSR().Get(csr)&(math.MaxUint64-c.GetRegister(rs1)))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101: // ------------------------------------------------------------------------ CSRRWI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrwi", c.LogI(rd), c.LogI(rs1), csr))
				if rd != Rzero {
					c.SetRegister(rd, c.GetCSR().Get(csr))
				}
				c.GetCSR().Set(csr, rs1)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b110: // ------------------------------------------------------------------------ CSRRSI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrsi", c.LogI(rd), c.LogI(rs1), csr))
				c.SetRegister(rd, c.GetCSR().Get(csr))
				if csr != 0x00 {
					c.GetCSR().Set(csr, c.GetCSR().Get(csr)|rs1)
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b111: // ------------------------------------------------------------------------ CSRRCI
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s csr: %#016x", "csrrci", c.LogI(rd), c.LogI(rs1), csr))
				c.SetRegister(rd, c.GetCSR().Get(csr))
				if csr != 0x00 {
					c.GetCSR().Set(csr, c.GetCSR().Get(csr)&(math.MaxUint64-rs1))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b0011011:
			rd, rs1, imm := IType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ ADDIW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s imm: %#016x", "addiw", c.LogI(rd), c.LogI(rs1), imm))
				c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(imm)))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b001: // ------------------------------------------------------------------------ SLLIW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s imm: %#016x", "slliw", c.LogI(rd), c.LogI(rs1), imm))
				if InstructionPart(imm, 5, 5) != 0x00 {
					return 0, ErrAbnormalInstruction
				}
				c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<imm), 31))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SRLIW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s imm: %#016x", "srliw", c.LogI(rd), c.LogI(rs1), imm))
					if InstructionPart(imm, 5, 5) != 0x00 {
						return 0, ErrAbnormalInstruction
					}
					shamt := InstructionPart(imm, 0, 4)
					c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>shamt), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRAIW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s imm: %#016x", "sraiw", c.LogI(rd), c.LogI(rs1), imm))
					if InstructionPart(imm, 5, 5) != 0x00 {
						return 0, ErrAbnormalInstruction
					}
					shamt := InstructionPart(imm, 0, 4)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>shamt))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			}
		case 0b0111011:
			rd, rs1, rs2 := RType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- ADDW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "addw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "mulw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SUBW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "subw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))-int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b001: // ------------------------------------------------------------------------ SLLW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "sllw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
				s := c.GetRegister(rs2) & 0x1f
				c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<s), 31))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b100: // ------------------------------------------------------------------------ DIVW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "divw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
				if c.GetRegister(rs2) == 0 {
					c.SetRegister(rd, math.MaxUint64)
				} else {
					c.SetRegister(rd, SignExtend(uint64(int32(c.GetRegister(rs1))/int32(c.GetRegister(rs2))), 31))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SRLW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "srlw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					s := c.GetRegister(rs2) & 0x1f
					c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>s), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIVUW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "divuw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, math.MaxUint64)
					} else {
						c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRAW
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "sraw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b110: // ------------------------------------------------------------------------ REMW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "remw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
				if c.GetRegister(rs2) == 0 {
					c.SetRegister(rd, c.GetRegister(rs1))
				} else {
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))%int32(c.GetRegister(rs2))))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b111: // ------------------------------------------------------------------------ REMUW
				Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "remuw", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
				if c.GetRegister(rs2) == 0 {
					c.SetRegister(rd, c.GetRegister(rs1))
				} else {
					c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))%uint32(c.GetRegister(rs2))), 31))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b0101111:
			rd, rs1, rs2 := RType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b010:
				a := SignExtend(c.GetRegister(rs1), 31)
				switch InstructionPart(s, 27, 31) {
				case 0b00010: // ------------------------------------------------------------------ LR.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "lr.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ SC.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "sc.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if a == c.GetLoadReservation() {
						c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
						c.SetRegister(rd, 0)
					} else {
						c.SetRegister(rd, 1)
					}
					c.SetLoadReservation(0)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ AMOSWAP.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoswap.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000: // ------------------------------------------------------------------ AMOADD.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoadd.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v+uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100: // ------------------------------------------------------------------ AMOXOR.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoxor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v^uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100: // ------------------------------------------------------------------ AMOAND.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoand.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v&uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000: // ------------------------------------------------------------------ AMOOR.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v|uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b10000: // ------------------------------------------------------------------ AMOMIN.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomin.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b10100: // ------------------------------------------------------------------ AMOMAX.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomax.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11000: // ------------------------------------------------------------------ AMOMINU.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amominu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11100: // ------------------------------------------------------------------ AMOMAXU.W
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomaxu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			case 0b011:
				a := c.GetRegister(rs1)
				switch InstructionPart(s, 27, 31) {
				case 0b00010: // ------------------------------------------------------------------ LR.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "lr.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, v)
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ SC.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "sc.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if a == c.GetLoadReservation() {
						c.GetMemory().SetUint64(a, c.GetRegister(rs2))
						c.SetRegister(rd, 0)
					} else {
						c.SetRegister(rd, 1)
					}
					c.SetLoadReservation(0)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ AMOSWAP.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoswap.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000: // ------------------------------------------------------------------ AMOADD.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoadd.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v+c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100: // ------------------------------------------------------------------ AMOXOR.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoxor.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v^c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100: // ------------------------------------------------------------------ AMOAND.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoand.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v&c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000: // ------------------------------------------------------------------ AMOOR.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amoor.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v|c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b10000: // ------------------------------------------------------------------ AMOMIN.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomin.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b10100: // ------------------------------------------------------------------ AMOMAX.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomax.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11000: // ------------------------------------------------------------------ AMOMINU.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amominu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11100: // ------------------------------------------------------------------ AMOMAXU.D
					Debugln(fmt.Sprintf("% 10s rd: %s rs1: %s rs2: %s", "amomaxu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			}
		case 0b0000111:
			rd, rs1, imm := IType(s)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(s, 12, 14) {
			case 0b010: // ------------------------------------------------------------------------ FLW
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "flw", c.LogF(rd), c.LogI(rs1), imm))
				v, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				c.SetRegisterFloatAsFloat32(rd, math.Float32frombits(v))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ FLD
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s imm: %#016x", "fld", c.LogF(rd), c.LogI(rs1), imm))
				v, err := c.GetMemory().GetUint64(a)
				if err != nil {
					return 0, err
				}
				c.SetRegisterFloat(rd, v)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b0100111:
			rs1, rs2, imm := SType(s)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(s, 12, 14) {
			case 0b010: // ------------------------------------------------------------------------ FSW
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "fsw", c.LogI(rs1), c.LogF(rs2), imm))
				err := c.GetMemory().SetUint32(a, uint32(c.GetRegisterFloat(rs2)))
				if err != nil {
					return 0, err
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ FSD
				Debugln(fmt.Sprintf("% 10s rs1: %s rs2: %s imm: %#016x", "fsd", c.LogI(rs1), c.LogF(rs2), imm))
				err := c.GetMemory().SetUint64(a, c.GetRegisterFloat(rs2))
				if err != nil {
					return 0, err
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b1000011:
			rd, rs1, rs2, rs3 := R4Type(s)
			switch InstructionPart(s, 25, 26) {
			case 0b00: // ------------------------------------------------------------------------- FMADD.S
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fmadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
			case 0b01: // ------------------------------------------------------------------------- FMADD.D
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fmadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
		case 0b1000111:
			rd, rs1, rs2, rs3 := R4Type(s)
			switch InstructionPart(s, 25, 26) {
			case 0b00: // ------------------------------------------------------------------------- FMSUB.S
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fmsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
			case 0b01: // ------------------------------------------------------------------------- FMSUB.D
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fmsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
		case 0b1001011:
			rd, rs1, rs2, rs3 := R4Type(s)
			switch InstructionPart(s, 25, 26) {
			case 0b00: // ------------------------------------------------------------------------- FNMSUB.S
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fnmsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
			case 0b01: // ------------------------------------------------------------------------- FNMSUB.D
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fnmsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
		case 0b1001111:
			rd, rs1, rs2, rs3 := R4Type(s)
			switch InstructionPart(s, 25, 26) {
			case 0b00: // ------------------------------------------------------------------------- FNMADD.S
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fnmadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
			case 0b01: // ------------------------------------------------------------------------- FNMADD.D
				Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s rs3: %s", "fnmadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2), c.LogF(rs3)))
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
		case 0b1010011:
			rd, rs1, rs2 := RType(s)
			switch InstructionPart(s, 25, 26) {
			case 0b00:
				a := c.GetRegisterFloatAsFloat32(rs1)
				b := c.GetRegisterFloatAsFloat32(rs2)
				switch InstructionPart(s, 27, 31) {
				case 0b00000: // ------------------------------------------------------------------ FADD.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					d := a + b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d-a != b || d-b != a {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ FSUB.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b00010: // ------------------------------------------------------------------ FMUL.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmul.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					d := a * b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d/a != b || d/b != a || float64(a)*float64(b) != float64(d) {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ FDIV.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fdiv.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b01011: // ------------------------------------------------------------------ FSQRT.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsqrt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b00100:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FSGNJ.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnj.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FSGNJN.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnjn.s.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010: // ---------------------------------------------------------------- FSGNJX.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnjx.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(a)) != math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b00101:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FMIN.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmin.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					case 0b001: // ---------------------------------------------------------------- FMAX.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmax.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					}
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
					switch InstructionPart(s, 12, 14) {
					case 0b000:
						if (math.Signbit(float64(a)) && !math.Signbit(float64(b))) || a < b {
							c.SetRegisterFloatAsFloat32(rd, a)
						} else {
							c.SetRegisterFloatAsFloat32(rd, b)
						}
					case 0b001:
						if (!math.Signbit(float64(a)) && math.Signbit(float64(b))) || a > b {
							c.SetRegisterFloatAsFloat32(rd, a)
						} else {
							c.SetRegisterFloatAsFloat32(rd, b)
						}
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11000:
					switch InstructionPart(s, 20, 24) {
					case 0b00000: // -------------------------------------------------------------- FCVT.W.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.w.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00001: // -------------------------------------------------------------- FCVT.WU.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.wu.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00010: // -------------------------------------------------------------- FCVT.L.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.l.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00011: // -------------------------------------------------------------- FCVT.LU.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.lu.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b01000: // ------------------------------------------------------------------ FCVT.S.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.s.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					d := c.GetRegisterFloatAsFloat64(rs1)
					if math.IsNaN(d) {
						c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
					} else {
						c.SetRegisterFloatAsFloat32(rd, float32(d))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11100:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FMV.X.W
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmv.x.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegisterFloat(rs1))), 31))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FCLASS.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fclass.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						a := c.GetRegisterFloatAsFloat32(rs1)
						c.SetRegister(rd, FClassS(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b10100:
					var cond bool
					switch InstructionPart(s, 12, 14) {
					case 0b010: // ---------------------------------------------------------------- FEQ.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "feq.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if IsSNaN32(a) || IsSNaN32(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001: // ---------------------------------------------------------------- FLT.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "flt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000: // ---------------------------------------------------------------- FLE.S
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fle.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a <= b
						}
					}
					if cond {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11010:
					switch InstructionPart(s, 20, 24) {
					case 0b00000: // -------------------------------------------------------------- FCVT.S.W
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.s.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001: // -------------------------------------------------------------- FCVT.S.WU
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.s.wu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010: // -------------------------------------------------------------- FCVT.S.L
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.s.l", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011: // -------------------------------------------------------------- FCVT.S.LU
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.s.lu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110: // ------------------------------------------------------------------ FMV.W.X
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmv.w.x", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(uint32(c.GetRegister(rs1))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b01:
				a := c.GetRegisterFloatAsFloat64(rs1)
				b := c.GetRegisterFloatAsFloat64(rs2)
				switch InstructionPart(s, 27, 31) {
				case 0b00000: // ------------------------------------------------------------------ FADD.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a+b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ FSUB.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					if (math.Signbit(a) == math.Signbit(b)) && math.IsInf(a, 0) && math.IsInf(b, 0) {
						c.SetRegisterFloat(rd, NaN64)
						c.SetFloatFlag(FFlagsNV, 1)
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					c.SetRegisterFloatAsFloat64(rd, a-b)
					if big.NewFloat(0).Sub(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00010: // ------------------------------------------------------------------ FMUL.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmul.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a*b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ FDIV.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fdiv.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					if b == 0 {
						c.SetRegisterFloat(rd, NaN64)
						c.SetFloatFlag(FFlagsDZ, 1)
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					c.SetRegisterFloatAsFloat64(rd, a/b)
					if big.NewFloat(0).Quo(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01011: // ------------------------------------------------------------------ FSQRT.D
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsqrt.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					if a < 0 {
						c.SetRegisterFloat(rd, NaN64)
						c.SetFloatFlag(FFlagsNV, 1)
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					c.SetRegisterFloatAsFloat64(rd, math.Sqrt(a))
					d := big.NewFloat(0).Sqrt(big.NewFloat(a))
					if big.NewFloat(0).Mul(d, d).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FSGNJ.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnj.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FSGNJN.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnjn.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010: // ---------------------------------------------------------------- FSGNJX.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fsgnjx.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(a) != math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b00101:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FMIN.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmin.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					case 0b001: // ---------------------------------------------------------------- FMAX.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmax.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					}
					c.ClrFloatFlag()
					if math.IsNaN(a) && math.IsNaN(b) {
						c.SetRegisterFloat(rd, NaN64)
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					if math.IsNaN(a) {
						c.SetRegisterFloatAsFloat64(rd, b)
						if IsSNaN64(a) {
							c.SetFloatFlag(FFlagsNV, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					if math.IsNaN(b) {
						c.SetRegisterFloatAsFloat64(rd, a)
						if IsSNaN64(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
					switch InstructionPart(s, 12, 14) {
					case 0b000:
						if (math.Signbit(a) && !math.Signbit(b)) || a < b {
							c.SetRegisterFloatAsFloat64(rd, a)
						} else {
							c.SetRegisterFloatAsFloat64(rd, b)
						}
					case 0b001:
						if (!math.Signbit(a) && math.Signbit(b)) || a > b {
							c.SetRegisterFloatAsFloat64(rd, a)
						} else {
							c.SetRegisterFloatAsFloat64(rd, b)
						}
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11000:
					switch InstructionPart(s, 20, 24) {
					case 0b00000: // -------------------------------------------------------------- FCVT.W.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.w.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						d := c.GetRegisterFloatAsFloat64(rs1)
						if math.IsNaN(d) {
							c.SetRegister(rd, 0x7fffffff)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d > float64(math.MaxInt32) {
							c.SetRegister(rd, SignExtend(0x7fffffff, 31))
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d < float64(math.MinInt32) {
							c.SetRegister(rd, SignExtend(0x80000000, 31))
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						c.SetRegister(rd, SignExtend(uint64(int32(d)), 31))
						if math.Ceil(d) != d {
							c.SetFloatFlag(FFlagsNX, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001: // -------------------------------------------------------------- FCVT.WU.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.wu.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						d := c.GetRegisterFloatAsFloat64(rs1)
						if math.IsNaN(d) {
							c.SetRegister(rd, 0xffffffffffffffff)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d > float64(math.MaxUint32) {
							c.SetRegister(rd, SignExtend(0xffffffff, 31))
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d <= float64(-1) {
							c.SetRegister(rd, SignExtend(0x00000000, 31))
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						c.SetRegister(rd, SignExtend(uint64(uint32(d)), 31))
						if math.Ceil(d) != d {
							c.SetFloatFlag(FFlagsNX, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010: // -------------------------------------------------------------- FCVT.L.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.l.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						d := c.GetRegisterFloatAsFloat64(rs1)
						if math.IsNaN(d) {
							c.SetRegister(rd, 0x7fffffffffffffff)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d > float64(math.MaxInt64) {
							c.SetRegister(rd, 0x7fffffffffffffff)
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d < float64(math.MinInt64) {
							c.SetRegister(rd, 0x8000000000000000)
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						c.SetRegister(rd, uint64(int64(d)))
						if math.Ceil(d) != d {
							c.SetFloatFlag(FFlagsNX, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011: // -------------------------------------------------------------- FCVT.LU.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.lu.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						d := c.GetRegisterFloatAsFloat64(rs1)
						if math.IsNaN(d) {
							c.SetRegister(rd, 0xffffffffffffffff)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d > float64(math.MaxUint64) {
							c.SetRegister(rd, 0xffffffffffffffff)
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						if d <= float64(-1) {
							c.SetRegister(rd, 0x0000000000000000)
							c.SetFloatFlag(FFlagsNV, 1)
							c.SetPC(c.GetPC() + 4)
							return 1, nil
						}
						c.SetRegister(rd, uint64(d))
						if math.Ceil(d) != d {
							c.SetFloatFlag(FFlagsNX, 1)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b01000: // ------------------------------------------------------------------ FCVT.D.S
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.d.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					d := c.GetRegisterFloatAsFloat32(rs1)
					if math.IsNaN(float64(d)) {
						c.SetRegisterFloat(rd, NaN64)
					} else {
						c.SetRegisterFloatAsFloat64(rd, float64(d))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b10100:
					var cond bool
					switch InstructionPart(s, 12, 14) {
					case 0b010: // ---------------------------------------------------------------- FEQ.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "feq.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if IsSNaN64(a) || IsSNaN64(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001: // ---------------------------------------------------------------- FLT.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "flt.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(a) || math.IsNaN(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000: // ---------------------------------------------------------------- FLE.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fle.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(a) || math.IsNaN(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a <= b
						}
					}
					if cond {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11100:
					switch InstructionPart(s, 12, 14) {
					case 0b000: // ---------------------------------------------------------------- FMV.X.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmv.x.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegister(rd, c.GetRegisterFloat(rs1))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FCLASS.D
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fclass.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						a := c.GetRegisterFloatAsFloat64(rs1)
						c.SetRegister(rd, FClassD(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11010:
					switch InstructionPart(s, 20, 24) {
					case 0b00000: // -------------------------------------------------------------- FCVT.D.W
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.d.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001: // -------------------------------------------------------------- FCVT.D.WU
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.d.wu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010: // -------------------------------------------------------------- FCVT.D.L
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.d.l", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011: // -------------------------------------------------------------- FCVT.D.LU
						Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fcvt.d.lu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110: // ------------------------------------------------------------------ FMV.D.X
					Debugln(fmt.Sprintf("% 10s rd : %s rs1: %s rs2: %s", "fmv.d.x", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.SetRegisterFloat(rd, c.GetRegister(rs1))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			}
		}
	}
	return 0, nil
}

func IsQNaN32(f float32) bool {
	return math.IsNaN(float64(f)) && math.Float32bits(f)&0x00400000 != 0x00
}

func IsSNaN32(f float32) bool {
	return math.IsNaN(float64(f)) && math.Float32bits(f)&0x00400000 == 0x00
}

func IsSubmoduleFloat32(f float32) bool {
	b := math.Float32bits(f)
	return b&0x7f800000 == 0 && b&0x000fffff != 0
}

func IsQNaN64(f float64) bool {
	return math.IsNaN(f) && math.Float64bits(f)&0x0008000000000000 != 0x00
}

func IsSNaN64(f float64) bool {
	return math.IsNaN(f) && math.Float64bits(f)&0x0008000000000000 == 0x00
}

func IsSubmoduleFloat64(f float64) bool {
	b := math.Float64bits(f)
	return b&0x7ff0000000000000 == 0 && b&0x000fffffffffffff != 0
}

func FClassS(f float32) uint64 {
	s := math.Float32bits(f)&(1<<31) != 0
	if IsSNaN32(f) {
		return 0b01_00000000
	}
	if IsQNaN32(f) {
		return 0b10_00000000
	}
	if s {
		if f < -math.MaxFloat32 {
			return 0b00_00000001
		} else if f == 0 {
			return 0b00_00001000
		} else if IsSubmoduleFloat32(f) {
			return 0b00_00000100
		} else {
			return 0b00_00000010
		}
	}
	if f > math.MaxFloat32 {
		return 0b00_10000000
	} else if f == 0 {
		return 0b00_00010000
	} else if IsSubmoduleFloat32(f) {
		return 0b00_00100000
	} else {
		return 0b00_01000000
	}
}

func FClassD(f float64) uint64 {
	s := math.Signbit(f)
	if IsSNaN64(f) {
		return 0b01_00000000
	}
	if IsQNaN64(f) {
		return 0b10_00000000
	}
	if s {
		if f < -math.MaxFloat64 {
			return 0b00_00000001
		} else if f == 0 {
			return 0b00_00001000
		} else if IsSubmoduleFloat64(f) {
			return 0b00_00000100
		} else {
			return 0b00_00000010
		}
	}
	if f > math.MaxFloat64 {
		return 0b00_10000000
	} else if f == 0 {
		return 0b00_00010000
	} else if IsSubmoduleFloat64(f) {
		return 0b00_00100000
	} else {
		return 0b00_01000000
	}
}
