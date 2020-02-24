package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	cCompile = flag.String("gcc", "/root/app/riscv/bin/riscv64-unknown-elf-gcc", "")
)

func call(name string, arg ...string) {
	log.Println("$", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Panicln(err)
	}
}

func main() {
	flag.Parse()
	// call(*cCompile, "-o", "./bin/fuzz_32i", "-march=rv32i", "-mabi=ilp32", "./res/fuzz.c")
	call(*cCompile, "-o", "./bin/fuzz", "./res/fuzz.c")
}
