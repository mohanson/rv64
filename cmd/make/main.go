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
	os.Mkdir("build", 0755)
	os.Mkdir("build/res", 0755)
	os.Mkdir("build/res/program", 0755)
	call(*cCompile, "-o", "./build/res/program/minimal", "./res/program/minimal.c")
}
