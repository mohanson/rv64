package rv64

import (
	"errors"
)

// https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf
// Chapter 20
// Table 20.1

const (
	Rzero = 0  // Hard-wired zero
	Rra   = 1  // Return address
	Rsp   = 2  // Stack pointer
	Rgp   = 3  // Global pointer
	Rtp   = 4  // Thread pointer
	Rt0   = 5  // Temporary/alternate link register
	Rt1   = 6  // Temporaries
	Rt2   = 7  // Temporaries
	Rs0fp = 8  // Saved register/frame pointer
	Rs1   = 9  // Saved register
	Ra0   = 10 // Function arguments/return values
	Ra1   = 11 // Function arguments/return values
	Ra2   = 12 // Function arguments
	Ra3   = 13 // Function arguments
	Ra4   = 14 // Function arguments
	Ra5   = 15 // Function arguments
	Ra6   = 16 // Function arguments
	Ra7   = 17 // Function arguments
	Rs2   = 18 // Saved registers
	Rs3   = 19 // Saved registers
	Rs4   = 20 // Saved registers
	Rs5   = 21 // Saved registers
	Rs6   = 22 // Saved registers
	Rs7   = 23 // Saved registers
	Rs8   = 24 // Saved registers
	Rs9   = 25 // Saved registers
	Rs10  = 26 // Saved registers
	Rs11  = 27 // Saved registers
	Rt3   = 28 // Temporaries
	Rt4   = 29 // Temporaries
	Rt5   = 30 // Temporaries
	Rt6   = 31 // Temporaries
)

const (
	Rft0  = 0  // FP temporaries
	Rft1  = 1  // FP temporaries
	Rft2  = 2  // FP temporaries
	Rft3  = 3  // FP temporaries
	Rft4  = 4  // FP temporaries
	Rft5  = 5  // FP temporaries
	Rft6  = 6  // FP temporaries
	Rft7  = 7  // FP temporaries
	Rfs0  = 8  // FP saved registers
	Rfs1  = 9  // FP saved registers
	Rfa0  = 10 // FP arguments/return values
	Rfa1  = 11 //  FP arguments/return values
	Rfa2  = 12 // FP arguments
	Rfa3  = 13 // FP arguments
	Rfa4  = 14 // FP arguments
	Rfa5  = 15 // FP arguments
	Rfa6  = 16 // FP arguments
	Rfa7  = 17 // FP arguments
	Rfs2  = 18 // FP saved registers
	Rfs3  = 19 // FP saved registers
	Rfs4  = 20 // FP saved registers
	Rfs5  = 21 // FP saved registers
	Rfs6  = 22 // FP saved registers
	Rfs7  = 23 // FP saved registers
	Rfs8  = 24 // FP saved registers
	Rfs9  = 25 // FP saved registers
	Rfs10 = 26 // FP saved registers
	Rfs11 = 27 // FP saved registers
	Rft8  = 28 // FP temporaries
	Rft9  = 29 // FP temporaries
	Rft10 = 30 // FP temporaries
	Rft11 = 31 // FP temporaries
)

const (
	CSRfflags  = 0x001 // Floating-Point Accrued Exceptions.
	CSRfrm     = 0x002 // Floating-Point Dynamic Rounding Mode.
	CSRfcsr    = 0x003 // Floating-Point Control and Status Register (frm + fflags).
	CSRcycle   = 0xc00 // Cycle counter for RDCYCLE instruction.
	CSRtime    = 0xc01 // Timer for RDTIME instruction.
	CSRinstret = 0xc02 // Instructions-retired counter for RDINSTRET instruction.
)

const (
	// Invalid Operation
	// This exception is raised if the given operands are invalid for the operation to be performed. Examples are
	// (see IEEE 754, section 7):
	// Addition or subtraction: &infin; - &infin;. (But &infin; + &infin; = &infin;).
	// Multiplication: 0 &middot; &infin;.
	// Division: 0/0 or &infin;/&infin;.
	// Remainder: x REM y, where y is zero or x is infinite.
	// Square root if the operand is less than zero. More generally, any mathematical function evaluated outside its
	// domain produces this exception.
	// Conversion of a floating-point number to an integer or decimal string, when the number cannot be represented in
	// the target format (due to overflow, infinity, or NaN).
	// Conversion of an unrecognizable input string.
	// Comparison via predicates involving < or >, when one or other of the operands is NaN. You can prevent this
	// exception by using the unordered comparison functions instead; see FP Comparison Functions.
	// If the exception does not trap, the result of the operation is NaN.
	FFlagsNV uint64 = 0x10
	// Division by Zero
	// This exception is raised when a finite nonzero number is divided by zero. If no trap occurs the result is
	// either +&infin; or -&infin;, depending on the signs of the operands.
	FFlagsDZ uint64 = 0x08
	// Overflow
	// This exception is raised whenever the result cannot be represented as a finite value in the precision format of
	// the destination. If no trap occurs the result depends on the sign of the intermediate result and the current
	// rounding mode (IEEE 754, section 7.3):
	// Round to nearest carries all overflows to &infin; with the sign of the intermediate result.
	// Round toward 0 carries all overflows to the largest representable finite number with the sign of the
	// intermediate result.
	// Round toward -&infin; carries positive overflows to the largest representable finite number and negative
	// overflows to -&infin;.
	// Round toward &infin; carries negative overflows to the most negative representable finite number and positive
	// overflows to &infin;.
	// Whenever the overflow exception is raised, the inexact exception is also raised.
	FFlagsOF uint64 = 0x04
	// Underflow
	// The underflow exception is raised when an intermediate result is too small to be calculated accurately, or if
	// the operation’s result rounded to the destination precision is too small to be normalized.
	// When no trap is installed for the underflow exception, underflow is signaled (via the underflow flag) only when
	// both tininess and loss of accuracy have been detected. If no trap handler is installed the operation continues
	// with an imprecise small value, or zero if the destination precision cannot hold the small exact result.
	FFlagsUF uint64 = 0x02
	// Inexact
	// This exception is signalled if a rounded result is not exact (such as when calculating the square root of two)
	// or a result overflows without an overflow trap.
	FFlagsNX uint64 = 0x01
)

var (
	ErrAbnormalEcall              = errors.New("Abnormal ecall")
	ErrAbnormalInstruction        = errors.New("Abnormal instruction")
	ErrMisalignedInstructionFetch = errors.New("Misaligned instruction fetch")
	ErrOutOfMemory                = errors.New("Out of memory")
	ErrReservedInstruction        = errors.New("Reserved instruction")
)

var (
	NaN32 uint32 = 0x7fc00000
	NaN64 uint64 = 0x7ff8000000000000
)

var (
	LogLevel = 0
)
