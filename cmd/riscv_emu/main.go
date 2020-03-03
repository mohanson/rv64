package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"log"

	"github.com/mohanson/riscv"
)

const cDebug = 1

type CPU struct {
	Inner *riscv.CPU
}

func (c *CPU) pushString(s string) {
	bs := append([]byte(s), 0x00)
	c.Inner.SetRegister(riscv.Rsp, c.Inner.GetRegister(riscv.Rsp)-uint64(len(bs)))
	for i, b := range bs {
		c.Inner.GetMemory().Set(c.Inner.GetRegister(riscv.Rsp)+uint64(i), []byte{b})
	}
}

func (c *CPU) pushUint64(v uint64) {
	c.Inner.SetRegister(riscv.Rsp, c.Inner.GetRegister(riscv.Rsp))
	mem := make([]byte, 8)
	binary.LittleEndian.PutUint64(mem, v)
	c.Inner.GetMemory().Set(c.Inner.GetRegister(riscv.Rsp), mem)
}

func (c *CPU) FetchInstruction() []byte {
	a, err := c.Inner.GetMemory().Get(c.Inner.GetPC(), 2)
	if err != nil {
		log.Panicln(err)
	}
	b := riscv.InstructionLengthEncoding(a)
	instructionBytes, err := c.Inner.GetMemory().Get(c.Inner.GetPC(), uint64(b))
	if err != nil {
		log.Panicln(err)
	}
	return instructionBytes
}

var cStep = flag.Int64("steps", 250, "")

func (c *CPU) Run() {
	flag.Parse()
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	i := 0
	for {
		if c.Inner.GetStatus() == 1 {
			log.Println("Exit:", c.Inner.GetSystem().Code())
			break
		}
		if i > int(*cStep) {
			break
		}
		data := c.FetchInstruction()
		log.Println("==========")
		if len(data) == 2 {
			log.Printf("%08b %08b\n", data[1], data[0])
		} else if len(data) == 4 {
			log.Printf("%08b %08b %08b %08b\n", data[3], data[2], data[1], data[0])
		} else {
			log.Panicln("")
		}
		var s uint64 = 0
		for i := 0; i < 32; i++ {
			s += c.Inner.GetRegister(i)
		}
		log.Println(i, c.Inner.GetPC(), s)

		if len(data) == 4 {
			var s uint64 = 0
			for i := len(data) - 1; i >= 0; i-- {
				s += uint64(data[i]) << (8 * i)
			}
			n, err := riscv.ExecuterRV64I(c.Inner, s)
			if err != nil {
				log.Panicln(err)
			}
			if n != 0 {
				i += 1
				continue
			}
		}
		log.Panicln("")
	}
}

var (
	cArgs = []string{"main"}
	cEnvs = []string{}
)

func main() {
	flag.Parse()

	inner := &riscv.CPU{}
	inner.SetMemory(riscv.NewMemoryLinear(4 * 1024 * 1024))
	inner.SetSystem(&riscv.SystemStandard{})
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
		if s.Flags&elf.SHF_ALLOC == 0 {
			continue
		}
		mem := make([]byte, s.Size)
		if _, err := s.ReadAt(mem, 0); err != nil {
			log.Panicln(err)
		}
		cpu.Inner.GetMemory().Set(s.Addr, mem)
	}
	cpu.Inner.SetRegister(riscv.Rsp, cpu.Inner.GetMemory().Len())

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

	addr := []uint64{0}
	// for i := len(cEnvs) - 1; i >= 0; i-- {
	// 	cpu.pushString(cEnvs[i])
	// 	addr = append(addr, cpu.ModuleBase.RG[riscv.Rsp])
	// }
	// addr = append(addr, 0)
	for i := len(cArgs) - 1; i >= 0; i-- {
		cpu.pushString(cArgs[i])
		addr = append(addr, cpu.Inner.GetRegister(riscv.Rsp))
	}
	// Align the stack to 8 bytes
	cpu.Inner.SetRegister(riscv.Rsp, cpu.Inner.GetRegister(riscv.Rsp)&^0x7)
	for _, a := range addr {
		cpu.pushUint64(a)
	}
	cpu.pushUint64(uint64(len(cArgs)))
	cpu.Run()
}
