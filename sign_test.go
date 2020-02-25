// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package riscv

import "testing"

func TestSignExtension(t *testing.T) {
	tests := []struct {
		in   uint64
		bit  int
		want uint64
	}{
		{in: 0x1, bit: 1, want: 0x1},
		{in: 0x3, bit: 1, want: 0xffffffffffffffff},
		{in: 0x1, bit: 0, want: 0xffffffffffffffff},
		{in: 0x800, bit: 11, want: 0xfffffffffffff800},
		{in: 0x7ff, bit: 11, want: 0x7ff},
		{in: 0x1000, bit: 12, want: 0xfffffffffffff000},
		{in: 0xfff, bit: 12, want: 0xfff},
		{in: 0x80000000, bit: 31, want: 0xffffffff80000000},
		{in: 0x7fffffff, bit: 31, want: 0x7fffffff},
	}
	for _, tt := range tests {
		got := signExtend(tt.in, tt.bit)
		if got != tt.want {
			t.Errorf("signExtend(%#x, %d) = %#x; want %#x", tt.in, tt.bit, got, tt.want)
		}
	}
}
