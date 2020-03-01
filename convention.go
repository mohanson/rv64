package riscv

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

var (
	ErrAbnormalInstruction = errors.New("Abnormal instruction")
	ErrReservedInstruction = errors.New("Reserved instruction")
	ErrAbnormalEcall       = errors.New("Abnormal ecall")
)
