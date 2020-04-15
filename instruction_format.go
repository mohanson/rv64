package rv64

// |31|30|29|28|27|26|25|24|23|22|21|20|19|18|17|16|15|14|13|12|11|10|09|08|07|06|05|04|03|02|01|00|
// |funct7              |rs2           |rs1           |funct3  |rd            |opcode              | R-type
// |rs3           |funct|rs2           |rs1           |funct3  |rd            |opcode              | R4-type
// |imm[11:0]                          |rs1           |funct3  |rd            |opcode              | I-type
// |imm[11:5]           |rs2           |rs1           |funct3  |imm[4:0]      |opcode              | S-type
// |imm[31:12]                                                 |rd            |opcode              | U-type
// |12|imm[10:5]        |rs2           |rs1           |funct3  |imm[4:1]   |11|opcode              | B-type
// |20|imm[10:1]                    |11|imm[19:12]             |rd            |opcode              | J-type
//
// Except for the 5-bit immediates used in CSR instructions (Chapter 9), immediates are always sign-extended, and are
// generally packed towards the leftmost available bits in the instruction and have been allocated to reduce hardware
// complexity. In particular, the sign bit for all immediates is always in bit 31 of the instruction to speed
// sign-extension circuitry.

func RType(i uint64) (rd uint64, rs1 uint64, rs2 uint64) {
	rd = InstructionPart(i, 7, 11)
	rs1 = InstructionPart(i, 15, 19)
	rs2 = InstructionPart(i, 20, 24)
	return
}

func IType(i uint64) (rd uint64, rs1 uint64, imm uint64) {
	rd = InstructionPart(i, 7, 11)
	rs1 = InstructionPart(i, 15, 19)
	imm = InstructionPart(i, 20, 31)
	return
}

func SType(i uint64) (rs1 uint64, rs2 uint64, imm uint64) {
	rs1 = InstructionPart(i, 15, 19)
	rs2 = InstructionPart(i, 20, 24)
	imm = SignExtend(InstructionPart(i, 25, 31)<<5|InstructionPart(i, 7, 11), 11)
	return
}

func UType(i uint64) (rd uint64, imm uint64) {
	rd = InstructionPart(i, 7, 11)
	imm = SignExtend(InstructionPart(i, 12, 31)<<12, 31)
	return
}

func BType(i uint64) (rs1 uint64, rs2 uint64, imm uint64) {
	rs1 = InstructionPart(i, 15, 19)
	rs2 = InstructionPart(i, 20, 24)
	imm = SignExtend(InstructionPart(i, 31, 31)<<12|InstructionPart(i, 7, 7)<<11|InstructionPart(i, 25, 30)<<5|InstructionPart(i, 8, 11)<<1, 12)
	return
}

func JType(i uint64) (rd uint64, imm uint64) {
	rd = InstructionPart(i, 7, 11)
	imm = SignExtend(InstructionPart(i, 31, 31)<<20|InstructionPart(i, 12, 19)<<12|InstructionPart(i, 20, 20)<<11|InstructionPart(i, 21, 30)<<1, 19)
	return
}

func R4Type(i uint64) (rd uint64, rs1 uint64, rs2 uint64, rs3 uint64) {
	rd = InstructionPart(i, 7, 11)
	rs1 = InstructionPart(i, 15, 19)
	rs2 = InstructionPart(i, 20, 24)
	rs3 = InstructionPart(i, 27, 31)
	return
}
