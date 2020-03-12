package rv64

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
		case 0b011_0111: // ---------------------------------------------------------------------- LUI
			rd, imm := UType(s)
			imm = SignExtend(imm, 31)
			DebuglnUType("LUI", rd, imm)
			c.SetRegister(rd, imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b001_0111: // ---------------------------------------------------------------------- AUIPC
			rd, imm := UType(s)
			imm = SignExtend(imm, 31)
			DebuglnUType("AUIPC", rd, imm)
			c.SetRegister(rd, c.GetPC()+imm)
			c.SetPC(c.GetPC() + 4)
			return 1, nil
		case 0b110_1111: // ---------------------------------------------------------------------- JAL
			rd, imm := JType(s)
			imm = SignExtend(imm, 19)
			DebuglnJType("JAL", rd, imm)
			c.SetRegister(rd, c.GetPC()+4)
			c.SetPC(c.GetPC() + imm)
			return 1, nil
		case 0b110_0111: // ---------------------------------------------------------------------- JALR
			rd, rs1, imm := IType(s)
			imm = SignExtend(imm, 11)
			DebuglnIType("JALR", rd, rs1, imm)
			c.SetRegister(rd, c.GetPC()+4)
			c.SetPC(((c.GetRegister(rs1) + imm) >> 1) << 1)
			return 1, nil
		case 0b110_0011:
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
		case 0b000_0011:
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
		}
	}
	return 0, nil
}
