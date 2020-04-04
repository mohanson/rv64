package main

import (
	"debug/elf"
	"flag"
	"log"
	"os"

	"github.com/mohanson/rv64"
)

var (
	flDebug = flag.Bool("d", false, "Debug")
)

func main() {
	flag.Parse()
	if *flDebug {
		rv64.LogLevel = 1
	}
	cpu := rv64.NewCPU()
	cpu.SetFasten(rv64.NewLinear(4 * 1024 * 1024))
	cpu.SetSystem(rv64.NewSystemStandard())
	cpu.SetCSR(rv64.NewCSRStandard())

	f, err := elf.Open(flag.Arg(0))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	for _, p := range f.Progs {
		// Specifies a loadable segment, described by p_filesz and p_memsz. The bytes from the file are mapped to the
		// beginning of the memory segment. If the segment's memory size (p_memsz) is larger than the file size
		// (p_filesz), the extra bytes are defined to hold the value 0 and to follow the segment's initialized area.
		// The file size can not be larger than the memory size. Loadable segment entries in the program header table
		// appear in ascending order, sorted on the p_vaddr member.
		if p.ProgHeader.Type == elf.PT_LOAD {
			mem := make([]byte, p.Memsz)
			p.ReadAt(mem[0:p.Filesz], 0)
			cpu.GetMemory().SetByte(p.Vaddr, mem)
		}
	}
	cpu.SetPC(f.Entry)
	cpu.SetRegister(rv64.Rsp, cpu.GetMemory().Len())

	// Command line parameters, distribution of environment variables on the stack:
	//
	// | envs[1]     | SP Base
	// | envs[0]     |
	// | argv[1]     |
	// | argv[0]     |
	// | \0          |
	// | envs[1].ptr |
	// | envs[0].ptr |
	// | \0          |
	// | argv[1].ptr |
	// | argv[0].ptr |
	// | argc        |
	argList := flag.Args()
	envList := []string{}
	envPtrs := []uint64{}
	argPtrs := []uint64{}
	for i := len(envList) - 1; i >= 0; i-- {
		cpu.PushString(envList[i])
		envPtrs = append(envPtrs, cpu.GetRegister(rv64.Rsp))
	}
	for i := len(argList) - 1; i >= 0; i-- {
		cpu.PushString(argList[i])
		argPtrs = append(argPtrs, cpu.GetRegister(rv64.Rsp))
	}
	cpu.PushUint8(0)
	for i := 0; i < len(envPtrs); i++ {
		cpu.PushUint64(envPtrs[i])
	}
	cpu.PushUint8(0)
	for i := 0; i < len(argPtrs); i++ {
		cpu.PushUint64(argPtrs[i])
	}
	cpu.PushUint64(uint64(len(argList)))

	os.Exit(int(cpu.Run()))
}
