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
	c.Inner.Register[riscv.Rsp] -= uint64(len(bs))
	for i, b := range bs {
		c.Inner.Memory[c.Inner.Register[riscv.Rsp]+uint64(i)] = b
	}
}

func (c *CPU) pushUint64(v uint64) {
	c.Inner.Register[riscv.Rsp] -= 8
	binary.LittleEndian.PutUint64(c.Inner.Memory[c.Inner.Register[riscv.Rsp]:c.Inner.Register[riscv.Rsp]+8], v)
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

var cStep = flag.Int64("steps", 200, "")

func (c *CPU) Run() {
	flag.Parse()
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	i := 0
	for {
		c.Inner.Register[riscv.Rzero] = 0x00
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
			s += c.Inner.Register[i]
		}
		log.Println(i, c.Inner.PC, s)
		if len(data) == 4 {
			var s uint64 = 0
			for i := len(data) - 1; i >= 0; i-- {
				s += uint64(data[i]) << (8 * i)
			}

			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[0], c.ModuleBase.RG[1], c.ModuleBase.RG[2], c.ModuleBase.RG[3])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[4], c.ModuleBase.RG[5], c.ModuleBase.RG[6], c.ModuleBase.RG[7])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[8], c.ModuleBase.RG[9], c.ModuleBase.RG[10], c.ModuleBase.RG[11])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[12], c.ModuleBase.RG[13], c.ModuleBase.RG[14], c.ModuleBase.RG[15])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[16], c.ModuleBase.RG[17], c.ModuleBase.RG[18], c.ModuleBase.RG[19])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[20], c.ModuleBase.RG[21], c.ModuleBase.RG[22], c.ModuleBase.RG[23])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[24], c.ModuleBase.RG[25], c.ModuleBase.RG[26], c.ModuleBase.RG[27])
			// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[28], c.ModuleBase.RG[29], c.ModuleBase.RG[30], c.ModuleBase.RG[31])
			n, err := riscv.ExecuterRV32I(c.Inner, s)
			if err != nil {
				log.Panicln(err)
			}
			if n != 0 {
				i += 1
				continue
			}
			n, err = riscv.ExecuterRV64I(c.Inner, s)
			if err != nil {
				log.Panicln(err)
			}
			if n != 0 {
				i += 1
				continue
			}
		}

		s = 0
		for i := len(data) - 1; i >= 0; i-- {
			s += uint64(data[i]) << (8 * i)
		}
		n, err := riscv.ExecuterC(c.Inner, s)
		if err != nil {
			log.Panicln(err)
		}
		if n != 0 {
			i += 1
			continue
		}
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[0], c.ModuleBase.RG[1], c.ModuleBase.RG[2], c.ModuleBase.RG[3])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[4], c.ModuleBase.RG[5], c.ModuleBase.RG[6], c.ModuleBase.RG[7])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[8], c.ModuleBase.RG[9], c.ModuleBase.RG[10], c.ModuleBase.RG[11])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[12], c.ModuleBase.RG[13], c.ModuleBase.RG[14], c.ModuleBase.RG[15])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[16], c.ModuleBase.RG[17], c.ModuleBase.RG[18], c.ModuleBase.RG[19])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[20], c.ModuleBase.RG[21], c.ModuleBase.RG[22], c.ModuleBase.RG[23])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[24], c.ModuleBase.RG[25], c.ModuleBase.RG[26], c.ModuleBase.RG[27])
		// log.Printf("%02x %02x %02x %02x\n", c.ModuleBase.RG[28], c.ModuleBase.RG[29], c.ModuleBase.RG[30], c.ModuleBase.RG[31])
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
			Register: [32]uint64{},
			Memory:   make([]byte, 4*1024*1024),
			PC:       0,
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
	cpu.Inner.Register[riscv.Rsp] = uint64(len(cpu.Inner.Memory))

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
		addr = append(addr, cpu.Inner.Register[riscv.Rsp])
	}
	// Align the stack to 8 bytes
	cpu.Inner.Register[riscv.Rsp] &^= 0x7
	for _, a := range addr {
		cpu.pushUint64(a)
	}
	cpu.pushUint64(uint64(len(cArgs)))
	cpu.Run()
}
