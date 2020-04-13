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
			return aluC.fldsp(c, i)
		case 0b10_010:
			return aluC.lwsp(c, i)
		case 0b10_011:
			return aluC.ldsp(c, i)
		case 0b10_100:
			switch InstructionPart(i, 12, 12) {
			case 0:
				if InstructionPart(i, 2, 6) == Rzero {
					return aluC.jr(c, i)
				} else {
					return aluC.mv(c, i)
				}
			case 1:
				rs1 := InstructionPart(i, 7, 11)
				rs2 := InstructionPart(i, 2, 6)
				if rs2 != Rzero {
					return aluC.add(c, i)
				}
				if rs1 != Rzero {
					return aluC.jalr(c, i)
				}
				return aluC.ebreak(c, i)
			}
		case 0b10_101:
			return aluC.fsdsp(c, i)
		case 0b10_110:
			return aluC.swsp(c, i)
		case 0b10_111:
			return aluC.sdsp(c, i)
		}
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
			switch funct3 {
			case 0b010:
				switch InstructionPart(i, 27, 31) {
				case 0b00010:
					return aluA.lrw(c, i)
				case 0b00011:
					return aluA.scw(c, i)
				case 0b00001:
					return aluA.amoswapw(c, i)
				case 0b00000:
					return aluA.amoaddw(c, i)
				case 0b00100:
					return aluA.amoxorw(c, i)
				case 0b01100:
					return aluA.amoandw(c, i)
				case 0b01000:
					return aluA.amoorw(c, i)
				case 0b10000:
					return aluA.amominw(c, i)
				case 0b10100:
					return aluA.amomaxw(c, i)
				case 0b11000:
					return aluA.amominuw(c, i)
				case 0b11100:
					return aluA.amomaxuw(c, i)
				}
			case 0b011:
				switch InstructionPart(i, 27, 31) {
				case 0b00010:
					return aluA.lrd(c, i)
				case 0b00011:
					return aluA.scd(c, i)
				case 0b00001:
					return aluA.amoswapd(c, i)
				case 0b00000:
					return aluA.amoaddd(c, i)
				case 0b00100:
					return aluA.amoxord(c, i)
				case 0b01100:
					return aluA.amoandd(c, i)
				case 0b01000:
					return aluA.amoord(c, i)
				case 0b10000:
					return aluA.amomind(c, i)
				case 0b10100:
					return aluA.amomaxd(c, i)
				case 0b11000:
					return aluA.amominud(c, i)
				case 0b11100:
					return aluA.amomaxud(c, i)
				}
			}
		case 0b0000111:
			switch funct3 {
			case 0b010:
				return aluF.flw(c, i)
			case 0b011:
				return aluD.fld(c, i)
			}
		case 0b0100111:
			switch funct3 {
			case 0b010:
				return aluF.fsw(c, i)
			case 0b011:
				return aluD.fsd(c, i)
			}
		case 0b1000011:
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				return aluF.fmadds(c, i)
			case 0b01:
				return aluD.fmaddd(c, i)
			}
		case 0b1000111:
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				return aluF.fmsubs(c, i)
			case 0b01:
				return aluD.fmsubd(c, i)
			}
		case 0b1001011:
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				return aluF.fnmsubs(c, i)
			case 0b01:
				return aluD.fnmsubd(c, i)
			}
		case 0b1001111:
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				return aluF.fnmadds(c, i)
			case 0b01:
				return aluD.fnmaddd(c, i)
			}
		case 0b1010011:
			switch InstructionPart(i, 25, 26) {
			case 0b00:
				switch InstructionPart(i, 27, 31) {
				case 0b00000:
					return aluF.fadds(c, i)
				case 0b00001:
					return aluF.fsubs(c, i)
				case 0b00010:
					return aluF.fmuls(c, i)
				case 0b00011:
					return aluF.fdivs(c, i)
				case 0b01011:
					return aluF.fsqrts(c, i)
				case 0b00100:
					switch funct3 {
					case 0b000:
						return aluF.fsgnjs(c, i)
					case 0b001:
						return aluF.fsgnjns(c, i)
					case 0b010:
						return aluF.fsgnjxs(c, i)
					}
				case 0b00101:
					switch funct3 {
					case 0b000:
						return aluF.fmins(c, i)
					case 0b001:
						return aluF.fmaxs(c, i)
					}
				case 0b11000:
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						return aluF.fcvtws(c, i)
					case 0b00001:
						return aluF.fcvtwus(c, i)
					case 0b00010:
						return aluF.fcvtls(c, i)
					case 0b00011:
						return aluF.fcvtlus(c, i)
					}
				case 0b01000:
					return aluD.fcvtsd(c, i)
				case 0b11100:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						return aluF.fmvxw(c, i)
					case 0b001:
						return aluF.fclasss(c, i)
					}
				case 0b10100:
					switch InstructionPart(i, 12, 14) {
					case 0b010:
						return aluF.feqs(c, i)
					case 0b001:
						return aluF.flts(c, i)
					case 0b000:
						return aluF.fles(c, i)
					}
				case 0b11010:
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						return aluF.fcvtsw(c, i)
					case 0b00001:
						return aluF.fcvtswu(c, i)
					case 0b00010:
						return aluF.fcvtsl(c, i)
					case 0b00011:
						return aluF.fcvtslu(c, i)
					}
				case 0b11110:
					return aluF.fmvwx(c, i)
				}
			case 0b01:
				switch InstructionPart(i, 27, 31) {
				case 0b00000:
					return aluD.faddd(c, i)
				case 0b00001:
					return aluD.fsubd(c, i)
				case 0b00010:
					return aluD.fmuld(c, i)
				case 0b00011:
					return aluD.fdivd(c, i)
				case 0b01011:
					return aluD.fsqrtd(c, i)
				case 0b00100:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						return aluD.fsgnjd(c, i)
					case 0b001:
						return aluD.fsgnjnd(c, i)
					case 0b010:
						return aluD.fsgnjxd(c, i)
					}
				case 0b00101:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						return aluD.fmind(c, i)
					case 0b001:
						return aluD.fmaxd(c, i)
					}
				case 0b11000:
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						return aluD.fcvtwd(c, i)
					case 0b00001:
						return aluD.fcvtwud(c, i)
					case 0b00010:
						return aluD.fcvtld(c, i)
					case 0b00011:
						return aluD.fcvtlud(c, i)
					}
				case 0b01000:
					return aluD.fcvtds(c, i)
				case 0b10100:
					switch InstructionPart(i, 12, 14) {
					case 0b010:
						return aluD.feqd(c, i)
					case 0b001:
						return aluD.fltd(c, i)
					case 0b000:
						return aluD.fled(c, i)
					}
				case 0b11100:
					switch InstructionPart(i, 12, 14) {
					case 0b000:
						return aluD.fmvxd(c, i)
					case 0b001:
						return aluD.fclassd(c, i)
					}
				case 0b11010:
					switch InstructionPart(i, 20, 24) {
					case 0b00000:
						return aluD.fcvtdw(c, i)
					case 0b00001:
						return aluD.fcvtdwu(c, i)
					case 0b00010:
						return aluD.fcvtdl(c, i)
					case 0b00011:
						return aluD.fcvtdlu(c, i)
					}
				case 0b11110:
					return aluD.fmvdx(c, i)
				}
			}
		}
	}
	return 0, ErrAbnormalInstruction
}
