package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mohanson/rv64"
)

type CPU struct {
	Inner *rv64.CPU
}

func (c *CPU) pushString(s string) {
	bs := append([]byte(s), 0x00)
	c.Inner.SetRegister(rv64.Rsp, c.Inner.GetRegister(rv64.Rsp)-uint64(len(bs)))
	for i, b := range bs {
		c.Inner.GetMemory().SetByte(c.Inner.GetRegister(rv64.Rsp)+uint64(i), []byte{b})
	}
}

func (c *CPU) pushUint64(v uint64) {
	c.Inner.SetRegister(rv64.Rsp, c.Inner.GetRegister(rv64.Rsp)-8)
	mem := make([]byte, 8)
	binary.LittleEndian.PutUint64(mem, v)
	c.Inner.GetMemory().SetByte(c.Inner.GetRegister(rv64.Rsp), mem)
}

var (
	cStep  = flag.Int64("steps", -1, "")
	cDebug = flag.Bool("d", false, "Debug")
)

func (c *CPU) Run() uint8 {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	i := 0
	for {
		if c.Inner.GetStatus() == 1 {
			rv64.Debugln("Exit:", c.Inner.GetSystem().Code())
			return c.Inner.GetSystem().Code()
		}
		if i > int(*cStep) && *cStep > 0 {
			break
		}
		data, err := c.Inner.PipelineInstructionFetch()
		if err != nil {
			rv64.Panicln(err)
		}
		rv64.Debugln("==========")
		if len(data) == 2 {
			rv64.Debugln(fmt.Sprintf("%08b %08b", data[1], data[0]))
		} else if len(data) == 4 {
			rv64.Debugln(fmt.Sprintf("%08b %08b %08b %08b", data[3], data[2], data[1], data[0]))
		} else {
			rv64.Panicln("")
		}
		var s uint64 = 0
		for i := 0; i < 32; i++ {
			s += c.Inner.GetRegister(uint64(i))
		}
		rv64.Debugln(i, c.Inner.GetPC(), s)

		n, err := c.Inner.PipelineExecute(data)
		if err != nil {
			log.Panicln(err)
		}
		if n != 0 {
			i += 1
			continue
		}

		if len(data) == 4 {
			var s uint64 = 0
			for i := len(data) - 1; i >= 0; i-- {
				s += uint64(data[i]) << (8 * i)
			}

			// n, err = rv64.ExecuterF(c.Inner, s)
			// if err != nil {
			// 	log.Panicln(err)
			// }
			// if n != 0 {
			// 	i += 1
			// 	// c.Inner.SetCSR(rv64.Rdcycle, c.Inner.GetCSR(rv64.Rdcycle)+n)
			// 	// c.Inner.SetCSR(rv64.Rdtime, c.Inner.GetCSR(rv64.Rdtime)+n)
			// 	// c.Inner.SetCSR(rv64.Rdinstret, c.Inner.GetCSR(rv64.Rdtime)+n)
			// 	continue
			// }
		}
		log.Panicln("")
	}
	return 0
}

var (
	cArgs = []string{"main"}
	cEnvs = []string{}
)

func main() {
	flag.Parse()
	if *cDebug == true {
		rv64.LogLevel = 1
	}
	inner := rv64.NewCPU()
	inner.SetMemory(rv64.NewMemoryLinear(4 * 1024 * 1024))
	inner.SetSystem(rv64.NewSystemStandard())
	cpu := &CPU{
		Inner: inner,
	}

	f, err := elf.Open(flag.Arg(0))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	cpu.Inner.SetPC(f.Entry)

	for _, s := range f.Sections {
		if s.Flags&elf.SHF_ALLOC != 0 && s.Type&elf.SHT_NOBITS == 0 {
			mem := make([]byte, s.Size)
			if _, err := s.ReadAt(mem, 0); err != nil {
				log.Panicln(err)
			}
			cpu.Inner.GetMemory().SetByte(s.Addr, mem)
		}
	}
	cpu.Inner.SetRegister(rv64.Rsp, cpu.Inner.GetMemory().Len())

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
		cpu.pushString(cArgs[i])
		addr = append(addr, cpu.Inner.GetRegister(rv64.Rsp))
	}
	cpu.Inner.GetMemory().SetUint8(cpu.Inner.GetRegister(rv64.Rsp), 0)
	cpu.Inner.SetRegister(rv64.Rsp, cpu.Inner.GetRegister(rv64.Rsp)-1)
	for i := len(addr) - 1; i >= 0; i-- {
		cpu.pushUint64(addr[i])
	}
	cpu.pushUint64(uint64(len(cArgs)))
	if cpu.Inner.GetRegister(rv64.Rsp) != 4194282 {
		log.Panicln("")
	}
	// Align the stack to 16 bytes
	cpu.Inner.SetRegister(rv64.Rsp, cpu.Inner.GetRegister(rv64.Rsp)&0xfffffff0)
	if cpu.Inner.GetRegister(rv64.Rsp) != 4194272 {
		log.Panicln("")
	}
	os.Exit(int(cpu.Run()))
}
