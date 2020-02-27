package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"log"

	"github.com/mohanson/riscv"
)

const cDebug = 1

type CPU struct {
	ModuleBase *riscv.RegisterRV64I
	Mem        []byte
}

func (c *CPU) pushString(s string) {
	bs := append([]byte(s), 0x00)
	c.ModuleBase.RG[riscv.Rsp] -= uint64(len(bs))
	for i, b := range bs {
		c.Mem[c.ModuleBase.RG[riscv.Rsp]+uint64(i)] = b
	}
}

func (c *CPU) pushUint64(v uint64) {
	c.ModuleBase.RG[riscv.Rsp] -= 8
	binary.LittleEndian.PutUint64(c.Mem[c.ModuleBase.RG[riscv.Rsp]:c.ModuleBase.RG[riscv.Rsp]+8], v)
}

func (c *CPU) FetchInstruction() []byte {
	if (c.ModuleBase.PC + 2) > uint64(len(c.Mem)) {
		log.Panicln("Out of memory")
	}
	a := c.Mem[c.ModuleBase.PC : c.ModuleBase.PC+2]
	b := riscv.InstructionLengthEncoding(a)
	instructionBytes := c.Mem[c.ModuleBase.PC : c.ModuleBase.PC+uint64(b)]
	return instructionBytes
}

var cStep = flag.Int64("steps", 20, "")

func (c *CPU) Run() {
	flag.Parse()
	i := 0
	for {
		if i > int(*cStep) {
			break
		}
		data := c.FetchInstruction()
		log.Println("==========")
		if len(data) == 2 {
			log.Printf("%08b %08b\n", data[1], data[0])
		} else if len(data) == 4 {
			log.Printf("%08b %08b %08b %08b\n", data[3], data[2], data[1], data[0])
		} else {
			log.Panicln("")
		}
		var s uint64 = 0
		for i := 0; i < 32; i++ {
			s += c.ModuleBase.RG[i]
		}
		log.Println(i, c.ModuleBase.PC, s)
		if len(data) == 4 {
			var s uint64 = 0
			for i := len(data) - 1; i >= 0; i-- {
				s += uint64(data[i]) << (8 * i)
			}

			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[0], c.ModuleBase.RG[1], c.ModuleBase.RG[2], c.ModuleBase.RG[3])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[4], c.ModuleBase.RG[5], c.ModuleBase.RG[6], c.ModuleBase.RG[7])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[8], c.ModuleBase.RG[9], c.ModuleBase.RG[10], c.ModuleBase.RG[11])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[12], c.ModuleBase.RG[13], c.ModuleBase.RG[14], c.ModuleBase.RG[15])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[16], c.ModuleBase.RG[17], c.ModuleBase.RG[18], c.ModuleBase.RG[19])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[20], c.ModuleBase.RG[21], c.ModuleBase.RG[22], c.ModuleBase.RG[23])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[24], c.ModuleBase.RG[25], c.ModuleBase.RG[26], c.ModuleBase.RG[27])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[28], c.ModuleBase.RG[29], c.ModuleBase.RG[30], c.ModuleBase.RG[31])
			if riscv.ExecuterRV64I(c.ModuleBase, s) != 0 {
				i += 1
				continue
			}
		}

		s = 0
		for i := len(data) - 1; i >= 0; i-- {
			s += uint64(data[i]) << (8 * i)
		}
		if riscv.ExecuterC(c.ModuleBase, c.Mem, s) != 0 {
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[0], c.ModuleBase.RG[1], c.ModuleBase.RG[2], c.ModuleBase.RG[3])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[4], c.ModuleBase.RG[5], c.ModuleBase.RG[6], c.ModuleBase.RG[7])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[8], c.ModuleBase.RG[9], c.ModuleBase.RG[10], c.ModuleBase.RG[11])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[12], c.ModuleBase.RG[13], c.ModuleBase.RG[14], c.ModuleBase.RG[15])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[16], c.ModuleBase.RG[17], c.ModuleBase.RG[18], c.ModuleBase.RG[19])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[20], c.ModuleBase.RG[21], c.ModuleBase.RG[22], c.ModuleBase.RG[23])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[24], c.ModuleBase.RG[25], c.ModuleBase.RG[26], c.ModuleBase.RG[27])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[28], c.ModuleBase.RG[29], c.ModuleBase.RG[30], c.ModuleBase.RG[31])
			i += 1
			continue
		}

		// switch len(data) {
		// case 2:
		// 	switch data[0] & 0x03 {
		// 	// case 00:
		// 	// Illegal instruction
		// 	// C.ADDI4SPN
		// 	// C.FLD
		// 	// C.LQ
		// 	// C.LW
		// 	// C.FLW
		// 	// C.LD
		// 	// Reserved
		// 	// C.FSD
		// 	// C.SQ
		// 	// C.SW
		// 	// C.FSW
		// 	// C.SD
		// 	// case 01:
		// 	// C.NOP
		// 	// C.ADDI
		// 	// C.JAL
		// 	// C.ADDIW
		// 	// C.LI
		// 	// C.ADDI16SP
		// 	// C.LUI
		// 	// C.SRLI
		// 	// C.SRLI64
		// 	// C.SRAI
		// 	// C.SRAI64
		// 	// C.ANDI
		// 	// C.SUB
		// 	// C.XOR
		// 	// C.OR
		// 	// C.AND
		// 	// C.SUBW
		// 	// C.ADDW
		// 	// C.J
		// 	// C.BEQZ
		// 	// C.BNEZ
		// 	// case 02:
		// 	// C.SLLI
		// 	// C.SLLI64
		// 	// C.FLDSP
		// 	// C.LQSP
		// 	// C.LWSP
		// 	// C.FLWSP
		// 	// C.LDSP
		// 	// C.JR
		// 	// C.MV
		// 	// C.EBREAK
		// 	// C.JALR
		// 	// C.ADD
		// 	// C.FSDSP
		// 	// C.SQSP
		// 	// C.SWSP
		// 	// C.FSWSP
		// 	// C.SDSP
		// 	default:
		// 		log.Panicln("")
		// 	}
		// case 4:
		// 	switch data[0] & 0b01111111 {
		// 	// RV32I Base Instruction Set
		// 	// case 0b00110111: // imm[31:12] rd 0110111 LUI
		// 	// AUIPC
		// 	// case 0b00010111:
		// 	// 	_, rd, imm := riscv.UType(data)
		// 	// 	riscv.PrintlnUType("AUIPC", rd, imm)
		// 	// 	c.ModuleBase.RG[rd] = c.ModuleBase.PC + uint64(imm)
		// 	// case 0b01101111: // imm[20|10:1|11|19:12] rd 1101111 JAL
		// 	// case 0b01100111: // imm[11:0] rs1 000 rd 1100111 JALR
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU
		// 	// case 0b01100011: // imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU
		// 	// case 0b00000011: // imm[11:0] rs1 000 rd 0000011 LB
		// 	// case 0b00000011: // imm[11:0] rs1 001 rd 0000011 LH
		// 	// case 0b00000011: // imm[11:0] rs1 010 rd 0000011 LW
		// 	// case 0b00000011: // imm[11:0] rs1 100 rd 0000011 LBU
		// 	// case 0b00000011: // imm[11:0] rs1 101 rd 0000011 LHU
		// 	// case 0b00100011: // imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB
		// 	// case 0b00100011: // imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH
		// 	// case 0b00100011: // imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW
		// 	// ADDI
		// 	// case 0b00010011:
		// 	// 	_, rd, funct3, rs1, imm := riscv.IType(data)
		// 	// 	riscv.PrintlnIType("ADDI", rd, funct3, rs1, imm)
		// 	// 	c.ModuleBase.RG[rd] = c.ModuleBase.RG[rs1] + uint64(imm)
		// 	// case 0b00010011: // imm[11:0] rs1 010 rd 0010011 SLTI
		// 	// case 0b00010011: // imm[11:0] rs1 011 rd 0010011 SLTIU
		// 	// case 0b00010011: // imm[11:0] rs1 100 rd 0010011 XORI
		// 	// case 0b00010011: // imm[11:0] rs1 110 rd 0010011 ORI
		// 	// case 0b00010011: // imm[11:0] rs1 111 rd 0010011 ANDI
		// 	// case 0b00010011: // 0000000 shamt rs1 001 rd 0010011 SLLI
		// 	// case 0b00010011: // 0000000 shamt rs1 101 rd 0010011 SRLI
		// 	// case 0b00010011: // 0100000 shamt rs1 101 rd 0010011 SRAI
		// 	// case 0b00110011: // 0000000 rs2 rs1 000 rd 0110011 ADD
		// 	// case 0b00110011: // 0100000 rs2 rs1 000 rd 0110011 SUB
		// 	// case 0b00110011: // 0000000 rs2 rs1 001 rd 0110011 SLL
		// 	// case 0b00110011: // 0000000 rs2 rs1 010 rd 0110011 SLT
		// 	// case 0b00110011: // 0000000 rs2 rs1 011 rd 0110011 SLTU
		// 	// case 0b00110011: // 0000000 rs2 rs1 100 rd 0110011 XOR
		// 	// case 0b00110011: // 0000000 rs2 rs1 101 rd 0110011 SRL
		// 	// case 0b00110011: // 0100000 rs2 rs1 101 rd 0110011 SRA
		// 	// case 0b00110011: // 0000000 rs2 rs1 110 rd 0110011 OR
		// 	// case 0b00110011: // 0000000 rs2 rs1 111 rd 0110011 AND
		// 	// case 0b00001111: // 0000 pred succ 00000 000 00000 0001111 FENCE
		// 	// case 0b00001111: // 0000 0000 0000 00000 001 00000 0001111 FENCE.I
		// 	// case 0b01110011: // 000000000000 00000 000 00000 1110011 ECALL
		// 	// case 0b01110011: // 000000000001 00000 000 00000 1110011 EBREAK
		// 	// case 0b01110011: // csr rs1 001 rd 1110011 CSRRW
		// 	// case 0b01110011: // csr rs1 010 rd 1110011 CSRRS
		// 	// case 0b01110011: // csr rs1 011 rd 1110011 CSRRC
		// 	// case 0b01110011: // csr zimm 101 rd 1110011 CSRRWI
		// 	// case 0b01110011: // csr zimm 110 rd 1110011 CSRRSI
		// 	// case 0b01110011: // csr zimm 111 rd 1110011 CSRRCI
		// 	// RV64I Base Instruction Set (in addition to RV32I)
		// 	// case 0b00000011: // imm[11:0] rs1 110 rd 0000011 LWU
		// 	// case 0b00000011: // imm[11:0] rs1 011 rd 0000011 LD
		// 	// case 0b00100011: // imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD
		// 	// case 0b00010011: // 000000 shamt rs1 001 rd 0010011 SLLI
		// 	// case 0b00010011: // 000000 shamt rs1 101 rd 0010011 SRLI
		// 	// case 0b00010011: // 010000 shamt rs1 101 rd 0010011 SRAI
		// 	// case 0b00011011: // imm[11:0] rs1 000 rd 0011011 ADDIW
		// 	// case 0b00011011: // 0000000 shamt rs1 001 rd 0011011 SLLIW
		// 	// case 0b00011011: // 0000000 shamt rs1 101 rd 0011011 SRLIW
		// 	// case 0b00011011: // 0100000 shamt rs1 101 rd 0011011 SRAIW
		// 	// case 0b00111011: // 0000000 rs2 rs1 000 rd 0111011 ADDW
		// 	// case 0b00111011: // 0100000 rs2 rs1 000 rd 0111011 SUBW
		// 	// case 0b00111011: // 0000000 rs2 rs1 001 rd 0111011 SLLW
		// 	// case 0b00111011: // 0000000 rs2 rs1 101 rd 0111011 SRLW
		// 	// case 0b00111011: // 0100000 rs2 rs1 101 rd 0111011 SRAW
		// 	// RV32M Standard Extension
		// 	// case 0b00110011: // 0000001 rs2 rs1 000 rd 0110011 MUL
		// 	// case 0b00110011: // 0000001 rs2 rs1 001 rd 0110011 MULH
		// 	// case 0b00110011: // 0000001 rs2 rs1 010 rd 0110011 MULHSU
		// 	// case 0b00110011: // 0000001 rs2 rs1 011 rd 0110011 MULHU
		// 	// case 0b00110011: // 0000001 rs2 rs1 100 rd 0110011 DIV
		// 	// case 0b00110011: // 0000001 rs2 rs1 101 rd 0110011 DIVU
		// 	// case 0b00110011: // 0000001 rs2 rs1 110 rd 0110011 REM
		// 	// case 0b00110011: // 0000001 rs2 rs1 111 rd 0110011 REMU
		// 	// RV64M Standard Extension (in addition to RV32M)
		// 	// case 0b00111011: // 0000001 rs2 rs1 000 rd 0111011 MULW
		// 	// case 0b00111011: // 0000001 rs2 rs1 100 rd 0111011 DIVW
		// 	// case 0b00111011: // 0000001 rs2 rs1 101 rd 0111011 DIVUW
		// 	// case 0b00111011: // 0000001 rs2 rs1 110 rd 0111011 REMW
		// 	// case 0b00111011: // 0000001 rs2 rs1 111 rd 0111011 REMUW
		// 	// RV32A Standard Extension
		// 	// case 0b00101111: // 00010 aq rl 00000 rs1 010 rd 0101111 LR.W
		// 	// case 0b00101111: // 00011 aq rl rs2 rs1 010 rd 0101111 SC.W
		// 	// case 0b00101111: // 00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W
		// 	// case 0b00101111: // 00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W
		// 	// case 0b00101111: // 00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W
		// 	// case 0b00101111: // 01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W
		// 	// case 0b00101111: // 01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W
		// 	// case 0b00101111: // 10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W
		// 	// case 0b00101111: // 10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W
		// 	// case 0b00101111: // 11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W
		// 	// case 0b00101111: // 11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W
		// 	// RV64A Standard Extension (in addition to RV32A)
		// 	// case 0b00101111: // 00010 aq rl 00000 rs1 011 rd 0101111 LR.D
		// 	// case 0b00101111: // 00011 aq rl rs2 rs1 011 rd 0101111 SC.D
		// 	// case 0b00101111: // 00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D
		// 	// case 0b00101111: // 00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D
		// 	// case 0b00101111: // 00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D
		// 	// case 0b00101111: // 01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D
		// 	// case 0b00101111: // 01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D
		// 	// case 0b00101111: // 10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D
		// 	// case 0b00101111: // 10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D
		// 	// case 0b00101111: // 11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D
		// 	// case 0b00101111: // 11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D
		// 	// RV32F Standard Extension
		// 	// case 0b00000111: // imm[11:0] rs1 010 rd 0000111 FLW
		// 	// case 0b00100111: // imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW
		// 	// case 0b01000011: // rs3 00 rs2 rs1 rm rd 1000011 FMADD.S
		// 	// case 0b01000111: // rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S
		// 	// case 0b01001011: // rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S
		// 	// case 0b01001111: // rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S
		// 	// case 0b01010011: // 0000000 rs2 rs1 rm rd 1010011 FADD.S
		// 	// case 0b01010011: // 0000100 rs2 rs1 rm rd 1010011 FSUB.S
		// 	// case 0b01010011: // 0001000 rs2 rs1 rm rd 1010011 FMUL.S
		// 	// case 0b01010011: // 0001100 rs2 rs1 rm rd 1010011 FDIV.S
		// 	// case 0b01010011: // 0101100 00000 rs1 rm rd 1010011 FSQRT.S
		// 	// case 0b01010011: // 0010000 rs2 rs1 000 rd 1010011 FSGNJ.S
		// 	// case 0b01010011: // 0010000 rs2 rs1 001 rd 1010011 FSGNJN.S
		// 	// case 0b01010011: // 0010000 rs2 rs1 010 rd 1010011 FSGNJX.S
		// 	// case 0b01010011: // 0010100 rs2 rs1 000 rd 1010011 FMIN.S
		// 	// case 0b01010011: // 0010100 rs2 rs1 001 rd 1010011 FMAX.S
		// 	// case 0b01010011: // 1100000 00000 rs1 rm rd 1010011 FCVT.W.S
		// 	// case 0b01010011: // 1100000 00001 rs1 rm rd 1010011 FCVT.WU.S
		// 	// case 0b01010011: // 1110000 00000 rs1 000 rd 1010011 FMV.X.W
		// 	// case 0b01010011: // 1010000 rs2 rs1 010 rd 1010011 FEQ.S
		// 	// case 0b01010011: // 1010000 rs2 rs1 001 rd 1010011 FLT.S
		// 	// case 0b01010011: // 1010000 rs2 rs1 000 rd 1010011 FLE.S
		// 	// case 0b01010011: // 1110000 00000 rs1 001 rd 1010011 FCLASS.S
		// 	// case 0b01010011: // 1101000 00000 rs1 rm rd 1010011 FCVT.S.W
		// 	// case 0b01010011: // 1101000 00001 rs1 rm rd 1010011 FCVT.S.WU
		// 	// case 0b01010011: // 1111000 00000 rs1 000 rd 1010011 FMV.W.X
		// 	// RV64F Standard Extension (in addition to RV32F)
		// 	// case 0b01010011: // 1100000 00010 rs1 rm rd 1010011 FCVT.L.S
		// 	// case 0b01010011: // 1100000 00011 rs1 rm rd 1010011 FCVT.LU.S
		// 	// case 0b01010011: // 1101000 00010 rs1 rm rd 1010011 FCVT.S.L
		// 	// case 0b01010011: // 1101000 00011 rs1 rm rd 1010011 FCVT.S.LU
		// 	// RV32D Standard Extension
		// 	// case 0b00000111: // imm[11:0] rs1 011 rd 0000111 FLD
		// 	// case 0b00100111: // imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD
		// 	// case 0b01000011: // rs3 01 rs2 rs1 rm rd 1000011 FMADD.D
		// 	// case 0b01000111: // rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D
		// 	// case 0b01001011: // rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D
		// 	// case 0b01001111: // rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D
		// 	// case 0b01010011: // 0000001 rs2 rs1 rm rd 1010011 FADD.D
		// 	// case 0b01010011: // 0000101 rs2 rs1 rm rd 1010011 FSUB.D
		// 	// case 0b01010011: // 0001001 rs2 rs1 rm rd 1010011 FMUL.D
		// 	// case 0b01010011: // 0001101 rs2 rs1 rm rd 1010011 FDIV.D
		// 	// case 0b01010011: // 0101101 00000 rs1 rm rd 1010011 FSQRT.D
		// 	// case 0b01010011: // 0010001 rs2 rs1 000 rd 1010011 FSGNJ.D
		// 	// case 0b01010011: // 0010001 rs2 rs1 001 rd 1010011 FSGNJN.D
		// 	// case 0b01010011: // 0010001 rs2 rs1 010 rd 1010011 FSGNJX.D
		// 	// case 0b01010011: // 0010101 rs2 rs1 000 rd 1010011 FMIN.D
		// 	// case 0b01010011: // 0010101 rs2 rs1 001 rd 1010011 FMAX.D
		// 	// case 0b01010011: // 0100000 00001 rs1 rm rd 1010011 FCVT.S.D
		// 	// case 0b01010011: // 0100001 00000 rs1 rm rd 1010011 FCVT.D.S
		// 	// case 0b01010011: // 1010001 rs2 rs1 010 rd 1010011 FEQ.D
		// 	// case 0b01010011: // 1010001 rs2 rs1 001 rd 1010011 FLT.D
		// 	// case 0b01010011: // 1010001 rs2 rs1 000 rd 1010011 FLE.D
		// 	// case 0b01010011: // 1110001 00000 rs1 001 rd 1010011 FCLASS.D
		// 	// case 0b01010011: // 1100001 00000 rs1 rm rd 1010011 FCVT.W.D
		// 	// case 0b01010011: // 1100001 00001 rs1 rm rd 1010011 FCVT.WU.D
		// 	// case 0b001010011: // 1101001 00000 rs1 rm rd 1010011 FCVT.D.W
		// 	// case 0b01010011: // 1101001 00001 rs1 rm rd 1010011 FCVT.D.WU
		// 	// RV64D Standard Extension (in addition to RV32D)
		// 	// case 0b01010011: // 1100001 00010 rs1 rm rd 1010011 FCVT.L.D
		// 	// case 0b01010011: // 1100001 00011 rs1 rm rd 1010011 FCVT.LU.D
		// 	// case 0b01010011: // 1110001 00000 rs1 000 rd 1010011 FMV.X.D
		// 	// case 0b01010011: // 1101001 00010 rs1 rm rd 1010011 FCVT.D.L
		// 	// case 0b01010011: // 1101001 00011 rs1 rm rd 1010011 FCVT.D.LU
		// 	// case 0b01010011: // 1111001 00000 rs1 000 rd 1010011 FMV.D.X
		// 	default:
		// 		log.Panicln("")
		// 	}
		// default:
		// 	log.Println("")
		// }

		// 0b10001110_00001001
		// INSTR: [ instruction 0x8e09 rs1=0xc rs2=0xa rd=0xc imm=0(0x0) func=sub ]

		log.Panicln("")
	}
}

var (
	cArgs = []string{"main"}
	cEnvs = []string{}
)

func main() {
	flag.Parse()

	cpu := &CPU{
		ModuleBase: &riscv.RegisterRV64I{
			RG: [32]uint64{},
			PC: 0,
		},
		Mem: make([]byte, 4*1024*1024),
	}

	f, err := elf.Open(flag.Arg(0))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	cpu.ModuleBase.PC = f.Entry

	for _, s := range f.Sections {
		if s.Flags&elf.SHF_ALLOC == 0 {
			continue
		}
		if _, err := s.ReadAt(cpu.Mem[s.Addr:s.Addr+s.Size], 0); err != nil {
			log.Panicln(err)
		}
	}
	cpu.ModuleBase.RG[riscv.Rsp] = uint64(len(cpu.Mem))

	// Command line parameters, distribution of environment variables on the stack:
	//
	// | envs[1]     | SP Base
	// | envs[0]     |
	// | argv[1]     |
	// | argv[0]     |
	// | \0          |
	// | envs[1].ptr |
	// | envs[0].ptr |
	// | \0          |
	// | argv[1].ptr |
	// | argv[0].ptr |
	// | argc        |

	addr := []uint64{0}
	// for i := len(cEnvs) - 1; i >= 0; i-- {
	// 	cpu.pushString(cEnvs[i])
	// 	addr = append(addr, cpu.ModuleBase.RG[riscv.Rsp])
	// }
	// addr = append(addr, 0)
	for i := len(cArgs) - 1; i >= 0; i-- {
		cpu.pushString(cArgs[i])
		addr = append(addr, cpu.ModuleBase.RG[riscv.Rsp])
	}
	// Align the stack to 8 bytes
	cpu.ModuleBase.RG[riscv.Rsp] &^= 0x7
	for _, a := range addr {
		cpu.pushUint64(a)
	}
	cpu.pushUint64(uint64(len(cArgs)))
	cpu.Run()
}
