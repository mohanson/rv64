package riscv

func ExecuterRV64I(r *RegisterRV64I, m []byte, i uint64) int {
	switch {
	// case 0b00000011: // imm[11:0] rs1 110 rd 0000011 LWU
	// case 0b00000011: // imm[11:0] rs1 011 rd 0000011 LD
	// case 0b00100011: // imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD
	// case 0b00010011: // 000000 shamt rs1 001 rd 0010011 SLLI
	// case 0b00010011: // 000000 shamt rs1 101 rd 0010011 SRLI
	// case 0b00010011: // 010000 shamt rs1 101 rd 0010011 SRAI
	// case 0b00011011: // imm[11:0] rs1 000 rd 0011011 ADDIW
	// case 0b00011011: // 0000000 shamt rs1 001 rd 0011011 SLLIW
	// case 0b00011011: // 0000000 shamt rs1 101 rd 0011011 SRLIW
	// case 0b00011011: // 0100000 shamt rs1 101 rd 0011011 SRAIW
	// case 0b00111011: // 0000000 rs2 rs1 000 rd 0111011 ADDW
	// case 0b00111011: // 0100000 rs2 rs1 000 rd 0111011 SUBW
	// case 0b00111011: // 0000000 rs2 rs1 001 rd 0111011 SLLW
	// case 0b00111011: // 0000000 rs2 rs1 101 rd 0111011 SRLW
	// case 0b00111011: // 0100000 rs2 rs1 101 rd 0111011 SRAW
	}
	return 0
}
