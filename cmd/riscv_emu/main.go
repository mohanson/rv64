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
		c.Inner.Memory[c.Inner.GetRegister(riscv.Rsp)+uint64(i)] = b
	}
}

func (c *CPU) pushUint64(v uint64) {
	c.Inner.SetRegister(riscv.Rsp, c.Inner.GetRegister(riscv.Rsp))
	binary.LittleEndian.PutUint64(c.Inner.Memory[c.Inner.GetRegister(riscv.Rsp):c.Inner.GetRegister(riscv.Rsp)+8], v)
}

func (c *CPU) FetchInstruction() []byte {
	if (c.Inner.PC + 2) > uint64(len(c.Inner.Memory)) {
		log.Panicln("Out of memory")
	}
	a := c.Inner.Memory[c.Inner.PC : c.Inner.PC+2]
	b := riscv.InstructionLengthEncoding(a)
	instructionBytes := c.Inner.Memory[c.Inner.PC : c.Inner.PC+uint64(b)]
	return instructionBytes
}

var cStep = flag.Int64("steps", 250, "")

func (c *CPU) Run() {
	flag.Parse()
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	i := 0
	for {
		if c.Inner.Stop {
			log.Println("Exit:", c.Inner.System.(*riscv.SystemStandard).ExitCode)
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
		log.Println(i, c.Inner.PC, s)

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

	cpu := &CPU{
		Inner: &riscv.CPU{
			System: &riscv.SystemStandard{},
			Memory: make([]byte, 4*1024*1024),
		},
	}

	f, err := elf.Open(flag.Arg(0))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	cpu.Inner.PC = f.Entry

	for _, s := range f.Sections {
		if s.Flags&elf.SHF_ALLOC == 0 {
			continue
		}
		if _, err := s.ReadAt(cpu.Inner.Memory[s.Addr:s.Addr+s.Size], 0); err != nil {
			log.Panicln(err)
		}
	}
	cpu.Inner.SetRegister(riscv.Rsp, uint64(len(cpu.Inner.Memory)))

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
