// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

import (
	"math/big"
	"testing"
	"github.com/copernet/copernicus/model/blockindex"
)

// TestBigToCompact ensures BigToCompact converts big integers to the expected
// compact representation.
func TestBigToCompact(t *testing.T) {
	tests := []struct {
		in  int64
		out uint32
	}{
		{0, 0},
		{-1, 25231360},
	}

	for x, test := range tests {
		n := big.NewInt(test.in)
		r := BigToCompact(n)
		if r != test.out {
			t.Errorf("TestBigToCompact test #%d failed: got %d want %d\n",
				x, r, test.out)
			return
		}
	}
}

// TestCompactToBig ensures CompactToBig converts numbers using the compact
// representation to the expected big intergers.
func TestCompactToBig(t *testing.T) {
	tests := []struct {
		in  uint32
		out int64
	}{
		{10000000, 0},
	}

	for x, test := range tests {
		n := CompactToBig(test.in)
		want := big.NewInt(test.out)
		if n.Cmp(want) != 0 {
			t.Errorf("TestCompactToBig test #%d failed: got %d want %d\n",
				x, n.Int64(), want.Int64())
			return
		}
	}
}

// TestCalcWork ensures CalcWork calculates the expected work value from values
// in compact representation.
func TestGetBlockProof(t *testing.T) {
	tests := []struct {
		in  uint32
		out int64
	}{
		{10000000, 0},
	}

	for x, test := range tests {
		bits := uint32(test.in)
		bi := new(blockindex.BlockIndex)
		bi.SetNull()
		bi.Header.Bits = bits

		r := GetBlockProof(bi)
		if r.Int64() != test.out {
			t.Errorf("TestCalcWork test #%d failed: got %v want %d\n",
				x, r.Int64(), test.out)
			return
		}
	}
}
