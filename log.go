package rv64

import (
	"fmt"
	"log"
)

func Println(v ...interface{}) {
	log.Println(v...)
}

func Debugln(v ...interface{}) {
	log.Println(v...)
}

func Panicln(v ...interface{}) {
	log.Panicln(v...)
}

func DebuglnRType(i string, rd uint64, rs1 uint64, rs2 uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x rs1: 0x%02x rs2: 0x%02x", i, rd, rs1, rs2))
}

func DebuglnIType(i string, rd uint64, rs1 uint64, imm uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x rs1: 0x%02x imm: 0x%04x", i, rd, rs1, imm))
}

func DebuglnSType(i string, rs1 uint64, rs2 uint64, imm uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rs1: 0x%02x rs2: 0x%02x imm: 0x%04x", i, rs1, rs2, imm))
}

func DebuglnBType(i string, rs1 uint64, rs2 uint64, imm uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rs1: 0x%02x rs2: 0x%02x imm: 0x%04x", i, rs1, rs2, imm))
}

func DebuglnUType(i string, rd uint64, imm uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x imm: 0x%04x", i, rd, imm))
}

func DebuglnJType(i string, rd uint64, imm uint64) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x imm: 0x%04x", i, rd, imm))
}
