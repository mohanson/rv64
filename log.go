package rv64

import (
	"fmt"
	"log"
)

func Println(v ...interface{}) {
	log.Println(v...)
}

func Debugln(v ...interface{}) {
	if LogLevel > 0 {
		log.Println(v...)
	}
}

func Panicln(v ...interface{}) {
	log.Panicln(v...)
}

func I(c *CPU, n uint64) string {
	return fmt.Sprintf("%#02x(%#016x)", n, c.GetRegister(n))
}

func F(c *CPU, n uint64) string {
	return fmt.Sprintf("%#02x(%#016x)", n, c.GetRegisterFloat(n))
}

func DebuglnRType(i string, rd uint64, rs1 uint64, rs2 uint64) {
	Debugln(fmt.Sprintf("% 10s rd: %#02x rs1: %#02x rs2: %#02x", i, rd, rs1, rs2))
}

func DebuglnR4Type(i string, rd uint64, rs1 uint64, rs2 uint64, rs3 uint64) {
	Debugln(fmt.Sprintf("% 10s rd: %#02x rs1: %#02x rs2: %#02x rs3: %#02x", i, rd, rs1, rs2, rs3))
}

func DebuglnIType(i string, rd uint64, rs1 uint64, imm uint64) {
	Debugln(fmt.Sprintf("% 10s rd: %#02x rs1: %#02x imm: %#04x", i, rd, rs1, imm))
}

func DebuglnSType(i string, rs1 uint64, rs2 uint64, imm uint64) {
	Debugln(fmt.Sprintf("% 10s rs1: %#02x rs2: %#02x imm: %#04x", i, rs1, rs2, imm))
}

func DebuglnBType(i string, rs1 uint64, rs2 uint64, imm uint64) {
	Debugln(fmt.Sprintf("% 10s rs1: %#02x rs2: %#02x imm: %#04x", i, rs1, rs2, imm))
}

func DebuglnUType(i string, rd uint64, imm uint64) {
	Debugln(fmt.Sprintf("% 10s rd: %#02x imm: %#04x", i, rd, imm))
}

func DebuglnJType(i string, rd uint64, imm uint64) {
	Debugln(fmt.Sprintf("% 10s rd: %#02x imm: %#04x", i, rd, imm))
}
