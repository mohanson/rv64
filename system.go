package riscv

var (
	_ System = (*SystemStandard)(nil)
)

type System interface {
	HandleCall(*CPU) (int, error)
}

type SystemStandard struct {
	ExitCode uint8
}

func (s *SystemStandard) HandleCall(c *CPU) (int, error) {
	code := c.GetRegister(Ra7)
	switch code {
	case 0x005d:
		s.ExitCode = uint8(c.GetRegister(Ra0))
		c.Stop = true
		return 1, nil
	}
	return 0, ErrAbnormalEcall
}
