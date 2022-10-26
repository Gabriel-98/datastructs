// LICENCE NOT YET DEFINED.

// coverage: 98.9% of statements
//
// The next functions cannot be fully covered in the tests because the
// execution of some statements depends on the type of architecture on
// which the code will be executed (32 or 64 bits). The functions with
// these limitations of coverage are:
//
//     - IntToBitBlock
//     - UintToBitBlock
//
package bitblock


import (
	"testing"
	"unsafe"
)


// checkBitBlockSize checks that bitBlock.Size() returns the correct size,
// by comparing it to the size it should be. If the value returned by
// bitBlock.Size() is incorrect, an error describing it will be printed.
func checkBitBlockSize(t *testing.T, bitBlock *BitBlock, correctSize int) bool {
	if s := bitBlock.Size(); s != correctSize {
		t.Errorf("got bitBlock.Size() = %d, want bitBlock.Size() = %d", s, correctSize)
		return false
	}
	return true
}

// checkBitBlockValues checks that bitBlock has set the bits to the
// correct value, by comparing it to the values of a slice of bools; if
// not, an error describing it will be printed.
func checkBitBlockValues(t *testing.T, bitBlock *BitBlock, correct []bool) bool {
	if !checkBitBlockSize(t, bitBlock, len(correct)) {
		return false
	}
	for i := 0; i < len(correct); i++ {
		if b := bitBlock.Get(i); b != correct[i] {
			t.Errorf("got bitBlock.Get(%d) = %t, want bitBlock.Get(%d) = %t", i, b, i, correct[i])
			return false
		}
	}
	return true
}

// checkPaddingBits checks that the padding bits are set to 0.
// If any of the padding bits are set to 1, an error describing
// it will be printed.
func checkPaddingBits(t *testing.T, bitBlock *BitBlock) bool {
	bitBlock2 := BytesToBitBlock(bitBlock.bits, len(bitBlock.bits) * 8)
	for i := bitBlock.Size(); i < bitBlock2.Size(); i++ {
		if bitBlock2.Get(i) {
			t.Errorf("the %d-th padding bit of bitBlock is true, padding bits must be set to false", i - bitBlock.Size())
			return false
		}
	}
	return true
}

// Test the functions to create a new BitBlock: NewZeroBitBlock and BytesToBitBlock.
// Test the BitBlock methods: Get, Set0, Set1, Set, ToBinaryString, RemoveFirstBits,
// RemoveLastBits, GetSubBlock and ToBytes.
func TestBitBlock(t *testing.T) {
	// Binary string of 135 characters.
	s := "010010111100110001101010110111110010100110001001111000110010101000101010101010011100011110010101011100000111001101101011001000110000010"

	// Test the NewZeroBitBlock() function.
	t.Run("NewZeroBitBlock", func(t *testing.T) {
		for _, size := range []int{0, 1, 8, 10, 16, 32, 64, 128, 170, 15, 90, 256, 200, 250, 10} {
			bools := make([]bool, size)
			bitBlock := NewZeroBitBlock(size)
			if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
				t.Fatalf("the BitBlock obtained by calling NewZeroBitBlock(%d) is wrong, want a BitBlock with all bits set to 0", size)
			}
		}
		for _, size := range []int{-1, -9, -10, -15, -100, -200, -128, -256, -64, -31} {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to NewZeroBitBlock(%d) did not panic", size)
					}
				}()
				NewZeroBitBlock(size)
			}()
		}
	})

	// Creates a BitBlock with all bits set to 0 and a bool slice to
	// validate that the bits in the BitBlock match the correct values
	// after each call to a BitBlock method.
	bools := make([]bool, len(s))
	bitBlock := NewZeroBitBlock(len(s))
	if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
		t.Fatalf("the BitBlock obtained by calling NewZeroBitBlock(%d) is wrong, want a BitBlock with all bits set to 0", len(s))
	}
	
	// Initialize the BitBlock with the values of s.
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			bitBlock.Set(i, false)
			bools[i] = false
		} else {
			bitBlock.Set(i, true)
			bools[i] = true
		}
		if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
			t.Fatalf("incosistency after call Set(%d, %t)", i, bools[i])
		}
	}

	// Test the Get() method.
	// This test is just to check that Get() panics if an invalid position
	// is passed for the corresponding BitBlock.
	for _, pos := range []int{135, -1, -8, -16, -64, -32, 138, 145, 162, 200, 256} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to Get(%d) on a BitBlock of size %d did not panic", pos, bitBlock.Size())
				}
			}()
			bitBlock.Get(pos)
		}()
	}

	// Test the Set1() method.
	for _, pos := range []int{7, 1, 13, 114, 120, 127, 15, 13, 0, 8, 64, 43, 63, 132, 63, 48} {
		bitBlock.Set1(pos); bools[pos] = true
		if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
			t.Fatalf("incosistency after call Set1(%d)", pos)
		}
	}
	for _, pos := range []int{170, -50, 135, 137, -4, -10, -8, 142, 160} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to Set1(%d) on a BitBlock of size = %d did not panic", pos, bitBlock.Size())
				}
			}()
			bitBlock.Set1(pos)
		}()
	}

	// Test the Set0() method.
	for _, pos := range []int{50, 7, 127, 9, 0, 0, 8, 13, 9, 2, 9, 83, 63, 134} {
		bitBlock.Set0(pos); bools[pos] = false
		if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
			t.Fatalf("incosistency after call Set0(%d)", pos)
		}
	}
	for _, pos := range []int{150, 135, -10, -1, 200, 193, -4, -30, 171} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to Set0(%d) on a BitBlock of size = %d did not panic", pos, bitBlock.Size())
				}
			}()
			bitBlock.Set0(pos)
		}()
	}

	// Test the Set() method.
	func() {
		type Update struct { pos int; value bool }
		updates := []Update{ Update{0, true}, Update{58, true}, Update{134, true}, Update{134, false}, Update{134, true}, Update{120, false} }
		for _, update := range updates {
			pos, value := update.pos, update.value
			bitBlock.Set(pos, value); bools[pos] = value
			if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
				t.Fatalf("incosistency after call Set(%d, %t)", pos, value)
			}
		}
	}()
	func() {
		type Update struct { pos int; value bool }
		updates := []Update{ Update{135, false}, Update{140, false}, Update{-9, false}, Update{150, true} }
		for _, update := range updates {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to Set(%d, %t) on a BitBlock of size = %d did not panic", update.pos, update.value, bitBlock.Size())
					}
				}()
				bitBlock.Set(update.pos, update.value)
			}()
		}
	}()

	// Test the ToBinaryString() method.
	t.Run("ToBinaryString", func(t *testing.T) {
		binstr := bitBlock.ToBinaryString()
		for i := 0; i < len(binstr); i++ {
			if bools[i] {
				if binstr[i] == '0' {
					t.Fatalf("ToBinaryString() returns an invalid binary string, got binstr[%d] = %c and bools[%d] = %t, want the value of binstr[%d] and the value of bools[%d] to be equivalent", i, binstr[i], i, bools[i], i, i)
				}
			} else {
				if binstr[i] == '1' {
					t.Fatalf("ToBinaryString() returns an invalid binary string, got binstr[%d] = %c and bools[%d] = %t, want the value of binstr[%d] and the value of bools[%d] to be equivalent", i, binstr[i], i, bools[i], i, i)
				}
			}
		}
	})

	// Test the RemoveFirstBits() method.
	t.Run("RemoveFirstBits", func(t *testing.T) {
		for k := 0; k <= len(s); k++ {
			bitBlock2 := bitBlock.RemoveFirstBits(k)
			if ok := checkBitBlockValues(t, bitBlock2, bools[k:]); !ok {
				t.Fatalf("inconsistency found after calling bitBlock.RemoveFirstBits(%d)", k)
			}
			if k < len(s) {
				if bitBlock2.Get(0) {
					bitBlock2.Set(0, false)
					if bitBlock2.Get(0) == bitBlock.Get(k) {
						t.Fatalf("problem found after calling bitBlock.RemoveFirstBits(%d), modifying position 0 in bitBlock2 had effect on value at position k = %d of bitBlock", k, k)
					}
					bitBlock2.Set(0, true)
				} else {
					bitBlock2.Set(0, true)
					if bitBlock2.Get(0) == bitBlock.Get(k) {
						t.Fatalf("problem found after calling bitBlock.RemoveFirstBits(%d), modifying position 0 in bitBlock2 had effect on value at position k = %d of bitBlock", k, k)
					}
					bitBlock2.Set(0, false)
				}
			}
		}
		for _, k := range []int{-1, 136, 170, -8} {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to RemoveFirstBits(%d) on a BitBlock of size %d did not panic", k, bitBlock.Size())
					}
				}()
				bitBlock.RemoveFirstBits(k)
			}()
		}
	})

	// Test the RemoveLastBits() method.
	t.Run("RemoveLastBits", func(t *testing.T) {
		for k := 0; k <= len(s); k++ {
			bitBlock2 := bitBlock.RemoveLastBits(k)
			if ok := checkBitBlockValues(t, bitBlock2, bools[:len(bools)-k]); !ok {
				t.Fatalf("inconsistency found after calling bitBlock.RemoveLastBits(%d)", k)
			}
			if k < len(s) {
				if bitBlock2.Get(0) {
					bitBlock2.Set(0, false)
					if bitBlock2.Get(0) == bitBlock.Get(0) {
						t.Fatalf("problem found after calling bitBlock.RemoveLastBits(%d), modifying position 0 in bitBlock2 had effect on the value at position 0 of bitBlock", k)
					}
					bitBlock2.Set(0, true)
				} else {
					bitBlock2.Set(0, true)
					if bitBlock2.Get(0) == bitBlock.Get(0) {
						t.Fatalf("problem found after calling bitBlock.RemoveLastBits(%d), modifying position 0 in bitBlock2 had effect on the value at position 0 of bitBlock", k)
					}
					bitBlock2.Set(0, false)
				}
			}
		}
		for _, k := range []int{200, -5, 140} {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to RemoveLastBits(%d) on a BitBlock of size %d did not panic", k, bitBlock.Size())
					}
				}()
				bitBlock.RemoveLastBits(k)
			}()
		}
	})

	// Test the GetSubBlock() method.
	t.Run("GetSubBlock", func(t *testing.T) {
		for l := 0; l <= len(s); l++ {
			for r := l; r <= len(s); r++ {
				bitBlock2 := bitBlock.GetSubBlock(l,r)
				if ok := checkBitBlockValues(t, bitBlock2, bools[l:r]); !ok {
					t.Fatalf("inconsistency found after calling bitBlock.GetSubBlock(%d,%d)", l, r)
				}
				if l < r {
					if bitBlock2.Get(0) {
						bitBlock2.Set(0, false)
						if bitBlock2.Get(0) == bitBlock.Get(l) {
							t.Fatalf("problem found after calling bitBlock.GetSubBlock(%d,%d), modifying position 0 in bitBlock2 had effect on the value at position l = %d of bitBlock", l, r, l)
						}
						bitBlock2.Set(0, true)
					} else {
						bitBlock2.Set(0, true)
						if bitBlock2.Get(0) == bitBlock.Get(l) {
							t.Fatalf("problem found after calling bitBlock.GetSubBlock(%d,%d), modifying position 0 in bitBlock2 had effect on the value at position l = %d of bitBlock", l, r, l)
						}
						bitBlock2.Set(0, false)
					}
				}
			}
		}
		
		type Range struct { start int; end int }
		ranges := []Range{ Range{100, 136}, Range{-7, -2}, Range{0, 140}, Range{-20, 15}, Range{-9, 240}, Range{2, 170}, Range{90, 140} }
		for _, r := range ranges {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to GetSubBlock(%d, %d) on a BitBlock of size %d did not panic", r.start, r.end, bitBlock.Size())
					}
				}()
				bitBlock.GetSubBlock(r.start, r.end)
			}()
		}
	})

	// Test the BytesToBitBlock() function.
	t.Run("BytesToBitBlock", func(t *testing.T) {
		bitBlock2 := BytesToBitBlock(bitBlock.bits, 70)
		if ok := checkBitBlockValues(t, bitBlock2, bools[:70]); !ok {
			t.Fatalf("the BitBlock obtained by calling BytesToBitBlock(bitBlock.bits, %d) is wrong, want a BitBlock with the first %d bits equal to the first %d bits of bitBlock", 70, 70, 70)
		}
		if bitBlock2.Get(0) {
			bitBlock2.Set(0, false)
			if bitBlock2.Get(0) == bitBlock.Get(0) {
				t.Fatalf("problem found with the BitBlock returned by BytesToBitBlock(bitBlock.bits, %d), modifying position 0 in bitBlock2 had effect on the value at position 0 of bitBlock", 70)
			}
			bitBlock2.Set(0, true)
		} else {
			bitBlock2.Set(0, true)
			if bitBlock2.Get(0) == bitBlock.Get(0) {
				t.Fatalf("problem found with the BitBlock returned by BytesToBitBlock(bitBlock.bits, %d), modifying position 0 in bitBlock2 had effect on the value at position 0 of bitBlock", 70)
			}
			bitBlock2.Set(0, false)
		}
		func() {
			type Update struct { pos int; value bool }
			updates := []Update{ Update{70, true}, Update{-5, false}, Update{-8, false}, Update{100, true}, Update{128, true}, Update{77, false} }
			for _, update := range updates {
				pos, value := update.pos, update.value
				func() {
					defer func() {
						panicMessage := recover()
						if panicMessage == nil {
							t.Fatalf("the call to Set(%d, %t) on bitBlock2 which has size = %d did not panic", pos, value, bitBlock2.Size())
						}
					}()
					bitBlock2.Set(pos, value)
				}()
			}
		}()
		for _, pos := range []int{70, -5, -1, -8, 100, 128, 127, 82, 75, -28} {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to Get(%d) on bitBlock2 which has size = %d did not panic", pos, bitBlock2.Size())
					}
				}()
				bitBlock2.Get(pos)
			}()
		}
		for _, size := range []int{-1, -9, -10, -15, -100, -200, -128, -256, -64, -31} {
			func() {
				defer func() {
					panicMessage := recover()
					if panicMessage == nil {
						t.Fatalf("the call to BytesToBitBlock(bitBlock.bits, %d) did not panic", size)
					}
				}()
				BytesToBitBlock(bitBlock.bits, size)
			}()
		}
	})

	// Test the ToBytes() method.
	t.Run("ToBytes", func(t *testing.T) {
		bytes := bitBlock.ToBytes()
		bitBlock2 := BytesToBitBlock(bytes, bitBlock.Size())
		bytes2 := bitBlock2.ToBytes()

		if len(bytes) != len(bytes2) {
			t.Fatalf("len(bytes) = %d is different from len(bytes2) = %d, want len(bytes) = len(bytes2)", len(bytes), len(bytes2))
		}
		if len(bytes) > 0 {
			if &bytes[0] == &bytes2[0] {
				t.Fatalf("address of bytes[0] = %p is the same as the address of bytes2[0] = %p, want the address of bytes[0] to be different than the address of bytes2[0]", bytes, bytes2)
			}
		}
		for i := 0; i < len(bytes); i++ {
			if bytes[i] != bytes2[i] {
				t.Fatalf("bytes[%d] = %d is different from bytes2[%d] = %d, want bytes[%d] = bytes2[%d]", i, bytes[i], i, bytes2[i], i, i)
			}
		}
	})
}

// Test the Clone() method of the BitBlock type.
func TestBitBlockClone(t *testing.T) {
	type Test struct{ id string; size int; bytes []byte }

	// Test cases.	
	tests := []Test{
		Test{ id: "0000", size: 82, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0001", size: 103, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0002", size: 88, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0003", size: 65, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0004", size: 3, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0005", size: 0, bytes: []byte{45, 232, 0, 1, 245, 87, 255, 1, 64, 127, 184} },
		Test{ id: "0006", size: 20, bytes: []byte{} },
		Test{ id: "0007", size: 16, bytes: []byte{} },
		Test{ id: "0008", size: 0, bytes: []byte{} },
		Test{ id: "0009", size: 3, bytes: []byte{} },
	}

	for _, test := range tests {
		// Slice of bytes that will be used to create the original BitBlock
		// and the size that it will be.
		bytes := test.bytes
		size := test.size
		t.Run(test.id, func(t *testing.T) {
			// Create a BitBlock from a slice of bytes, and from this BitBlock
			// another is created using the Clone() method.
			bitBlock := BytesToBitBlock(bytes, size)
			clonedBitBlock := bitBlock.Clone()

			// Test that the BitBlock returned by the Clone() method will be
			// different of nil and the bits contained in an slice will also be
			// different of nil.
			if clonedBitBlock == nil {
				t.Fatalf("clonedBitBlock = nil, want clonedBitBlock != nil")
			}
			if clonedBitBlock.bits == nil {
				t.Fatalf("clonedBitBlock.bits = nil, want clonedBitBlock.bits != nil")
			}

			// Test that the length of bits in bitBlock and clonedBitBlock is
			// the same.
			if len(bitBlock.bits) != len(clonedBitBlock.bits) {
				t.Fatalf("len(bitBlock.bits) = %d is different from len(clonedBitBlock.bits) = %d, want len(bitBlock.bits) = len(clonedBitBlock.bits)", len(bitBlock.bits), len(clonedBitBlock.bits))
			}

			// Test that bitBlock and clonedBitBlock have the same size.
			if s1, s2 := bitBlock.Size(), clonedBitBlock.Size(); s1 != s2 {
				t.Fatalf("bitBlock.Size() = %d is different of clonedBitBlock.Size() = %d, want bitBlock.Size() = clonedBitBlock.Size()", s1, s2)
			}

			// Test that the bits of the original BitBlock and those of the
			// cloned BitBlock are stored at different addresses.
			if len(bitBlock.bits) > 0 {
				if &bitBlock.bits[0] == &clonedBitBlock.bits[0] {
					t.Fatalf("address of bitBlock.bits (%p) = address of clonedBitBlock.bits (%p), want address of bitBlock.bits != address of clonedBitBlock.bits", bitBlock.bits, clonedBitBlock.bits)
				}
			}

			// Test that the original BitBlock and the cloned BitBlock have
			// the bits set to the same values. The comparison is done at
			// the level of the underlying structure and also by comparing 
			// each bit by calling the Get function for both BitBlocks.
			for i := 0; i < len(bitBlock.bits); i++ {
				if bitBlock.bits[i] != clonedBitBlock.bits[i] {
					t.Fatalf("bitBlock.bits[%d] = %d is different of clonedBitBlock.bits[%d] = %d, want bitBlock.bits[%d] = clonedBitBlock.bits[%d]", i, bitBlock.bits[i], i, clonedBitBlock.bits[i], i, i)
				}
			}
			for i := 0; i < bitBlock.Size(); i++ {
				if b1, b2 := bitBlock.Get(i), clonedBitBlock.Get(i); b1 != b2 {
					t.Fatalf("bitBlock.Get(%d) = %t is different of clonedBitBlock.Get(%d) = %t, want bitBlock.Get(%d) = clonedBitBlock.Get(%d)", i, b1, i, b2, i, i)
				}
			}

			// Test whether a panic is raised when trying to get the value of
			// an invalid position for clonedBitBlock.
			for _, pos := range []int{clonedBitBlock.Size(), clonedBitBlock.Size() + 3, -1, -15} {
				func() {
					defer func() {
						panicMessage := recover()
						if panicMessage == nil {
							t.Fatalf("inconsistency with clonedBitBlock of size = %d, the call to Get(%d) on clonedBitBlock did not panic", clonedBitBlock.Size(), pos)
						}
					}()
					clonedBitBlock.Get(pos)
				}()
			}
		})
	}
}

// Test the Concatenate() function.
func TestBitBlockConcatenate(t *testing.T) {
	type Test struct { id string; s string; sizes []int }
	
	// Test cases.
	tests := []Test{
		Test{
			id: "0000",
			s: "1101001000011010111101010111010101110110101000011110111111010111110111110110001011011111010101101110101010110101101111110000111101010101011",
			sizes: []int{31, 42, 24, 29, 13},
		},
		Test{
			id: "0001",
			s: "01010101100101010101111000101010101",
			sizes: []int{14, 8, 13},
		},
		Test{
			id: "0002",
			s: "011",
			sizes: []int{3},
		},
		Test{
			id: "0003",
			s: "1101011010111101010101110001100101010101010101010111100001110010",
			sizes: []int{8,15,7,18,1,8,5,2},
		},
		Test{
			id: "0004",
			s: "",
			sizes: []int{},
		},
		Test{
			id: "0005",
			s: "10101001",
			sizes: []int{8},
		},
		Test{
			id: "0006",
			s: "01010111",
			sizes: []int{2,1,3,2},
		},
	}

	for _, test := range tests {
		// The binary string and the sizes that each of the sub BitBlocks will have.
		// The sum of the sizes of the sub BitBlocks will be equal to the length
		// of s.
		s := test.s
		sizes := test.sizes

		t.Run(test.id, func(t *testing.T) {
			// A slice of bools is created to be equivalent to the binary string
			// s and so that it can be tested whether the BitBlock has set the
			// correct values just by comparing it to that slice of bools.
			bools := make([]bool, len(s))
			for i := 0; i < len(s); i++ {
				if s[i] == '1' {
					bools[i] = true
				}
			}

			// Some BitBlocks are created such that each has bits corresponding
			// to a substring of s, and all BitBlocks concatenated in the same
			// order are equivalent to the bits represented by s.
			bitBlocks := make([]*BitBlock, len(sizes))
			for i := 0; i < len(bitBlocks); i++ {
				bitBlocks[i] = NewZeroBitBlock(sizes[i])
			}
			for i, accu := 0, 0; i < len(bitBlocks); i++ {
				for j := 0; j < bitBlocks[i].Size(); j, accu = j+1, accu+1 {
					if s[accu] == '0' {
						bitBlocks[i].Set(j, false)
					} else {
						bitBlocks[i].Set(j, true)
					}
				}
			}

			// Concatenates the BitBlocks in the same order and checks that
			// the BitBlock obtained is correct.
			bitBlock := Concatenate(bitBlocks...)
			if ok := checkBitBlockValues(t, bitBlock, bools); !ok {
				t.Fatalf("wrong concatenation of BitBlocks, want bits of bitBlocks[0] + bitBlocks[1] + bitBlocks[2] + ... = bits of bitBlock")
			}
		})
	}
}

// Test that the padding bits of the BitBlock returned by some method or function will
// be all set to false.
// The functions tested here are: NewZeroBitBlock, BytesToBitBlock and Concatenate.
// The methods tested here are: RemoveFirstBits, RemoveLastBits, GetSubBlock and Clone.
func TestPaddingBits(t *testing.T) {
	// Test the NewZeroBitBlock() function.
	t.Run("NewZeroBitBlock", func(t *testing.T) {
		for size := 0; size <= 200; size++ {
			bitBlock := NewZeroBitBlock(size)
			if ok := checkPaddingBits(t, bitBlock); !ok {
				t.Fatalf("the call to NewZeroBitBlock(%d) returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", size)
			}
		}
	})

	// Test the BytesToBitBlock() function.
	t.Run("BytesToBitBlock", func(t *testing.T) {
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size := 0; size <= 200; size++ {
			bitBlock := BytesToBitBlock(bytes, size)
			if ok := checkPaddingBits(t, bitBlock); !ok {
				t.Fatalf("the call to BytesToBitBlock(%v, %d) returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", bytes, size)
			}
 		}
	})

	// Test the RemoveFirstBits() method.
	t.Run("RemoveFirstBits", func(t *testing.T) {
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size := 0; size <= 8 * len(bytes); size++ {
			bitBlock := BytesToBitBlock(bytes, size)
			for k := 0; k <= size; k++ {
				bitBlock2 := bitBlock.RemoveFirstBits(k)
				if ok := checkPaddingBits(t, bitBlock2); !ok {
					t.Fatalf("the call to RemoveFirstBits(%d) on a BitBlock of size %d returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", k, size)
				}
			}
		}
	})

	// Test the RemoveLastBits() method.
	t.Run("RemoveLastBits", func(t *testing.T){
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size := 0; size <= 8 * len(bytes); size++ {
			bitBlock := BytesToBitBlock(bytes, size)
			for k := 0; k <= size; k++ {
				bitBlock2 := bitBlock.RemoveLastBits(k)
				if ok := checkPaddingBits(t, bitBlock2); !ok {
					t.Fatalf("the call to RemoveLastBits(%d) on a BitBlock of size %d returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", k, size)
				}
			}
		}
	})

	// Test the GetSubBlock() method.
	t.Run("GetSubBlock", func(t *testing.T){
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size := 0; size <= 8 * len(bytes); size++ {
			bitBlock := BytesToBitBlock(bytes, size)
			for l := 0; l <= size; l++ {
				for r := l; r <= size; r++ {
					bitBlock2 := bitBlock.GetSubBlock(l,r)
					if ok := checkPaddingBits(t, bitBlock2); !ok {
						t.Fatalf("the call to GetSubBlock(%d, %d) on a BitBlock of size %d returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", l, r, size)
					}
				}
			}
		}
	})

	// Test the Clone() method.
	t.Run("Clone", func(t *testing.T) {
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size := 0; size <= 8 * len(bytes); size++ {
			bitBlock := BytesToBitBlock(bytes, size)
			bitBlock2 := bitBlock.Clone()
			if ok := checkPaddingBits(t, bitBlock2); !ok {
				t.Fatalf("the call to Clone() on a BitBlock of size %d returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", size)
			}
		}
	})

	// Test the Concatenate() function.
	t.Run("Concatenate", func(t *testing.T) {
		bytes := []byte{15, 54, 127, 200, 0, 15, 95, 128, 127, 34, 19, 183, 255}
		for size1 := 0; size1 <= 100; size1++ {
			bitBlock := BytesToBitBlock(bytes, size1)
			for size2 := 0; size2 <= 25; size2++ {
				bitBlock2 := BytesToBitBlock(bytes, size2)
				for size3 := 0; size3 <= 100; size3++ {
					bitBlock3 := BytesToBitBlock(bytes, size3)
					bitBlock4 := Concatenate(bitBlock, bitBlock2, bitBlock3)
					if ok := checkPaddingBits(t, bitBlock4); !ok {
						t.Fatalf("the call to Concatenate() for BitBlocks of sizes %d, %d and %d returned a BitBlock with some padding bits set to true, want all the padding bits of the BitBlocks to be set to false", size1, size2, size3)
					}
				}
			}
		}
	})
}

// Test the functions to convert between integer numbers and BitBlock:
// - Uint8ToBitBlock, Uint16ToBitBlock, Uint32ToBitBlock, Uint64ToBitBlock, UintToBitBlock.
// - Int8ToBitBlock, Int16ToBitBlock, Int32ToBitBlock, Int64ToBitBlock, IntToBitBlock.
// - BitBlockToUint8, BitBlockToUint16, BitBlockToUint32, BitBlockToUint64, BitBlockToUint.
// - BitBlockToInt8, BitBlockToInt16, BitBlockToInt32, BitBlockToInt64, BitBlockToInt.
func TestConversionBetweenIntegerNumbersAndBitBlocks(t *testing.T) {
	// Test the functions Uint8ToBitBlock and BitBlockToUint8.
	uint8Numbers := []uint8{9, 0, 54, 128, 255, 127, 64, 65, 230, 1}
	for _, x := range uint8Numbers {
		bitBlock := Uint8ToBitBlock(x)
		var x2 uint8 = 0
		for i := 7; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Uint8ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToUint8(bitBlock)
		if x != x3 {
			t.Fatalf("the uint8 obtained by calling BitBlockToUint8(%s) = %d is wrong, want BitBlockToUint8(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}	
	for _, size := range []int{0, 7, 9, 10, 15, 16} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToUint8(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToUint8() with a BitBlock of size != 8 should panic", bitBlock.Size())
				}
			}()
			BitBlockToUint8(bitBlock)
		}()
	}
	
	// Test the functions Uint16ToBitBlock and BitBlockToUint16.
	uint16Numbers := []uint16{541, 3, 8941, 0, 65535, 65534, 3165, 65415, 19811, 132, 127}
	for _, x := range uint16Numbers {
		bitBlock := Uint16ToBitBlock(x)
		var x2 uint16 = 0
		for i := 15; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Uint16ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToUint16(bitBlock)
		if x != x3 {
			t.Fatalf("the uint16 obtained by calling BitBlockToUint16(%s) = %d is wrong, want BitBlockToUint16(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{15, 17, 24, 19, 8, 7, 13} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToUint16(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToUint16() with a BitBlock of size != 16 should panic", bitBlock.Size())
				}
			}()
			BitBlockToUint16(bitBlock)
		}()
	}

	// Test the functions Uint32ToBitBlock and BitBlockToUint32.
	uint32Numbers := []uint32{615, 8942132, 12314531, 0, 1, 41, 4294967295, 4294967114, 2015617849, 654641, 1064541654, 2313516547}
	for _, x := range uint32Numbers {
		bitBlock := Uint32ToBitBlock(x)
		var x2 uint32 = 0
		for i := 31; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Uint32ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToUint32(bitBlock)
		if x != x3 {
			t.Fatalf("the uint32 obtained by calling BitBlockToUint32(%s) = %d is wrong, want BitBlockToUint32(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{30, 35, 40, 64, 0, 8, 1, 15} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToUint32(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToUint32() with a BitBlock of size != 32 should panic", bitBlock.Size())
				}
			}()
			BitBlockToUint32(bitBlock)
		}()
	}

	// Test the functions Uint64ToBitBlock and BitBlockToUint64.
	uint64Numbers := []uint64{654151, 21332, 231, 0, 1, 46513321, 18446744073709551615, 18446744073709550923, 5645651, 374, 6517}
	for _, x := range uint64Numbers {
		bitBlock := Uint64ToBitBlock(x)
		var x2 uint64 = 0
		for i := 63; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Uint64ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToUint64(bitBlock)
		if x != x3 {
			t.Fatalf("the uint64 obtained by calling BitBlockToUint64(%s) = %d is wrong, want BitBlockToUint64(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{150, 200, 128, 32, 8, 0, 7, 15, 49, 63, 65} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToUint64(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToUint64() with a BitBlock of size != 64 should panic", bitBlock.Size())
				}
			}()
			BitBlockToUint64(bitBlock)
		}()
	}

	// Test the functions UintToBitBlock and BitBlockToUint.
	var uintNumbers []uint
	if int(unsafe.Sizeof(uint(0))) == 4 {
		uintNumbers = []uint{615, 8942132, 12314531, 0, 1, 41, 4294967295, 4294967114, 2015617849, 654641, 1064541654, 2313516547}
	} else {
		uintNumbers = []uint{654151, 21332, 231, 0, 1, 46513321, 18446744073709551615, 18446744073709550923, 5645651, 374, 6517}
	}
	for _, x := range uintNumbers {
		bitBlock := UintToBitBlock(x)
		var x2 uint = 0
		for i := int(unsafe.Sizeof(uint(0))) * 8 - 1; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling UintToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToUint(bitBlock)
		if x != x3 {
			t.Fatalf("the uint obtained by calling BitBlockToUint(%s) = %d is wrong, want BitBlockToUint(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{150, 200, 128, 31, 33, 8, 0, 7, 15, 49, 63, 65} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToUint(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToUint() with a BitBlock of size != size in bits of uint (32 or 64 depending on architecture) should panic", bitBlock.Size())
				}
			}()
			BitBlockToUint(bitBlock)
		}()
	}

	// Test the functions Int8ToBitBlock and BitBlockToInt8.
	int8Numbers := []int8{48, 127, 0, 7, -128, -9, -64, -100}
	for _, x := range int8Numbers {
		bitBlock := Int8ToBitBlock(x)
		var x2 int8 = 0
		for i := 7; i >= 0; i-- {
			x2 *= 2;
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Int8ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToInt8(bitBlock)
		if x != x3 {
			t.Fatalf("the int8 obtained by calling BitBlockToInt8(%s) = %d is wrong, want BitBlockToInt8(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{0, 7, 9, 10, 15, 16} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToInt8(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToInt8() with a BitBlock of size != 8 should panic", bitBlock.Size())
				}
			}()
			BitBlockToInt8(bitBlock)
		}()
	}

	// Test the functions Int16ToBitBlock and BitBlockToInt16.
	int16Numbers := []int16{5156, 0, 218, 213, 84, -32768, 32767, 32000, -32500, 465, 9, 27, -13, -84, -654, -4613}
	for _, x := range int16Numbers {
		bitBlock := Int16ToBitBlock(x)
		var x2 int16 = 0
		for i := 15; i >= 0; i-- {
			x2 *= 2;
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Int16ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToInt16(bitBlock)
		if x != x3 {
			t.Fatalf("the int16 obtained by calling BitBlockToInt16(%s) = %d is wrong, want BitBlockToInt16(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{15, 17, 24, 19, 8, 7, 13} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToInt16(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToInt16() with a BitBlock of size != 16 should panic", bitBlock.Size())
				}
			}()
			BitBlockToInt16(bitBlock)
		}()
	}

	// Test the functions Int32ToBitBlock and BitBlockToInt32.
	int32Numbers := []int32{641564, 0, 94131, 2, 16489894, 94123134, 2147483647, -2147483648, -1616515612, 2107616513, 711235156, -334815, -1324}
	for _, x := range int32Numbers {
		bitBlock := Int32ToBitBlock(x)
		var x2 int32 = 0
		for i := 31; i >= 0; i-- {
			x2 *= 2;
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Int32ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToInt32(bitBlock)
		if x != x3 {
			t.Fatalf("the int32 obtained by calling BitBlockToInt32(%s) = %d is wrong, want BitBlockToInt32(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{30, 35, 40, 64, 0, 8, 1, 15} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToInt32(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToInt32() with a BitBlock of size != 32 should panic", bitBlock.Size())
				}
			}()
			BitBlockToInt32(bitBlock)
		}()
	}

	// Test the functions Int64ToBitBlock and BitBlockToInt64.
	int64Numbers := []int64{6515645615, 0, 8, 1, -7, -30, 9223372036854775807, -9223372036854775808, 651515615156156165, -4561636165116, 61561561}
	for _, x := range int64Numbers {
		bitBlock := Int64ToBitBlock(x)
		var x2 int64 = 0
		for i := 63; i >= 0; i-- {
			x2 *= 2;
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling Int64ToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToInt64(bitBlock)
		if x != x3 {
			t.Fatalf("the int64 obtained by calling BitBlockToInt64(%s) = %d is wrong, want BitBlockToInt64(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{150, 200, 128, 32, 8, 0, 7, 15, 49, 63, 65} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToInt64(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToInt64() with a BitBlock of size != 64 should panic", bitBlock.Size())
				}
			}()
			BitBlockToInt64(bitBlock)
		}()
	}

	// Test the functions IntToBitBlock and BitBlockToInt.
	var intNumbers []int
	if int(unsafe.Sizeof(int(0))) == 4 {
		intNumbers = []int{641564, 0, 94131, 2, 16489894, 94123134, 2147483647, -2147483648, -1616515612, 2107616513, 711235156, -334815, -1324}
	} else {
		intNumbers = []int{6515645615, 0, 8, 1, -7, -30, 9223372036854775807, -9223372036854775808, 651515615156156165, -4561636165116, 61561561}
	}
	for _, x := range intNumbers {
		bitBlock := IntToBitBlock(x)
		var x2 int = 0
		for i := int(unsafe.Sizeof(int(0))) * 8 - 1; i >= 0; i-- {
			x2 *= 2
			if bitBlock.Get(i) {
				x2 += 1
			}
		}
		if x != x2 {
			t.Fatalf("the BitBlock obtained by calling IntToBitBlock(%d) = %s is wrong", x, bitBlock.ToBinaryString())
		}
		x3 := BitBlockToInt(bitBlock)
		if x != x3 {
			t.Fatalf("the int obtained by calling BitBlockToInt(%s) = %d is wrong, want BitBlockToInt(%s) = %d", bitBlock.ToBinaryString(), x3, bitBlock.ToBinaryString(), x)
		}
	}
	for _, size := range []int{150, 200, 128, 31, 33, 8, 0, 7, 15, 49, 63, 65} {
		bitBlock := NewZeroBitBlock(size)
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to BitBlockToInt(bitBlock) with bitBlock.Size() = %d did not panic, calls to BitBlockToInt() with a BitBlock of size != size in bits of int (32 or 64 depending on architecture) should panic", bitBlock.Size())
				}
			}()
			BitBlockToInt(bitBlock)
		}()
	}
}

// Test the functions to set the first or last bits of an integer
// number to 1 and the rest to 0:
// - FirstBitsSet1Uint8, FirstBitsSet1Uint32, FirstBitsSet1Uint64.
// - LastBitsSet1Uint8, LastBitsSet1Uint32, LastBitsSet1Uint64.
func TestSetBitsInNumbers(t *testing.T) {
	type Test struct { id string; x uint; binstr string }

	// Test the FirstBitsSet1Uint8() function.
	for k := 0; k <= 8; k++ {
		x := FirstBitsSet1Uint8(k)
		for i := 0; i < k; i++ {
			if (x & (1 << i)) == 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint8(%d) = %d, got the %d-th least significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 8; i++ {
			if (x & (1 << i)) > 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint8(%d) = %d, got the %d-th least significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{10, 15, -9, -1, 9, 42, 32, 64, 19, -4} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to FirstBitsSet1Uint8(%d) did not panic", k)
				}
			}()
			FirstBitsSet1Uint8(k)
		}()
	}

	// Test the FirstBitsSet1Uint32() function.
	for k := 0; k <= 32; k++ {
		x := FirstBitsSet1Uint32(k)
		for i := 0; i < k; i++ {
			if (x & (1 << i)) == 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint32(%d) = %d, got the %d-th least significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 32; i++ {
			if (x & (1 << i)) > 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint32(%d) = %d, got the %d-th least significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{-1, -5, 33, 64, 65, 48, -32, -64, -10} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to FirstBitsSet1Uint32(%d) did not panic", k)
				}
			}()
			FirstBitsSet1Uint32(k)
		}()
	}

	// Test the FirstBitsSet1Uint64() function.
	for k := 0; k <= 64; k++ {
		x := FirstBitsSet1Uint64(k)
		for i := 0; i < k; i++ {
			if (x & (1 << i)) == 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint64(%d) = %d, got the %d-th least significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 64; i++ {
			if (x & (1 << i)) > 0 {
				t.Fatalf("wrong answer for FirstBitsSet1Uint64(%d) = %d, got the %d-th least significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{-3, -9, -8, -10, 65, 70, 100, 128, 256, 80, 100, -64, -63, -65} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to FirstBitsSet1Uint64(%d) did not panic", k)
				}
			}()
			FirstBitsSet1Uint64(k)
		}()
	}

	// Test the LastBitsSet1Uint8() function.
	for k := 0; k <= 8; k++ {
		x := LastBitsSet1Uint8(k)
		for i := 0; i < k; i++ {
			if (x & (1 << (7-i))) == 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint8(%d) = %d, got the %d-th most significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 8; i++ {
			if (x & (1 << (7-i))) > 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint8(%d) = %d, got the %d-th most significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{10, 15, -9, -1, 9, 42, 32, 64, 19, -4} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to LastBitsSet1Uint8(%d) did not panic", k)
				}
			}()
			LastBitsSet1Uint8(k)
		}()
	}

	// Test the LastBitsSet1Uint32() function.
	for k := 0; k <= 32; k++ {
		x := LastBitsSet1Uint32(k)
		for i := 0; i < k; i++ {
			if (x & (1 << (31-i))) == 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint32(%d) = %d, got the %d-th most significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 32; i++ {
			if (x & (1 << (31-i))) > 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint32(%d) = %d, got the %d-th most significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{-1, -5, 33, 64, 65, 48, -32, -64, -10} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to LastBitsSet1Uint32(%d) did not panic", k)
				}
			}()
			LastBitsSet1Uint32(k)
		}()
	}

	// Test the LastBitsSet1Uint64() function.
	for k := 0; k <= 64; k++ {
		x := LastBitsSet1Uint64(k)
		for i := 0; i < k; i++ {
			if (x & (1 << (63-i))) == 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint64(%d) = %d, got the %d-th most significant bit set to 0, want that bit to be set to 1", k, x, i)
			}
		}
		for i := k; i < 64; i++ {
			if (x & (1 << (63-i))) > 0 {
				t.Fatalf("wrong answer for LastBitsSet1Uint64(%d) = %d, got the %d-th most significant bit set to 1, want that bit to be set to 0", k, x, i)
			}
		}
	}
	for _, k := range []int{-3, -9, -8, -10, 65, 70, 100, 128, 256, 80, 100, -64, -63, -65} {
		func() {
			defer func() {
				panicMessage := recover()
				if panicMessage == nil {
					t.Fatalf("the call to LastBitsSet1Uint64(%d) did not panic", k)
				}
			}()
			LastBitsSet1Uint64(k)
		}()
	}
}

// The functions to get a panic message are executed.
// The message returned by those functions is not checked.
func TestPanicMessages(t *testing.T) {
	panicMessageNegativeSize(-5)
	panicMessageInvalidValueOutOfRange(7, 16, 2)
	panicMessageInvalidIndexOverBitBlock(10, 12)
	panicMessageInvalidRangeOverBitBlock(10, 8, 13)
	panicMessageInvalidRangeOverBitBlock(10, 13, 8)
	panicMessageInvalidNumberOfBitsToDiscardOverBitBlock(10, 30)
	panicMessageInvalidBitBlockSizeToConvertToInteger("int32", 64)
}