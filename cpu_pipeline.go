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
			imm = SignExtend(imm, 31)
			DebuglnUType("LUI", rd, imm)
			c.SetRegister(rd, imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b0010111: // ----------------------------------------------------------------------- AUIPC
			rd, imm := UType(s)
			imm = SignExtend(imm, 31)
			DebuglnUType("AUIPC", rd, imm)
			c.SetRegister(rd, c.GetPC()+imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b1101111: // ----------------------------------------------------------------------- JAL
			rd, imm := JType(s)
			imm = SignExtend(imm, 19)
			DebuglnJType("JAL", rd, imm)
			c.SetRegister(rd, c.GetPC()+4)
			c.SetPC(c.GetPC() + imm)
			return 1, nil
		case 0b1100111: // ----------------------------------------------------------------------- JALR
			rd, rs1, imm := IType(s)
			imm = SignExtend(imm, 11)
			DebuglnIType("JALR", rd, rs1, imm)
			c.SetRegister(rd, c.GetPC()+4)
			c.SetPC(((c.GetRegister(rs1) + imm) >> 1) << 1)
			return 1, nil
		case 0b1100011:
			rs1, rs2, imm := BType(s)
			imm = SignExtend(imm, 12)
			var cond bool
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ BEQ
				DebuglnBType("BEQ", rs1, rs2, imm)
				cond = c.GetRegister(rs1) == c.GetRegister(rs2)
			case 0b001: // ------------------------------------------------------------------------ BNE
				DebuglnBType("BNE", rs1, rs2, imm)
				cond = c.GetRegister(rs1) != c.GetRegister(rs2)
			case 0b100: // ------------------------------------------------------------------------ BLT
				DebuglnBType("BLT", rs1, rs2, imm)
				cond = int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2))
			case 0b101: // ------------------------------------------------------------------------ BGE
				DebuglnBType("BGE", rs1, rs2, imm)
				cond = int64(c.GetRegister(rs1)) >= int64(c.GetRegister(rs2))
			case 0b110: // ------------------------------------------------------------------------ BLTU
				DebuglnBType("BLTU", rs1, rs2, imm)
				cond = c.GetRegister(rs1) < c.GetRegister(rs2)
			case 0b111: // ------------------------------------------------------------------------ BGEU
				DebuglnBType("BGEU", rs1, rs2, imm)
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
			imm = SignExtend(imm, 11)
			a := c.GetRegister(rs1) + imm
			var v uint64
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ LB
				DebuglnIType("LB", rd, rs1, imm)
				b, err := c.GetMemory().GetUint8(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 7)
			case 0b001: // ------------------------------------------------------------------------ LH
				DebuglnIType("LH", rd, rs1, imm)
				b, err := c.GetMemory().GetUint16(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 15)
			case 0b010: // ------------------------------------------------------------------------ LW
				DebuglnIType("LW", rd, rs1, imm)
				b, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				v = SignExtend(uint64(b), 31)
			case 0b011: // ------------------------------------------------------------------------ LD
				DebuglnIType("LD", rd, rs1, imm)
				b, err := c.GetMemory().GetUint64(a)
				if err != nil {
					return 0, err
				}
				v = b
			case 0b100: // ------------------------------------------------------------------------ LBU
				DebuglnIType("LBU", rd, rs1, imm)
				b, err := c.GetMemory().GetUint8(a)
				if err != nil {
					return 0, err
				}
				v = uint64(b)
			case 0b101: // ------------------------------------------------------------------------ LHU
				DebuglnIType("LHU", rd, rs1, imm)
				b, err := c.GetMemory().GetUint16(a)
				if err != nil {
					return 0, err
				}
				v = uint64(b)
			case 0b110: // ------------------------------------------------------------------------ LWU
				DebuglnIType("LWU", rd, rs1, imm)
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
			imm = SignExtend(imm, 11)
			a := c.GetRegister(rs1) + imm
			var err error
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ SB
				DebuglnIType("SB", rs1, rs2, imm)
				err = c.GetMemory().SetUint8(a, uint8(c.GetRegister(rs2)))
			case 0b001: // ------------------------------------------------------------------------ SH
				DebuglnIType("SH", rs1, rs2, imm)
				err = c.GetMemory().SetUint16(a, uint16(c.GetRegister(rs2)))
			case 0b010: // ------------------------------------------------------------------------ SW
				DebuglnIType("SW", rs1, rs2, imm)
				err = c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
			case 0b011: // ------------------------------------------------------------------------ SD
				DebuglnIType("SD", rs1, rs2, imm)
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
				imm = SignExtend(imm, 11)
				DebuglnIType("ADDI", rd, rs1, imm)
				c.SetRegister(rd, c.GetRegister(rs1)+imm)
			case 0b010: // ------------------------------------------------------------------------ SLTI
				imm = SignExtend(imm, 11)
				DebuglnIType("SLTI", rd, rs1, imm)
				if int64(c.GetRegister(rs1)) < int64(imm) {
					c.SetRegister(rd, 1)
				} else {
					c.SetRegister(rd, 0)
				}
			case 0b011: // ------------------------------------------------------------------------ SLTIU
				imm = uint64(SignExtend(imm, 11))
				DebuglnIType("SLTIU", rd, rs1, imm)
				if c.GetRegister(rs1) < imm {
					c.SetRegister(rd, 1)
				} else {
					c.SetRegister(rd, 0)
				}
			case 0b100: // ------------------------------------------------------------------------ XORI
				imm = SignExtend(imm, 11)
				DebuglnIType("XORI", rd, rs1, imm)
				c.SetRegister(rd, c.GetRegister(rs1)^imm)
			case 0b110: // ------------------------------------------------------------------------ ORI
				imm = SignExtend(imm, 11)
				DebuglnIType("ORI", rd, rs1, imm)
				c.SetRegister(rd, c.GetRegister(rs1)|imm)
			case 0b111: // ------------------------------------------------------------------------ ANDI
				imm = SignExtend(imm, 11)
				DebuglnIType("ANDI", rd, rs1, imm)
				c.SetRegister(rd, c.GetRegister(rs1)&imm)
			case 0b001: // ------------------------------------------------------------------------ SLLI
				imm = InstructionPart(imm, 0, 5)
				DebuglnIType("SLLI", rd, rs1, imm)
				c.SetRegister(rd, c.GetRegister(rs1)<<imm)
			case 0b101:
				switch InstructionPart(s, 26, 31) {
				case 0b000000: // ----------------------------------------------------------------- SRLI
					imm = InstructionPart(imm, 0, 5)
					DebuglnIType("SRLI", rd, rs1, imm)
					c.SetRegister(rd, c.GetRegister(rs1)>>imm)
				case 0b010000: // ----------------------------------------------------------------- SRAI
					imm = InstructionPart(imm, 0, 5)
					DebuglnIType("SRAI", rd, rs1, imm)
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>imm))
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
					DebuglnRType("ADD", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)+c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MUL
					DebuglnRType("MUL", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))*int64(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SUB
					DebuglnRType("SUB", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)-c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b001:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SLL
					DebuglnRType("SLL", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)<<InstructionPart(c.GetRegister(rs2), 0, 5))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULH
					DebuglnRType("MULH", rd, rs1, rs2)
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
					DebuglnRType("SLT", rd, rs1, rs2)
					if int64(c.GetRegister(rs1)) < int64(c.GetRegister(rs2)) {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULHSU
					DebuglnRType("MULHSU", rd, rs1, rs2)
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
					DebuglnRType("SLTU", rd, rs1, rs2)
					if c.GetRegister(rs1) < c.GetRegister(rs2) {
						c.SetRegister(rd, 1)
					} else {
						c.SetRegister(rd, 0)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULHU
					DebuglnRType("MULHU", rd, rs1, rs2)
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
					DebuglnRType("XOR", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)^c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIV
					DebuglnRType("DIV", rd, rs1, rs2)
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
					DebuglnRType("SRL", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)>>InstructionPart(c.GetRegister(rs2), 0, 5))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIVU
					DebuglnRType("DIVU", rd, rs1, rs2)
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, math.MaxUint64)
					} else {
						c.SetRegister(rd, c.GetRegister(rs1)/c.GetRegister(rs2))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRA
					DebuglnRType("SRA", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int64(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 5)))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b110:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- OR
					DebuglnRType("OR", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)|c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- REM
					DebuglnRType("REM", rd, rs1, rs2)
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
					DebuglnRType("AND", rd, rs1, rs2)
					c.SetRegister(rd, c.GetRegister(rs1)&c.GetRegister(rs2))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- REMU
					DebuglnRType("REMU", rd, rs1, rs2)
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
				Debugln(fmt.Sprintf("Instr: % 10s |", "FENCE"))
			case 0b001: // ---------------------------------------------------------------------- FENCE.I
				Debugln(fmt.Sprintf("Instr: % 10s |", "FENCE.I"))
			}
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b1110011:
			rd, rs1, imm := IType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000:
				switch InstructionPart(s, 20, 31) {
				case 0b000000000000: // ----------------------------------------------------------- ECALL
					DebuglnIType("ECALL", rd, rs1, imm)
					return c.GetSystem().HandleCall(c)
				case 0b000000000001: // ----------------------------------------------------------- EBREAK
					DebuglnIType("EBREAK", rd, rs1, imm)
					return 1, nil
				}
			case 0b001: // ------------------------------------------------------------------------ CSRRW
				DebuglnIType("CSRRW", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, c.GetRegister(rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b010: // ------------------------------------------------------------------------ CSRRS
				DebuglnIType("CSRRS", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, c.GetCSR(imm)|c.GetRegister(rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ CSRRC
				DebuglnIType("CSRRC", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, c.GetCSR(imm)&(math.MaxUint64-c.GetRegister(rs1)))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101: // ------------------------------------------------------------------------ CSRRWI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRWI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, rs1)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b110: // ------------------------------------------------------------------------ CSRRSI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRSI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, c.GetCSR(imm)|rs1)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b111: // ------------------------------------------------------------------------ CSRRCI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRCI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR(imm))
				c.SetCSR(imm, c.GetCSR(imm)&(math.MaxUint64-rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		}
	}
	return 0, nil
}
