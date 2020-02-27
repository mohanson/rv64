package riscv

import "log"

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 19

func RType(data []byte) (opcode int, rd int, funct3 int, rs1 int, rs2 int, funct7 int) {
	opcode = int(InstructionPart(data, 0, 6))
	rd = int(InstructionPart(data, 7, 11))
	funct3 = (int(data[1]) >> 4) & 0x07
	rs1 = ((int(data[2]) & 0x0f) << 1) | (int(data[1]) >> 7)
	rs2 = ((int(data[4]) & 0x01) << 3) | (int(data[2]) >> 5)
	funct7 = int(data[4]) >> 1
	return
}

func IType(data []byte) (opcode int, rd int, funct3 int, rs1 int, imm uint64) {
	opcode = int(InstructionPart(data, 0, 6))
	rd = (int(data[1])&0x0f)<<1 | (int(data[0]) >> 7)
	log.Println(data)
	log.Println("!!!", rd, int(InstructionPart(data, 7, 11)))
	funct3 = (int(data[1]) >> 4) & 0x07
	rs1 = ((int(data[2]) & 0x0f) << 1) | (int(data[1]) >> 7)
	imm = uint64((int(data[3]) << 4) | int(data[2])>>4)
	return
}

func SType(data []byte) (opcode int, imm uint64, funct3 int, rs1 int, rs2 int) {
	opcode = int(InstructionPart(data, 0, 6))
	imm = uint64(((int(data[3]) & 0xf7) << 4) | ((int(data[1]) & 0x0f) << 1) | ((int(data[0]) & 0x80) >> 7))
	funct3 = (int(data[1]) >> 4) & 0x07
	rs1 = ((int(data[2]) & 0x0f) << 1) | (int(data[1]) >> 7)
	rs2 = ((int(data[4]) & 0x01) << 3) | (int(data[2]) >> 5)
	return
}

func BType() {}

func UType(data []byte) (opcode int, rd int, imm uint32) {
	opcode = int(InstructionPart(data, 0, 6))
	rd = (int(data[1])&0x0f)<<1 | (int(data[0]) >> 7)
	imm = uint32(((int(data[3]) << 12) | (int(data[2]) << 4) | (int(data[1]) >> 4)) << 12)
	return
}

func JType() {}
