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
	var i uint64 = 0
	for j := len(data) - 1; j >= 0; j-- {
		i += uint64(data[j]) << (8 * j)
	}
	switch len(data) {
	case 2:
		opcode := InstructionPart(i, 0, 1)
		funct3 := InstructionPart(i, 13, 15)
		switch opcode<<3 | funct3 {
		case 00_000:
			return aluC.addi4spn(c, i)
		case 0b00_001:
			return aluC.fld(c, i)
		case 0b00_010:
			return aluC.lw(c, i)
		case 0b00_011:
			return aluC.ld(c, i)
		case 0b00_100:
			return 0, ErrReservedInstruction
		case 0b00_101:
			return aluC.fsd(c, i)
		case 0b00_110:
			return aluC.sw(c, i)
		case 0b00_111:
			return aluC.sd(c, i)
		case 0b01_000:
			return aluC.addi(c, i)
		case 0b01_001:
			return aluC.addiw(c, i)
		case 0b01_010:
			return aluC.li(c, i)
		case 0b01_011:
			if InstructionPart(i, 7, 11) == Rsp {
				return aluC.addi16sp(c, i)
			} else {
				return aluC.lui(c, i)
			}
		case 0b01_100:
			// misc-alu
			switch InstructionPart(i, 10, 11) {
			case 0b00:
				return aluC.srli(c, i)
			case 0b01:
				return aluC.srai(c, i)
			case 0b10:
				return aluC.andi(c, i)
			case 0b11:
				switch InstructionPart(i, 12, 12)<<2 | InstructionPart(i, 5, 6) {
				case 0b0_00:
					return aluC.sub(c, i)
				case 0b0_01:
					return aluC.xor(c, i)
				case 0b0_10:
					return aluC.or(c, i)
				case 0b0_11:
					return aluC.and(c, i)
				case 0b1_00:
					return aluC.subw(c, i)
				case 0b1_01:
					return aluC.addw(c, i)
				case 0b1_10:
					return 0, ErrReservedInstruction
				case 0b1_11:
					return 0, ErrReservedInstruction
				}
			}
		case 0b01_101:
			return aluC.j(c, i)
		case 0b01_110:
			return aluC.beqz(c, i)
		case 0b01_111:
			return aluC.bnez(c, i)
		case 0b10_000:
			return aluC.slli(c, i)
		case 0b10_001:
			Println("c.fldsp")
		case 0b10_010:
			Println("c.lwsp")
		case 0b10_011:
			Println("c.ldsp")
		case 0b10_100:
			// j[al]r/mv/add
			switch InstructionPart(i, 12, 12) {
			case 0:
				if InstructionPart(i, 2, 6) == 0 {
					return aluC.jr(c, i)
				} else {
					return aluC.mv(c, i)
				}
			case 1:
				l1 := InstructionPart(i, 7, 11)
				l2 := InstructionPart(i, 2, 6)
				if l1 == 0 && l2 == 0 {
					Println("c.ebreak")
				}
				if l1 != 0 {
					return aluC.jalr(c, i)
				}
				if l1 == 0 {
					Println("c.add")
				}
			}
		case 0b10_101:
			Println("c.fsdsp")
		case 0b10_110:
			Println("c.swsp")
		case 0b10_111:
			Println("c.sdsp")
		}
		Panicln("unreachable")
	case 4:
		opcode := InstructionPart(i, 0, 6)
		funct3 := InstructionPart(i, 12, 14)
		funct6 := InstructionPart(i, 26, 31)
		funct7 := InstructionPart(i, 25, 31)
		switch opcode {
		case 0b0110111:
			return aluI.lui(c, i)
		case 0b0010111:
			return aluI.aupic(c, i)
		case 0b1101111:
			return aluI.jal(c, i)
		case 0b1100111:
			return aluI.jalr(c, i)
		case 0b1100011:
			switch funct3 {
			case 0b000:
				return aluI.beq(c, i)
			case 0b001:
				return aluI.bne(c, i)
			case 0b100:
				return aluI.blt(c, i)
			case 0b101:
				return aluI.bge(c, i)
			case 0b110:
				return aluI.bltu(c, i)
			case 0b111:
				return aluI.bgeu(c, i)
			}
		case 0b0000011:
			switch funct3 {
			case 0b000:
				return aluI.lb(c, i)
			case 0b001:
				return aluI.lh(c, i)
			case 0b010:
				return aluI.lw(c, i)
			case 0b011:
				return aluI.ld(c, i)
			case 0b100:
				return aluI.lbu(c, i)
			case 0b101:
				return aluI.lhu(c, i)
			case 0b110:
				return aluI.lwu(c, i)
			}
		case 0b0100011:
			switch funct3 {
			case 0b000:
				return aluI.sb(c, i)
			case 0b001:
				return aluI.sh(c, i)
			case 0b010:
				return aluI.sw(c, i)
			case 0b011:
				return aluI.sd(c, i)
			}
		case 0b0010011:
			switch funct3 {
			case 0b000:
				return aluI.addi(c, i)
			case 0b010:
				return aluI.slti(c, i)
			case 0b011:
				return aluI.sltiu(c, i)
			case 0b100:
				return aluI.xori(c, i)
			case 0b110:
				return aluI.ori(c, i)
			case 0b111:
				return aluI.andi(c, i)
			case 0b001:
				return aluI.slli(c, i)
			case 0b101:
				switch funct6 {
				case 0b000000:
					return aluI.srli(c, i)
				case 0b010000:
					return aluI.srai(c, i)
				}
			}
		case 0b0110011:
			switch funct3 {
			case 0b000:
				switch funct7 {
				case 0b0000000:
					return aluI.add(c, i)
				case 0b0000001:
					return aluM.mul(c, i)
				case 0b0100000:
					return aluI.sub(c, i)
				}
			case 0b001:
				switch funct7 {
				case 0b0000000:
					return aluI.sll(c, i)
				case 0b0000001:
					return aluM.mulh(c, i)
				}
			case 0b010:
				switch funct7 {
				case 0b0000000:
					return aluI.slt(c, i)
				case 0b0000001:
					return aluM.mulhsu(c, i)
				}
			case 0b011:
				switch funct7 {
				case 0b0000000:
					return aluI.sltu(c, i)
				case 0b0000001:
					return aluM.mulhu(c, i)
				}
			case 0b100:
				switch funct7 {
				case 0b0000000:
					return aluI.xor(c, i)
				case 0b0000001:
					return aluM.div(c, i)
				}
			case 0b101:
				switch funct7 {
				case 0b0000000:
					return aluI.srl(c, i)
				case 0b0000001:
					return aluM.divu(c, i)
				case 0b0100000:
					return aluI.sra(c, i)
				}
			case 0b110:
				switch funct7 {
				case 0b0000000:
					return aluI.or(c, i)
				case 0b0000001:
					return aluM.rem(c, i)
				}
			case 0b111:
				switch funct7 {
				case 0b0000000:
					return aluI.and(c, i)
				case 0b0000001:
					return aluM.remu(c, i)
				}
			}
		case 0b0001111:
			switch funct3 {
			case 0b000:
				return aluI.fence(c, i)
			case 0b001:
				return aluZifencei.fencei(c, i)
			}
		case 0b1110011:
			switch funct3 {
			case 0b000:
				switch InstructionPart(i, 20, 31) {
				case 0b000000000000:
					return aluI.ecall(c, i)
				case 0b000000000001:
					return aluI.ebreak(c, i)
				}
			case 0b001:
				return aluZicsr.csrrw(c, i)
			case 0b010:
				return aluZicsr.csrrs(c, i)
			case 0b011:
				return aluZicsr.csrrc(c, i)
			case 0b101:
				return aluZicsr.csrrwi(c, i)
			case 0b110:
				return aluZicsr.csrrsi(c, i)
			case 0b111:
				return aluZicsr.csrrci(c, i)
			}
		case 0b0011011:
			switch funct3 {
			case 0b000:
				return aluI.addiw(c, i)
			case 0b001:
				return aluI.slliw(c, i)
			case 0b101:
				switch funct7 {
				case 0b0000000:
					return aluI.srliw(c, i)
				case 0b0100000:
					return aluI.sraiw(c, i)
				}
			}
		case 0b0111011:
			switch funct3 {
			case 0b000:
				switch funct7 {
				case 0b0000000:
					return aluI.addw(c, i)
				case 0b0000001:
					return aluM.mulw(c, i)
				case 0b0100000:
					return aluI.subw(c, i)
				}
			case 0b001:
				return aluI.sllw(c, i)
			case 0b100:
				return aluM.divw(c, i)
			case 0b101:
				switch funct7 {
				case 0b0000000:
					return aluI.srlw(c, i)
				case 0b0000001:
					return aluM.divuw(c, i)
				case 0b0100000:
					return aluI.sraw(c, i)
				}
			case 0b110:
				return aluM.remw(c, i)
			case 0b111:
				return aluM.remuw(c, i)
			}
		case 0b0101111:
			rd, rs1, rs2 := RType(i)
			switch InstructionPart(i, 12, 14) {
			case 0b010:
				a := SignExtend(c.GetRegister(rs1), 31)
				switch InstructionPart(i, 27, 31) {
				case 0b00010:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "lr.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sc.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if a == c.GetLoadReservation() {
						c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
						c.SetRegister(rd, 0)
					} else {
						c.SetRegister(rd, 1)
					}
					c.SetLoadReservation(0)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoswap.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoadd.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v+uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoxor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v^uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoand.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v&uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoor.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint32(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint32(a, v|uint32(c.GetRegister(rs2)))
					c.SetRegister(rd, SignExtend(uint64(v), 31))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b10000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomin.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b10100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomax.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amominu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomaxu.w", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				switch InstructionPart(i, 27, 31) {
				case 0b00010:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "lr.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.SetRegister(rd, v)
					c.SetLoadReservation(a)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "sc.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					if a == c.GetLoadReservation() {
						c.GetMemory().SetUint64(a, c.GetRegister(rs2))
						c.SetRegister(rd, 0)
					} else {
						c.SetRegister(rd, 1)
					}
					c.SetLoadReservation(0)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoswap.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoadd.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v+c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoxor.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v^c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amoand.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
					v, err := c.GetMemory().GetUint64(a)
					if err != nil {
						return 0, err
					}
					c.GetMemory().SetUint64(a, v&c.GetRegister(rs2))
					c.SetRegister(rd, v)
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b01000:
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
				case 0b10000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomin.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b10100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomax.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amominu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
				case 0b11100:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "amomaxu.d", c.LogI(rd), c.LogI(rs1), c.LogI(rs2)))
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
			rd, rs1, imm := IType(i)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(i, 12, 14) {
			case 0b010:
				Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "flw", c.LogF(rd), c.LogI(rs1), imm))
				v, err := c.GetMemory().GetUint32(a)
				if err != nil {
					return 0, err
				}
				c.SetRegisterFloatAsFloat32(rd, math.Float32frombits(v))
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011:
				Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s imm: ____(%#016x)", c.GetPC(), "fld", c.LogF(rd), c.LogI(rs1), imm))
				v, err := c.GetMemory().GetUint64(a)
				if err != nil {
					return 0, err
				}
				c.SetRegisterFloat(rd, v)
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b0100111:
			rs1, rs2, imm := SType(i)
			a := c.GetRegister(rs1) + imm
			switch InstructionPart(i, 12, 14) {
			case 0b010:
				Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "fsw", c.LogI(rs1), c.LogF(rs2), imm))
				err := c.GetMemory().SetUint32(a, uint32(c.GetRegisterFloat(rs2)))
				if err != nil {
					return 0, err
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			case 0b011:
				Debugln(fmt.Sprintf("%#08x % 10s rs1: %s rs2: %s imm: ____(%#016x)", c.GetPC(), "fsd", c.LogI(rs1), c.LogF(rs2), imm))
				err := c.GetMemory().SetUint64(a, c.GetRegisterFloat(rs2))
				if err != nil {
					return 0, err
				}
				c.SetPC(c.GetPC() + 4)
				return 1, nil
			}
		case 0b1000011:
			rd, rs1, rs2, rs3 := R4Type(i)
			switch InstructionPart(i, 25, 26) {
			case 0b00:
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
			case 0b01:
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
		case 0b1000111:
			rd, rs1, rs2, rs3 := R4Type(i)
			switch InstructionPart(i, 25, 26) {
			case 0b00:
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
			case 0b01:
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
		case 0b1001011:
			rd, rs1, rs2, rs3 := R4Type(i)
			switch InstructionPart(i, 25, 26) {
			case 0b00:
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
			case 0b01:
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
		case 0b1001111:
			rd, rs1, rs2, rs3 := R4Type(i)
			switch InstructionPart(i, 25, 26) {
			case 0b00:
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
			case 0b01:
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
		case 0b1010011:
			rd, rs1, rs2 := RType(i)
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				a := c.GetRegisterFloatAsFloat32(rs1)
				b := c.GetRegisterFloatAsFloat32(rs2)
				switch InstructionPart(i, 27, 31) {
				case 0b00000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fadd.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					d := a + b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d-a != b || d-b != a {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsub.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b00010:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmul.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					d := a * b
					c.SetRegisterFloatAsFloat32(rd, d)
					if d/a != b || d/b != a || float64(a)*float64(b) != float64(d) {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fdiv.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b01011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsqrt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnj.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjn.s.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjx.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(float64(a)) != math.Signbit(float64(b)) {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)|0x80000000))
						} else {
							c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(math.Float32bits(a)&0x7fffffff))
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b00101:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmin.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmax.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
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
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
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
					case 0b00001:
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
					case 0b00010:
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
					case 0b00011:
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
				case 0b01000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					d := c.GetRegisterFloatAsFloat64(rs1)
					if math.IsNaN(d) {
						c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(NaN32))
					} else {
						c.SetRegisterFloatAsFloat32(rd, float32(d))
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b11100:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.x.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegister(rd, SignExtend(uint64(uint32(c.GetRegisterFloat(rs1))), 31))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fclass.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						a := c.GetRegisterFloatAsFloat32(rs1)
						c.SetRegister(rd, FClassS(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b10100:
					var cond bool
					switch InstructionPart(i, 12, 14) {
					case 0b010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "feq.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if IsSNaN32(a) || IsSNaN32(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "flt.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fle.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.wu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.l", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.s.lu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat32(rd, float32(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.w.x", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.SetRegisterFloat(rd, 0xffffffff00000000|uint64(uint32(c.GetRegister(rs1))))
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				}
			case 0b01:
				a := c.GetRegisterFloatAsFloat64(rs1)
				b := c.GetRegisterFloatAsFloat64(rs2)
				switch InstructionPart(i, 27, 31) {
				case 0b00000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fadd.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a+b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00001:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsub.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b00010:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmul.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					c.ClrFloatFlag()
					c.SetRegisterFloatAsFloat64(rd, a*b)
					if big.NewFloat(0).Add(big.NewFloat(a), big.NewFloat(b)).Acc() != big.Exact {
						c.SetFloatFlag(FFlagsNX, 1)
					}
					c.SetPC(c.GetPC() + 4)
					return 1, nil
				case 0b00011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fdiv.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b01011:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsqrt.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnj.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjn.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fsgnjx.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.Signbit(a) != math.Signbit(b) {
							c.SetRegisterFloat(rd, math.Float64bits(a)|0x8000000000000000)
						} else {
							c.SetRegisterFloat(rd, math.Float64bits(a)&0x7fffffffffffffff)
						}
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b00101:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmin.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmax.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
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
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.w.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.wu.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.l.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					case 0b00011:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.lu.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
				case 0b01000:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.d.s", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
					case 0b010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "feq.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if IsSNaN64(a) || IsSNaN64(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a == b
						}
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "flt.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						if math.IsNaN(a) || math.IsNaN(b) {
							c.SetFloatFlag(FFlagsNV, 1)
						} else {
							cond = a < b
						}
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fle.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.x.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegister(rd, c.GetRegisterFloat(rs1))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fclass.d", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						a := c.GetRegisterFloatAsFloat64(rs1)
						c.SetRegister(rd, FClassD(a))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11010:
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.d.w", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(int32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00001:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.d.wu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(uint32(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00010:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.d.l", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(int64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					case 0b00011:
						Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fcvt.d.lu", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
						c.SetRegisterFloatAsFloat64(rd, float64(uint64(c.GetRegister(rs1))))
						c.SetPC(c.GetPC() + 4)
						return 1, nil
					}
				case 0b11110:
					Debugln(fmt.Sprintf("%#08x % 10s  rd: %s rs1: %s rs2: %s", c.GetPC(), "fmv.d.x", c.LogF(rd), c.LogF(rs1), c.LogF(rs2)))
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
