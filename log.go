package riscv

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

func DebuglnRType(i string, rd int, rs1 int, rs2 int) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x rs1: 0x%02x rs2: 0x%02x", i, rd, rs1, rs2))
}

func DebuglnIType(i string, rd int, rs1 int, imm int32) {
	Debugln(fmt.Sprintf("Instr: % 10s | rd: 0x%02x rs1: 0x%02x imm: 0x%04x", i, rd, rs1, imm))
}
func DebuglnSType() {}
func DebuglnBType() {}
func DebuglnUType(i string, rd int, imm uint32) {
	log.Printf("Instr: % 10s | rd: 0x%02x imm: 0x%04x", i, rd, imm)
}
func DebuglnJType() {}
