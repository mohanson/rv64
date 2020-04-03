package main

import (
	"debug/elf"
	"flag"
	"log"
	"os"

	"github.com/mohanson/rv64"
)

var (
	cStep  = flag.Int64("steps", -1, "")
	cDebug = flag.Bool("d", false, "Debug")
)

var (
	cArgs = []string{"main"}
	cEnvs = []string{}
)

func main() {
	flag.Parse()
	if *cDebug == true {
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

	for _, s := range f.Sections {
		if s.Flags&elf.SHF_ALLOC != 0 && s.Type != elf.SHT_NOBITS {
			mem := make([]byte, s.Size)
			if _, err := s.ReadAt(mem, 0); err != nil {
				log.Panicln(err)
			}
			cpu.GetMemory().SetByte(s.Addr, mem)
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

	addr := []uint64{}
	// for i := len(cEnvs) - 1; i >= 0; i-- {
	// 	cpu.pushString(cEnvs[i])
	// 	addr = append(addr, cpu.ModuleBase.RG[riscv.Rsp])
	// }
	// addr = append(addr, 0)
	for i := len(cArgs) - 1; i >= 0; i-- {
		cpu.PushString(cArgs[i])
		addr = append(addr, cpu.GetRegister(rv64.Rsp))
	}
	cpu.GetMemory().SetUint8(cpu.GetRegister(rv64.Rsp), 0)
	cpu.SetRegister(rv64.Rsp, cpu.GetRegister(rv64.Rsp)-1)
	for i := len(addr) - 1; i >= 0; i-- {
		cpu.PushUint64(addr[i])
	}
	cpu.PushUint64(uint64(len(cArgs)))
	if cpu.GetRegister(rv64.Rsp) != 4194282 {
		log.Panicln("")
	}
	// Align the stack to 16 bytes
	cpu.SetRegister(rv64.Rsp, cpu.GetRegister(rv64.Rsp)&0xfffffff0)
	if cpu.GetRegister(rv64.Rsp) != 4194272 {
		log.Panicln("")
	}
	os.Exit(int(cpu.Execute()))
}
