package riscv

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 19

func RType(i uint64) (opcode int, rd int, funct3 int, rs1 int, rs2 int, funct7 int) {
	opcode = int(InstructionPart(i, 0, 6))
	rd = int(InstructionPart(i, 7, 11))
	funct3 = int(InstructionPart(i, 12, 14))
	rs1 = int(InstructionPart(i, 15, 19))
	rs2 = int(InstructionPart(i, 20, 24))
	funct7 = int(InstructionPart(i, 25, 31))
	return
}

func IType(i uint64) (opcode int, rd int, funct3 int, rs1 int, imm uint64) {
	opcode = int(InstructionPart(i, 0, 6))
	rd = int(InstructionPart(i, 7, 11))
	funct3 = int(InstructionPart(i, 12, 14))
	rs1 = int(InstructionPart(i, 15, 19))
	imm = InstructionPart(i, 20, 31)
	return
}

func SType(i uint64) (opcode int, imm uint64, funct3 int, rs1 int, rs2 int) {
	opcode = int(InstructionPart(i, 0, 6))
	imm = InstructionPart(i, 25, 31)<<5 | InstructionPart(i, 7, 11)
	funct3 = int(InstructionPart(i, 12, 14))
	rs1 = int(InstructionPart(i, 15, 19))
	rs2 = int(InstructionPart(i, 20, 24))
	return
}

func BType() (opcode int, funct3 int, rs1 int, rs2 int, imm uint64) {
	return
}

func UType(i uint64) (opcode int, rd int, imm uint64) {
	opcode = int(InstructionPart(i, 0, 6))
	rd = int(InstructionPart(i, 7, 11))
	imm = InstructionPart(i, 12, 31) << 12
	return
}

func JType(i uint64) (opcode int, rd int, imm uint64) {
	opcode = int(InstructionPart(i, 0, 6))
	rd = int(InstructionPart(i, 7, 11))
	imm = InstructionPart(i, 31, 31)<<20 | InstructionPart(i, 12, 19)<<12 | InstructionPart(i, 20, 20)<<11 | InstructionPart(i, 21, 30)
	return
}
