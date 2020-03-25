package rv64

import (
	"math"
)

// NaN Boxing of Narrower Values
// When multiple ﬂoating-point precisions are supported, then valid values of narrower n-bit types, n < FLEN, are
// represented in the lower n bits of an FLEN-bit NaN value, in a process termed NaN-boxing. The upper bits of a valid
// NaN-boxed value must be all 1s. Valid NaN-boxed n-bit values therefore appear as negative quiet NaNs (qNaNs) when
// viewed as any wider m-bit value, n < m ≤ FLEN.
//
// Floating-point n-bit transfer operations move external values held in IEEE standard formats into and out of the f
// registers, and comprise ﬂoating-point loads and stores (FLn/FSn) and ﬂoatingpoint move
// instructions (FMV.n.X/FMV.X.n). A narrower n-bit transfer, n < FLEN, into the f registers will create a valid
// NaN-boxed value by setting all upper FLEN−n bits of the destination f register to 1. A narrower n-bit transfer out
// of the ﬂoating-point registers will transfer the lower n bits of the register ignoring the upper FLEN−n bits.
// Floating-point compute and sign-injection operations calculate results based on the FLEN-bit values held in the
// f registers. A narrow n-bit operation, where n < FLEN, checks that input operands are correctly NaN-boxed, i.e.,
// all upper FLEN−n bits are 1. If so, the n least-signiﬁcant bits of the input are used as the input value, otherwise
// the input value is treated as an n-bit canonical NaN. An n-bit ﬂoating-point result is written to the n
// least-signiﬁcant bits of the destination f register, with all 1s written to the uppermost FLEN−n bits to yield a
// legal NaN-boxed value. Conversions from integer to ﬂoating-point (e.g., FCVT.S.X), will NaN-box any results
// narrower than FLEN to ﬁll the FLEN-bit destination register. Conversions from narrower n-bit ﬂoatingpoint
// values to integer (e.g., FCVT.X.S) will check for legal NaN-boxing and treat the input as the n-bit canonical NaN
// if not a legal n-bit value.

func NaNBoxing(f float32) float64 {
	return math.Float64frombits(0xffffffff00000000 | uint64(math.Float32bits(f)))
}

func NaNGnixob(f float64) float32 {
	u := math.Float64bits(f)
	// The n least-significant bits of the input are used as the input value, otherwise the input value is treated as
	// an n-bit canonical NaN.
	if (u >> 32) != 0xffffffff {
		return math.Float32frombits(NaN32)
	}
	return math.Float32frombits(uint32(u))
}
