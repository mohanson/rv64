package rv64

var (
	_ System = (*SystemStandard)(nil)
)

type System interface {
	HandleCall(*CPU) (uint64, error)
	Code() uint8
}

type SystemStandard struct {
	ExitCode uint8
}

func (s *SystemStandard) HandleCall(c *CPU) (uint64, error) {
	code := c.GetRegister(Ra7)
	switch code {
	case 0x005d:
		s.ExitCode = uint8(c.GetRegister(Ra0))
		c.SetStatus(1)
		return 1, nil
	}
	return 0, ErrAbnormalEcall
}

func (s *SystemStandard) Code() uint8 {
	return s.ExitCode
}

func NewSystemStandard() *SystemStandard {
	return &SystemStandard{
		ExitCode: 0,
	}
}
