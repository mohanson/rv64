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
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, c.GetRegister(rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b010: // ------------------------------------------------------------------------ CSRRS
				DebuglnIType("CSRRS", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, c.GetCSR().Get(imm)|c.GetRegister(rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ CSRRC
				DebuglnIType("CSRRC", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, c.GetCSR().Get(imm)&(math.MaxUint64-c.GetRegister(rs1)))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101: // ------------------------------------------------------------------------ CSRRWI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRWI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, rs1)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b110: // ------------------------------------------------------------------------ CSRRSI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRSI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, c.GetCSR().Get(imm)|rs1)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b111: // ------------------------------------------------------------------------ CSRRCI
				rs1 = SignExtend(rs1, 4)
				DebuglnIType("CSRRCI", rd, rs1, imm)
				c.SetRegister(rd, c.GetCSR().Get(imm))
				c.GetCSR().Set(imm, c.GetCSR().Get(imm)&(math.MaxUint64-rs1))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b0011011:
			rd, rs1, imm := IType(s)
			switch InstructionPart(s, 12, 14) {
			case 0b000: // ------------------------------------------------------------------------ ADDIW
				imm = SignExtend(imm, 11)
				DebuglnIType("ADDIW", rd, rs1, imm)
				c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(imm)))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b001: // ------------------------------------------------------------------------ SLLIW
				if InstructionPart(imm, 5, 5) != 0x00 {
					return 0, ErrAbnormalInstruction
				}
				imm = InstructionPart(imm, 0, 4)
				DebuglnIType("SLLIW", rd, rs1, imm)
				c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<imm), 31))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b101:
				switch InstructionPart(s, 25, 31) {
				case 0b0000000: // ---------------------------------------------------------------- SRLIW
					if InstructionPart(imm, 5, 5) != 0x00 {
						return 0, ErrAbnormalInstruction
					}
					imm = InstructionPart(imm, 0, 4)
					DebuglnIType("SRLIW", rd, rs1, imm)
					c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>imm), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRAIW
					if InstructionPart(imm, 5, 5) != 0x00 {
						return 0, ErrAbnormalInstruction
					}
					imm = InstructionPart(imm, 0, 4)
					DebuglnIType("SRAIW", rd, rs1, imm)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>imm))
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
					DebuglnRType("ADDW", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))+int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- MULW
					DebuglnRType("MULW", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))*int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SUBW
					DebuglnRType("SUBW", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))-int32(c.GetRegister(rs2))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b001: // ------------------------------------------------------------------------ SLLW
				DebuglnRType("SLLW", rd, rs1, rs2)
				c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))<<InstructionPart(c.GetRegister(rs2), 0, 4)), 31))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b100: // ------------------------------------------------------------------------ DIVW
				DebuglnRType("DIVW", rd, rs1, rs2)
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
					DebuglnRType("SRLW", rd, rs1, rs2)
					c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0000001: // ---------------------------------------------------------------- DIVUW
					DebuglnRType("DIVUW", rd, rs1, rs2)
					if c.GetRegister(rs2) == 0 {
						c.SetRegister(rd, math.MaxUint64)
					} else {
						c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegister(rs1))/uint32(c.GetRegister(rs2))), 31))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b0100000: // ---------------------------------------------------------------- SRAW
					DebuglnRType("SRAW", rd, rs1, rs2)
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))>>InstructionPart(c.GetRegister(rs2), 0, 4)))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b110: // ------------------------------------------------------------------------ REMW
				DebuglnRType("REMW", rd, rs1, rs2)
				if c.GetRegister(rs2) == 0 {
					c.SetRegister(rd, c.GetRegister(rs1))
				} else {
					c.SetRegister(rd, uint64(int32(c.GetRegister(rs1))%int32(c.GetRegister(rs2))))
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b111: // ------------------------------------------------------------------------ REMUW
				DebuglnRType("REMUW", rd, rs1, rs2)
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
				switch InstructionPart(s, 27, 31) {
				case 0b00010: // ------------------------------------------------------------------ LR.W
					DebuglnRType("LR.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ SC.W
					DebuglnRType("SC.W", rd, rs1, rs2)
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
				case 0b00001: // ------------------------------------------------------------------ AMOSWAP.W
					DebuglnRType("AMOSWAP.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000: // ------------------------------------------------------------------ AMOADD.W
					DebuglnRType("AMOADD.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v+uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100: // ------------------------------------------------------------------ AMOXOR.W
					DebuglnRType("AMOXOR.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v^uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100: // ------------------------------------------------------------------ AMOAND.W
					DebuglnRType("AMOAND.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v&uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000: // ------------------------------------------------------------------ AMOOR.W
					DebuglnRType("AMOOR.W", rd, rs1, rs2)
					a := SignExtend(c.GetRegister(rs1), 31)
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v|uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b10000: // ------------------------------------------------------------------ AMOMIN.W
					DebuglnRType("AMOMIN.W", rd, rs1, rs2)
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
				case 0b10100: // ------------------------------------------------------------------ AMOMAX.W
					DebuglnRType("AMOMAX.W", rd, rs1, rs2)
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
				case 0b11000: // ------------------------------------------------------------------ AMOMINU.W
					DebuglnRType("AMOMINU.W", rd, rs1, rs2)
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
				case 0b11100: // ------------------------------------------------------------------ AMOMAXU.W
					DebuglnRType("AMOMAXU.W", rd, rs1, rs2)
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
			case 0b011:
				switch InstructionPart(s, 27, 31) {
				case 0b00010: // ------------------------------------------------------------------ LR.D
					DebuglnRType("LR.D", rd, rs1, rs2)
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, v)
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ SC.D
					DebuglnRType("SC.D", rd, rs1, rs2)
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
				case 0b00001: // ------------------------------------------------------------------ AMOSWAP.D
					DebuglnRType("AMOSWAP.D", rd, rs1, rs2)
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000: // ------------------------------------------------------------------ AMOADD.D
					DebuglnRType("AMOADD.D", rd, rs1, rs2)
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v+c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100: // ------------------------------------------------------------------ AMOXOR.D
					DebuglnRType("AMOXOR.D", rd, rs1, rs2)
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v^c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100: // ------------------------------------------------------------------ AMOAND.D
					DebuglnRType("AMOAND.D", rd, rs1, rs2)
					a := c.GetRegister(rs1)
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v&c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000: // ------------------------------------------------------------------ AMOOR.D
					DebuglnRType("AMOOR.D", rd, rs1, rs2)
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
					DebuglnRType("AMOMIN.D", rd, rs1, rs2)
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
				case 0b10100: // ------------------------------------------------------------------ AMOMAX.D
					DebuglnRType("AMOMAX.D", rd, rs1, rs2)
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
				case 0b11000: // ------------------------------------------------------------------ AMOMINU.D
					DebuglnRType("AMOMINU.D", rd, rs1, rs2)
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
				case 0b11100: // ------------------------------------------------------------------ AMOMAXU.D
					DebuglnRType("AMOMAXU.D", rd, rs1, rs2)
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
			}
		case 0b0000111:
			rd, rs1, imm := IType(s)
			imm = SignExtend(imm, 11)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(s, 12, 14) {
			case 0b010: // ------------------------------------------------------------------------ FLW
				DebuglnIType("FLW", rd, rs1, imm)
				v, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				c.SetRegisterFloatAsFloat32(rd, math.Float32frombits(v))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ FLD
				DebuglnIType("FLD", rd, rs1, imm)
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
			imm = SignExtend(imm, 11)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(s, 12, 14) {
			case 0b010: // ------------------------------------------------------------------------ FSW
				DebuglnSType("FSW", rs1, rs2, imm)
				err := c.GetMemory().SetUint32(a, uint32(c.GetRegisterFloat(rs2)))
				if err != nil {
					return 0, err
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011: // ------------------------------------------------------------------------ FSD
				DebuglnSType("FSD", rs1, rs2, imm)
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
				DebuglnR4Type("FMADD.S", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat32(rs1)
				b := c.GetRegisterFloatAsFLoat32(rs2)
				d := c.GetRegisterFloatAsFLoat32(rs3)
				r := a*b + d
				c.SetRegisterFloatAsFloat32(rd, r)
				if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
					c.SetFloatFlag(FFlagsNX, 1)
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b01: // ------------------------------------------------------------------------- FMADD.D
				DebuglnR4Type("FMADD.D", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat64(rs1)
				b := c.GetRegisterFloatAsFLoat64(rs2)
				d := c.GetRegisterFloatAsFLoat64(rs3)
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
				DebuglnR4Type("FMSUB.S", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat32(rs1)
				b := c.GetRegisterFloatAsFLoat32(rs2)
				d := c.GetRegisterFloatAsFLoat32(rs3)
				r := a*b - d
				c.SetRegisterFloatAsFloat32(rd, r)
				if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
					c.SetFloatFlag(FFlagsNX, 1)
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b01: // ------------------------------------------------------------------------- FMSUB.D
				DebuglnR4Type("FMSUB.D", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat64(rs1)
				b := c.GetRegisterFloatAsFLoat64(rs2)
				d := c.GetRegisterFloatAsFLoat64(rs3)
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
				DebuglnR4Type("FNMSUB.S", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat32(rs1)
				b := c.GetRegisterFloatAsFLoat32(rs2)
				d := c.GetRegisterFloatAsFLoat32(rs3)
				r := a*b - d
				c.SetRegisterFloatAsFloat32(rd, -r)
				if r+d != a*b || a*b-r != d || (r+d)/a != b || (r+d)/b != a {
					c.SetFloatFlag(FFlagsNX, 1)
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b01: // ------------------------------------------------------------------------- FNMSUB.D
				DebuglnR4Type("FNMSUB.D", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat64(rs1)
				b := c.GetRegisterFloatAsFLoat64(rs2)
				d := c.GetRegisterFloatAsFLoat64(rs3)
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
				DebuglnR4Type("FNMADD.S", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat32(rs1)
				b := c.GetRegisterFloatAsFLoat32(rs2)
				d := c.GetRegisterFloatAsFLoat32(rs3)
				r := a*b + d
				c.SetRegisterFloatAsFloat32(rd, -r)
				if r-d != a*b || r-a*b != d || (r-d)/a != b || (r-d)/b != a {
					c.SetFloatFlag(FFlagsNX, 1)
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b01: // ------------------------------------------------------------------------- FNMADD.D
				DebuglnR4Type("FNMADD.D", rd, rs1, rs2, rs3)
				c.ClrFloatFlag()
				a := c.GetRegisterFloatAsFLoat64(rs1)
				b := c.GetRegisterFloatAsFLoat64(rs2)
				d := c.GetRegisterFloatAsFLoat64(rs3)
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
				a := c.GetRegisterFloatAsFLoat32(rs1)
				b := c.GetRegisterFloatAsFLoat32(rs2)
				switch InstructionPart(s, 27, 31) {
				case 0b00000: // ------------------------------------------------------------------ FADD.S
					DebuglnRType("FADD.S", rd, rs1, rs2)
					c.ClrFloatFlag()
					d := a + b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d-a != b || d-b != a {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ FSUB.S
					DebuglnRType("FSUB.S", rd, rs1, rs2)
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
					DebuglnRType("FMUL.S", rd, rs1, rs2)
					c.ClrFloatFlag()
					d := a * b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d/a != b || d/b != a || float64(a)*float64(b) != float64(d) {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ FDIV.S
					DebuglnRType("FDIV.S", rd, rs1, rs2)
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
					DebuglnRType("FSQRT.D", rd, rs1, rs2)
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
						DebuglnRType("FSGNJ.S", rd, rs1, rs2)
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FSGNJN.S
						DebuglnRType("FSGNJ.S", rd, rs1, rs2)
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010: // ---------------------------------------------------------------- FSGNJX.S
						DebuglnRType("FSGNJ.S", rd, rs1, rs2)
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
						DebuglnRType("FMIN.D", rd, rs1, rs2)
					case 0b001: // ---------------------------------------------------------------- FMAX.D
						DebuglnRType("FMAX.D", rd, rs1, rs2)
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
						DebuglnRType("FCVT.W.S", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat32(rs1)
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
						DebuglnRType("FCVT.WU.S", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat32(rs1)
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
						DebuglnRType("FCVT.L.S", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat32(rs1)
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
						DebuglnRType("FCVT.LU.D", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat32(rs1)
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
					DebuglnRType("FCVT.S.D", rd, rs1, rs2)
					d := c.GetRegisterFloatAsFLoat64(rs1)
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
						DebuglnRType("FMV.X.W", rd, rs1, rs2)
						c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegisterFloat(rs1))), 31))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FCLASS.S
						DebuglnRType("FCLASS.S", rd, rs1, rs2)
						a := c.GetRegisterFloatAsFLoat32(rs1)
						c.SetRegister(rd, FClassS(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b10100:
					var cond bool
					switch InstructionPart(s, 12, 14) {
					case 0b010: // ---------------------------------------------------------------- FEQ.S
						DebuglnRType("FEQ.S", rd, rs1, rs2)
						if IsSNaN32(a) || IsSNaN32(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001: // ---------------------------------------------------------------- FLT.S
						DebuglnRType("FLT.S", rd, rs1, rs2)
						if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000: // ---------------------------------------------------------------- FLE.S
						DebuglnRType("FLE.S", rd, rs1, rs2)
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
						DebuglnRType("FCVT.S.W", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat32(rd, float32(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001: // -------------------------------------------------------------- FCVT.S.WU
						DebuglnRType("FCVT.S.WU", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat32(rd, float32(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010: // -------------------------------------------------------------- FCVT.S.L
						DebuglnRType("FCVT.S.L", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat32(rd, float32(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011: // -------------------------------------------------------------- FCVT.S.LU
						DebuglnRType("FCVT.S.LU", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat32(rd, float32(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110: // ------------------------------------------------------------------ FMV.W.X
					DebuglnRType("FMV.W.X", rd, rs1, rs2)
					c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(uint32(c.GetRegister(rs1))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b01:
				a := c.GetRegisterFloatAsFLoat64(rs1)
				b := c.GetRegisterFloatAsFLoat64(rs2)
				switch InstructionPart(s, 27, 31) {
				case 0b00000: // ------------------------------------------------------------------ FADD.D
					DebuglnRType("FADD.D", rd, rs1, rs2)
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a+b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001: // ------------------------------------------------------------------ FSUB.D
					DebuglnRType("FSUB.D", rd, rs1, rs2)
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
					DebuglnRType("FMUL.D", rd, rs1, rs2)
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a*b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011: // ------------------------------------------------------------------ FDIV.D
					DebuglnRType("FDIV.D", rd, rs1, rs2)
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
					DebuglnRType("FSQRT.D", rd, rs1, rs2)
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
						DebuglnRType("FSGNJ.D", rd, rs1, rs2)
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FSGNJN.D
						DebuglnRType("FSGNJN.D", rd, rs1, rs2)
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010: // ---------------------------------------------------------------- FSGNJX.D
						DebuglnRType("FSGNJX.D", rd, rs1, rs2)
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
						DebuglnRType("FMIN.D", rd, rs1, rs2)
					case 0b001: // ---------------------------------------------------------------- FMAX.D
						DebuglnRType("FMAX.D", rd, rs1, rs2)
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
						DebuglnRType("FCVT.W.D", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat64(rs1)
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
						DebuglnRType("FCVT.WU.D", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat64(rs1)
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
						DebuglnRType("FCVT.L.D", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat64(rs1)
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
						DebuglnRType("FCVT.LU.D", rd, rs1, rs2)
						d := c.GetRegisterFloatAsFLoat64(rs1)
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
					DebuglnRType("FCVT.D.S", rd, rs1, rs2)
					d := c.GetRegisterFloatAsFLoat32(rs1)
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
						DebuglnRType("FEQ.D", rd, rs1, rs2)
						if IsSNaN64(a) || IsSNaN64(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001: // ---------------------------------------------------------------- FLT.D
						DebuglnRType("FLT.D", rd, rs1, rs2)
						if math.IsNaN(a) || math.IsNaN(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000: // ---------------------------------------------------------------- FLE.D
						DebuglnRType("FLE.D", rd, rs1, rs2)
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
						DebuglnRType("FMV.X.D", rd, rs1, rs2)
						c.SetRegister(rd, c.GetRegisterFloat(rs1))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001: // ---------------------------------------------------------------- FCLASS.D
						DebuglnRType("FCLASS.D", rd, rs1, rs2)
						a := c.GetRegisterFloatAsFLoat64(rs1)
						c.SetRegister(rd, FClassD(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11010:
					switch InstructionPart(s, 20, 24) {
					case 0b00000: // -------------------------------------------------------------- FCVT.D.W
						DebuglnRType("FCVT.D.W", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat64(rd, float64(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001: // -------------------------------------------------------------- FCVT.D.WU
						DebuglnRType("FCVT.D.WU", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat64(rd, float64(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010: // -------------------------------------------------------------- FCVT.D.L
						DebuglnRType("FCVT.D.L", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat64(rd, float64(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011: // -------------------------------------------------------------- FCVT.D.LU
						DebuglnRType("FCVT.D.LU", rd, rs1, rs2)
						c.SetRegisterFloatAsFloat64(rd, float64(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110: // ------------------------------------------------------------------ FMV.D.X
					DebuglnRType("FMV.D.X", rd, rs1, rs2)
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
