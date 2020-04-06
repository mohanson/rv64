package rv64

import (
	"log"
)

func (c *CPU) Run() uint8 {
	for {
		if c.GetStatus() == 1 {
			Debugln("Exit:", c.GetSystem().Code())
			return c.GetSystem().Code()
		}
		data, err := c.PipelineInstructionFetch()
		if err != nil {
			Panicln(err)
		}

		// Debugln("----------------------------------------")
		// var s uint64 = 0
		// for i := 0; i < 32; i++ {
		// 	s += c.GetRegister(uint64(i))
		// }
		// Debugln(fmt.Sprintf("nums=%d, pc=%d, sum=%d", c.GetCSR().Get(CSRinstret), c.GetPC(), s))
		// if len(data) == 2 {
		// 	Debugln(fmt.Sprintf("%08b %08b", data[1], data[0]))
		// } else if len(data) == 4 {
		// 	Debugln(fmt.Sprintf("%08b %08b %08b %08b", data[3], data[2], data[1], data[0]))
		// } else {
		// 	Panicln("")
		// }

		n, err := c.PipelineExecute(data)
		if err != nil {
			log.Panicln(err)
		}

		c.GetCSR().Set(CSRcycle, c.GetCSR().Get(CSRcycle)+n)
		c.GetCSR().Set(CSRtime, c.GetCSR().Get(CSRtime)+n)
		c.GetCSR().Set(CSRinstret, c.GetCSR().Get(CSRinstret)+1)
	}
}
