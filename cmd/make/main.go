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
	cEmu       = "./bin/rv64"
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

func makeBinary() {
	os.Mkdir("bin", 0755)
	call("go", "build", "-o", "bin", "github.com/mohanson/rv64/cmd/make")
	call("go", "build", "-o", "bin", "github.com/mohanson/rv64/cmd/rv64")
}

func makeRiscvTests() {
	os.Chdir(*cTmp)
	defer os.Chdir(cPwd)
	if _, err := os.Stat("riscv-tests"); err == nil {
		return
	}
	call("git", "clone", "https://github.com/libraries/riscv-tests")
	os.Chdir("riscv-tests")
	defer os.Chdir("..")
	call("git", "submodule", "update", "--init", "--recursive")
	call("autoconf")
	call("./configure", "--prefix="+cRiscvTool)
	call("make", "isa")
}

func testRiscvTests() {
	m, err := filepath.Glob(filepath.Join(*cTmp, "riscv-tests", "isa", "rv64u[ima]-u-*"))
	if err != nil {
		log.Panicln(err)
	}
	for _, e := range m {
		if strings.HasSuffix(e, ".dump") {
			continue
		}
		call(cEmu, e)
	}
}

func main() {
	if cRiscvTool == "" {
		log.Panicln("$RISCV undefined")
	}
	flag.Parse()
	for _, e := range flag.Args() {
		if e == "make" {
			makeBinary()
		}
		if e == "test" {
			makeRiscvTests()
			testRiscvTests()
		}
	}
}
