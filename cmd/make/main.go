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

func makeExamples() {
	os.Mkdir("bin", 0755)
	os.Mkdir("bin/res", 0755)
	os.Mkdir("bin/res/program", 0755)
	call(*cGCC, "-o", "bin/res/program/andi", "res/program/andi.c")
	call(*cGCC, "-o", "bin/res/program/fib_args", "res/program/fib_args.c")
	call(*cGCC, "-o", "bin/res/program/fib", "res/program/fib.c")
	call(*cGCC, "-o", "bin/res/program/math", "res/program/math.c")
	call(*cGCC, "-o", "bin/res/program/minimal", "res/program/minimal.c")
}

func testExamples() {
	call(cEmu, "--", "bin/res/program/andi")
	call(cEmu, "--", "bin/res/program/fib_args", "10", "55")
	call(cEmu, "--", "bin/res/program/fib_args", "9", "34")
	call(cEmu, "--", "bin/res/program/fib_args", "8", "21")
	call(cEmu, "--", "bin/res/program/fib_args", "7", "13")
	call(cEmu, "--", "bin/res/program/fib_args", "6", "8")
	call(cEmu, "--", "bin/res/program/fib")
	call(cEmu, "--", "bin/res/program/math")
	call(cEmu, "--", "bin/res/program/minimal")
}

func makeRiscvTests() {
	os.Chdir("res")
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
	m, err := filepath.Glob(filepath.Join("res", "riscv-tests", "isa", "rv64u[imafd]-u-*"))
	if err != nil {
		log.Panicln(err)
	}
	for _, e := range m {
		if strings.HasSuffix(e, ".dump") {
			continue
		}
		call(cEmu, "--", e)
	}
}

func main() {
	if cRiscvTool == "" {
		log.Panicln("$RISCV undefined")
	}
	flag.Parse()
	if flag.NArg() == 0 {
		makeBinary()
		return
	}
	for _, e := range flag.Args() {
		switch e {
		case "make":
			makeBinary()
		case "test":
			makeRiscvTests()
			makeExamples()
			testRiscvTests()
			testExamples()
		}
	}
}
