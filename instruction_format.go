package riscv

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 19

func RType(i uint64) (rd int, rs1 int, rs2 int) {
	rd = int(InstructionPart(i, 7, 11))
	rs1 = int(InstructionPart(i, 15, 19))
	rs2 = int(InstructionPart(i, 20, 24))
	return
}

func IType(i uint64) (rd int, rs1 int, imm uint64) {
	rd = int(InstructionPart(i, 7, 11))
	rs1 = int(InstructionPart(i, 15, 19))
	imm = InstructionPart(i, 20, 31)
	return
}

func SType(i uint64) (rs1 int, rs2 int, imm uint64) {
	rs1 = int(InstructionPart(i, 15, 19))
	rs2 = int(InstructionPart(i, 20, 24))
	imm = InstructionPart(i, 25, 31)<<5 | InstructionPart(i, 7, 11)
	return
}

func BType(i uint64) (rs1 int, rs2 int, imm uint64) {
	rs1 = int(InstructionPart(i, 15, 19))
	rs2 = int(InstructionPart(i, 20, 24))
	imm = InstructionPart(i, 31, 31)<<12 | InstructionPart(i, 7, 7)<<11 | InstructionPart(i, 25, 30)<<5 | InstructionPart(i, 8, 11)<<1
	return
}

func UType(i uint64) (rd int, imm uint64) {
	rd = int(InstructionPart(i, 7, 11))
	imm = InstructionPart(i, 12, 31) << 12
	return
}

func JType(i uint64) (rd int, imm uint64) {
	rd = int(InstructionPart(i, 7, 11))
	imm = InstructionPart(i, 31, 31)<<20 | InstructionPart(i, 12, 19)<<12 | InstructionPart(i, 20, 20)<<11 | InstructionPart(i, 21, 30)<<1
	return
}
