package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	cTmp = flag.String("tmp", "/tmp", "")
	cGCC = flag.String("gcc", filepath.Join(cRiscvTool, "bin", "riscv64-unknown-elf-gcc"), "")
)

var (
	cPwd, _    = os.Getwd()
	cRiscvTool = os.Getenv("RISCV")
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

func makeRiscvTests() {
	os.Chdir(*cTmp)
	if _, err := os.Stat("riscv-tests"); err == nil {
		return
	}
	call("git", "clone", "https://github.com/nervosnetwork/riscv-tests")
	os.Chdir("riscv-tests")
	call("git", "submodule", "update", "--init", "--recursive")
	call("autoconf")
	call("./configure", "--prefix="+cRiscvTool)
	call("make", "isa")
}

func main() {
	if cRiscvTool == "" {
		log.Panicln("$RISCV undefined")
	}
	flag.Parse()
	if flag.Arg(0) == "test" {
		makeRiscvTests()
	}
}

// os.Mkdir("build", 0755)
// os.Mkdir("build/res", 0755)
// os.Mkdir("build/res/program", 0755)
// call(*cCompile, "-o", "./build/res/program/minimal", "./res/program/minimal.c")
