package riscv

import "log"

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 1.2

func InstructionLengthEncoding(b []byte) int {
	if len(b) != 2 {
		log.Panicln("")
	}
	// xxxxxxxxxxxxxxaa 16-bit, aa != 11
	if b[0]&0x03 != 0x03 {
		return 2
	}
	// xxxxxxxxxxxbbb11 32-bit, bbb != 111
	if b[0]&0x1c != 0x1c {
		return 4
	}
	// xxxxxxxxxx011111 48-bit
	if b[0]&0x20 == 0x00 {
		return 6
	}
	// xxxxxxxxx0111111 64-bit
	if b[0]&0x40 == 0x00 {
		return 8
	}
	// xnnnxxxxx1111111 (80+16*nnn)-bit, nnn != 111
	if b[1]&0x70 != 0x70 {
		n := (b[1] & 0x70) >> 4
		return int(10 + 2*n)
	}
	// x111xxxxx1111111 Reserved for ≥192-bits
	log.Panicln("Reserved for ≥192-bits")
	return 0
}
