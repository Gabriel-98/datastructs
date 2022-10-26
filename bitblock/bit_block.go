// LICENCE NOT YET DEFINED.

package bitblock


import (
	"strconv"
	"unsafe"
)


// panicMessageNegativeSize returns the message that should
// appear within a panic, which will be raised because an
// attempt was made to create a structure with a negative
// size.
//
// The message will indicate the value that was passed to be
// the size of some structure.
func panicMessageNegativeSize(size int) string {
	return "size (" + strconv.Itoa(size) + ") cannot be negative"
}

// panicMessageInvalidValueOutOfRange returns the message that
// should appear within a panic, which will be raised because
// some function or method was passed a value that is not within
// the range of valid values.
//
// The message will indicate the value passed to the function and
// the limits of the range of valid values.
func panicMessageInvalidValueOutOfRange(minValue int, maxValue int, value int) string {
	return "invalid value (" + strconv.Itoa(value) + "), only values between " + strconv.Itoa(minValue) + " and " + strconv.Itoa(maxValue) + " (both inclusive) are allowed"
}

// panicMessageInvalidIndexOverBitBlock returns the message
// that will appear within a panic that will be raised because
// an invalid index was passed to a method from BitBlock.
//
// The message will indicate the size of the BitBlock and the
// position that was attempted to be accessed or modified.
func panicMessageInvalidIndexOverBitBlock(size int, pos int) string {
	return "invalid index [" + strconv.Itoa(pos) + "] for BitBlock with size " + strconv.Itoa(size)
}

// panicMessageInvalidRangeOverBitBlock returns the message
// that will appear within a panic that will be raised because
// an invalid range was passed to a method from BitBlock.
//
// The message will indicate the limits of the range and if
// l <= r also the size of the BitBlock.
func panicMessageInvalidRangeOverBitBlock(size int, l int, r int) string {
	message := "invalid range [" + strconv.Itoa(l) + ":" + strconv.Itoa(r) + "] for BitBlock"
	if l > r {
		message += ", start of range [" + strconv.Itoa(l) + "] is greater than end of range [" + strconv.Itoa(r) + "]"
	} else {
		message += " with size " + strconv.Itoa(size)
	}
	return message
}

// panicMessageInvalidNumberOfBitsToDiscardOverBitBlock returns
// the message that should appear within a panic, which will be
// raised beacuse an invalid number of bits to discard within a
// BitBlock was passed to some function or method.
//
// The message will indicate the size of the BitBlock and the
// number of bits that were attempted to be discarded.
func panicMessageInvalidNumberOfBitsToDiscardOverBitBlock(bitBlockSize int, bitsToDiscard int) string {
	return "invalid number of bits to discard (" + strconv.Itoa(bitsToDiscard)  + ") in a BitBlock, this one must be non-negative and less than or equal to the size of the BitBlock (" + strconv.Itoa(bitBlockSize) + ")"
}

// panicMessageInvalidBitBlockSizeToConvertToInteger returns
// the message that should appear within a panic, which will
// be raised because an attempt was made to convert a BitBlock
// to a specific integer type, but the size of the BitBlock was
// different than the number of bits in the target data type.
//
// The message will indicate the name of the target data type and
// the size of the BitBlock that was attempted to be converted.
func panicMessageInvalidBitBlockSizeToConvertToInteger(typeName string, bitBlockSize int) string {
	return "invalid BitBlock size, BitBlock with size " + strconv.Itoa(bitBlockSize) + "cannot be converted to " + typeName
}

// FirstBitsSet1Uint8 returns an 8-bit unsigned integer
// (uint8) in which only the k least significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 8.
func FirstBitsSet1Uint8(k int) uint8 {
	switch true {
		case !(0 <= k && k <= 8):
			panic(panicMessageInvalidValueOutOfRange(0, 8, k))
		case k == 8:
			return 0xFF
		default:
			return (1 << k) - 1
	}
}

// FirstBitsSet1Uint32 returns a 32-bit unsigned integer
// (uint32) in which only the k least significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 32.
func FirstBitsSet1Uint32(k int) uint32 {
	switch true {
		case !(0 <= k && k <= 32):
			panic(panicMessageInvalidValueOutOfRange(0, 32, k))
		case k == 32:
			return 0xFFFFFFFF
		default:
			return (1 << k) - 1
	}
}

// FirstBitsSet1Uint64 returns a 64-bit unsigned integer
// (uint64) in which only the k least significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 64.
func FirstBitsSet1Uint64(k int) uint64 {
	switch true {
		case !(0 <= k && k <= 64):
			panic(panicMessageInvalidValueOutOfRange(0, 64, k))
		case k == 64:
			return 0xFFFFFFFFFFFFFFFF
		default:
			return (1 << k) - 1
	}
}

// LastBitsSet1Uint8 returns an 8-bit unsigned integer
// (uint8) in which only the k most significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 8.
func LastBitsSet1Uint8(k int) uint8 {
	switch true {
		case !(0 <= k && k <= 8):
			panic(panicMessageInvalidValueOutOfRange(0, 8, k))
		default:
			return 0xFF ^ FirstBitsSet1Uint8(8 - k)
	}
}

// LastBitsSet1Uint32 returns a 32-bit unsigned integer
// (uint32) in which only the k most significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 32.
func LastBitsSet1Uint32(k int) uint32 {
	switch true {
		case !(0 <= k && k <= 32):
			panic(panicMessageInvalidValueOutOfRange(0, 32, k))
		default:
			return 0xFFFFFFFF ^ FirstBitsSet1Uint32(32 - k)
	}
}

// LastBitsSet1Uint64 returns a 64-bit unsigned integer
// (uint64) in which only the k most significant bits
// are set to 1, the rest are set to 0. This function
// panics if k < 0 or k > 64.
func LastBitsSet1Uint64(k int) uint64 {
	switch true {
		case !(0 <= k && k <= 64):
			panic(panicMessageInvalidValueOutOfRange(0, 64, k))
		default:
			return 0xFFFFFFFFFFFFFFFF ^ FirstBitsSet1Uint64(64 - k)
	}
}

// A BitBlock represents a sequence of bits, which allows
// each bit to be read and modified individually.
//
// Each bit can be getted and setted as a bool value,
// but internally each bit occupies only 1 bit, unlike
// a bool type which occupies 1 byte.
//
// A BitBlock is a 0-based indexed structure and uses a
// slice of bytes as the underlying structure.
type BitBlock struct {
	bits []byte
	size int
}

// NewZeroBitBlock returns a new BitBlock with all bits
// set to 0. NewZeroBitBlock panics if size < 0.
func NewZeroBitBlock(size int) *BitBlock {
	if size < 0 {
		panic(panicMessageNegativeSize(size))
	}
	bits := make([]byte, (size + 7) / 8)
	return &BitBlock{
		bits: bits,
		size: size,
	}
}

// BytesToBitBlock returns a new BitBlock, which will contain a
// copy of the first size bits of src. If src does not have
// enough bits to fully set the required number of bits, the
// remaining bits will be set to 0. BytesToBitBlock panics if
// size < 0.
func BytesToBitBlock(src []byte, size int) *BitBlock {
	if size < 0 {
		panic(panicMessageNegativeSize(size))
	}
	bits := make([]byte, (size + 7) / 8)
	for i := 0; i < len(bits)-1 && i < len(src); i++ {
		bits[i] = src[i]
	}
	if len(src) >= len(bits) && len(bits) >= 1 {
		if (size & 7) == 0 {
			bits[len(bits) - 1] = src[len(bits) - 1]
		} else {
			bits[len(bits) - 1] = src[len(bits) - 1] & FirstBitsSet1Uint8(size & 7)
		}
	}
	return &BitBlock{
		bits: bits,
		size: size,
	}
}

// Get returns the value of the bit at position pos.
// If pos < 0 or pos >= block.Size(), Get panics.
func (block *BitBlock) Get(pos int) bool {
	if !(0 <= pos && pos < block.size) {
		panic(panicMessageInvalidIndexOverBitBlock(block.size, pos))
	}
	if (block.bits[pos >> 3] & (1 << (pos & 7))) > 0 {
		return true
	}
	return false
}

// Set0 sets the bit at position pos to 0.
// If pos < 0 or pos >= block.Size(), Set0 panics.
func (block *BitBlock) Set0(pos int) {
	if !(0 <= pos && pos < block.size) {
		panic(panicMessageInvalidIndexOverBitBlock(block.size, pos))
	}
	block.bits[pos >> 3] &= (0xFF ^ (1 << (pos & 7)))
}

// Set1 sets the bit at position pos to 1.
// If pos < 0 or pos >= block.Size(), Set1 panics.
func (block *BitBlock) Set1(pos int) {
	if !(0 <= pos && pos < block.size) {
		panic(panicMessageInvalidIndexOverBitBlock(block.size, pos))
	}
	block.bits[pos >> 3] |= (1 << (pos & 7))
}

// Set sets the bit at position pos to 1 or 0 depending
// on whether value == true or value == false respectively.
// If pos < 0 or pos >= block.Size(), Set panics.
func (block *BitBlock) Set(pos int, value bool) {
	if !(0 <= pos && pos < block.size) {
		panic(panicMessageInvalidIndexOverBitBlock(block.size, pos))
	}
	if value {
		block.Set1(pos)
	} else {
		block.Set0(pos)
	}
}

// Size returns the number of bits used by the BitBlock.
func (block *BitBlock) Size() int {
	return block.size
}

// GetSubBlock returns a new BitBlock containing a copy of
// the bits from position l to position r (including l, but
// excluding r). This method panics if l and r form an
// invalid range for this BitBlock.
func (block *BitBlock) GetSubBlock(l int, r int) *BitBlock {
	if !(0 <= l && l <= r && r <= block.size) {
		panic(panicMessageInvalidRangeOverBitBlock(block.size, l, r))
	}
	size := r-l
	bitBlock := NewZeroBitBlock(size)
	for pos:=0; pos < size; pos++ {
		bitBlock.Set(pos, block.Get(l + pos))
	}
	return bitBlock
}

// ToBytes returns a copy of the bits in this BitBlock as a
// slice of bytes.
// The size of the returned slice is the minimum necessary
// to contain at least block.Size() bits. The padding bits
// will be equal to 0.
func (block *BitBlock) ToBytes() []byte {
	bits := make([]byte, len(block.bits))
	copy(bits, block.bits)
	return bits
}

// Clone returns a new BitBlock containing a copy of the
// bits in this BitBlock.
func (block *BitBlock) Clone() *BitBlock {
	return &BitBlock{
		bits: block.ToBytes(),
		size: block.Size(),
	}
}

// RemoveFirstBits returns a new BitBlock containing a copy of
// the bits in this BitBlock, but without copying the first k bits.
// This method panics if k < 0 or k > block.Size().
func (block *BitBlock) RemoveFirstBits(k int) *BitBlock {
	if !(0 <= k && k <= block.size) {
		panic(panicMessageInvalidNumberOfBitsToDiscardOverBitBlock(block.size, k))
	}
	size := block.size - k
	bits := make([]byte, (size + 7) / 8)
	mask1 := LastBitsSet1Uint8(8 - (k & 7))
	mask2 := 0xFF ^ mask1
	for i, j := 0, (k / 8); i < len(bits); i,j = i+1, j+1 {
		bits[i] = (block.bits[j] & mask1) >> (k & 7)
		if j+1 < len(block.bits) {
			bits[i] |= (block.bits[j+1] & mask2) << (8 - (k & 7))
		}
	}
	return &BitBlock{
		bits: bits,
		size: size,
	}
}

// RemoveLastBits returns a new BitBlock containing a copy of
// the bits in this BitBlock, but without copying the last k bits.
// This method panics if k < 0 or k > block.Size().
func (block *BitBlock) RemoveLastBits(k int) *BitBlock {
	if !(0 <= k && k <= block.size) {
		panic(panicMessageInvalidNumberOfBitsToDiscardOverBitBlock(block.size, k))
	}
	return BytesToBitBlock(block.bits, block.size - k)
}

// ToBinaryString returns this BitBlock as a binary string.
func (block *BitBlock) ToBinaryString() string {
	binChars := make([]byte, block.size)
	for i := 0; i < block.size; i++ {
		if block.Get(i) {
			binChars[i] = '1'
		} else {
			binChars[i] = '0'
		}
	}
	return string(binChars)
}

// Concatenate receives multiple BitBlocks and returns a new
// BitBlock containing the bits from the other BitBlocks in
// the same order as they were passed to this method.
func Concatenate(bitBlocks ...*BitBlock) *BitBlock {
	size := 0
	for _, bitBlock := range bitBlocks {
		size += bitBlock.Size()
	}
	concatenatedBitBlock := NewZeroBitBlock(size)
	currentSize := 0
	for _, bitBlock := range bitBlocks {
		for i := 0; i < bitBlock.Size(); i++ {
			concatenatedBitBlock.Set(currentSize, bitBlock.Get(i))
			currentSize++
		}
	}
	return concatenatedBitBlock
}

// IntToBitBlock converts an integer to a BitBlock.
// The returned BitBlock will be either 32 or 64 bits depending
// on the type of architecture. If the architecture is 32 bits,
// then it returns a 32-bit BitBlock; otherwise it returns a
// 64-bit BitBlock.
// The number is stored in little endian format. 
func IntToBitBlock(x int) *BitBlock {
	switch unsafe.Sizeof(x) {
		case 4:
			return Int32ToBitBlock(int32(x))
		default:
			return Int64ToBitBlock(int64(x))
	}
}

// Int8ToBitBlock converts an 8-bit integer to an 8-bit BitBlock.
// The number is stored in little endian format.
func Int8ToBitBlock(x int8) *BitBlock {
	return Uint8ToBitBlock(uint8(x))
}

// Int16ToBitBlock converts a 16-bit integer to a 16-bit BitBlock.
// The number is stored in little endian format.
func Int16ToBitBlock(x int16) *BitBlock {
	return Uint16ToBitBlock(uint16(x))
}

// Int32ToBitBlock converts a 32-bit integer to a 32-bit BitBlock.
// The number is stored in little endian format.
func Int32ToBitBlock(x int32) *BitBlock {
	return Uint32ToBitBlock(uint32(x))
}

// Int64ToBitBlock converts a 64-bit integer to a 64-bit BitBlock.
// The number is stored in little endian format.
func Int64ToBitBlock(x int64) *BitBlock {
	return Uint64ToBitBlock(uint64(x))
}

// UintToBitBlock converts an unsigned integer to a BitBlock.
// The returned BitBlock will be either 32 or 64 bits depending
// on the type of architecture. If the architecture is 32 bits,
// then it returns a 32-bit BitBlock; otherwise it returns a
// 64-bit BitBlock.
// The number is stored in little endian format. 
func UintToBitBlock(x uint) *BitBlock {
	switch unsafe.Sizeof(x) {
		case 4:
			return Uint32ToBitBlock(uint32(x))
		default:
			return Uint64ToBitBlock(uint64(x))
	}
}

// Uint8ToBitBlock converts an 8-bit unsigned integer to an 8-bit BitBlock.
// The number is stored in little endian format.
func Uint8ToBitBlock(x uint8) *BitBlock {
	bits := []byte{x}
	return BytesToBitBlock(bits, 8)
}

// Uint16ToBitBlock converts a 16-bit unsigned integer to a 16-bit BitBlock.
// The number is stored in little endian format.
func Uint16ToBitBlock(x uint16) *BitBlock {
	bits := make([]byte, 2)
	for i := 0; i < len(bits); i++ {
		bits[i] = byte((x & (uint16(0xFF) << (8 * i))) >> (8 * i))
	}
	return BytesToBitBlock(bits, 16)
}

// Uint32ToBitBlock converts a 32-bit unsigned integer to a 32-bit BitBlock.
// The number is stored in little endian format.
func Uint32ToBitBlock(x uint32) *BitBlock {
	bits := make([]byte, 4)
	for i := 0; i < len(bits); i++ {
		bits[i] = byte((x & (uint32(0xFF) << (8 * i))) >> (8 * i))
	}
	return BytesToBitBlock(bits, 32)
}

// Uint64ToBitBlock converts a 64-bit unsigned integer to a 64-bit BitBlock.
// The number is stored in little endian format.
func Uint64ToBitBlock(x uint64) *BitBlock {
	bits := make([]byte, 8)
	for i := 0; i < len(bits); i++ {
		bits[i] = byte((x & (uint64(0xFF) << (8 * i))) >> (8 * i))
	}
	return BytesToBitBlock(bits, 64)
}

// BitBlockToInt converts a BitBlock to an integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from the size of an int variable (32 or 64 bits
// depending on the architecture).
func BitBlockToInt(bitBlock *BitBlock) int {
	intSize := int(unsafe.Sizeof(int(0))) * 8
	if bitBlock.Size() != intSize {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("int", bitBlock.Size()))
	}
	return int(BitBlockToUint(bitBlock))
}

// BitBlockToInt8 converts an 8-bit BitBlock to an 8-bit integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 8.
func BitBlockToInt8(bitBlock *BitBlock) int8 {
	if bitBlock.Size() != 8 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("int8", bitBlock.Size()))
	}
	return int8(BitBlockToUint8(bitBlock))
}

// BitBlockToInt16 converts a 16-bit BitBlock to a 16-bit integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 16.
func BitBlockToInt16(bitBlock *BitBlock) int16 {
	if bitBlock.Size() != 16 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("int16", bitBlock.Size()))
	}
	return int16(BitBlockToUint16(bitBlock))
}

// BitBlockToInt32 converts a 32-bit BitBlock to a 32-bit integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 32.
func BitBlockToInt32(bitBlock *BitBlock) int32 {
	if bitBlock.Size() != 32 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("int32", bitBlock.Size()))
	}
	return int32(BitBlockToUint32(bitBlock))
}

// BitBlockToInt64 converts a 64-bit BitBlock to a 64-bit integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 64.
func BitBlockToInt64(bitBlock *BitBlock) int64 {
	if bitBlock.Size() != 64 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("int64", bitBlock.Size()))
	}
	return int64(BitBlockToUint64(bitBlock))
}

// BitBlockToUint converts a BitBlock to an unsigned integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from the size of an uint variable (32 or 64 bits
// depending on the architecture).
func BitBlockToUint(bitBlock *BitBlock) uint {
	uintSize := int(unsafe.Sizeof(uint(0))) * 8
	if bitBlock.Size() != uintSize {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("uint", bitBlock.Size()))
	}
	var x uint = 0
	bytes := bitBlock.ToBytes()
	for i := 0; i < len(bytes); i++ {
		x = x | (uint(bytes[i]) << (8 * i))
	}
	return x
}

// BitBlockToUint8 converts an 8-bit BitBlock to an 8-bit unsigned integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 8.
func BitBlockToUint8(bitBlock *BitBlock) uint8 {
	if bitBlock.Size() != 8 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("uint8", bitBlock.Size()))
	}
	bytes := bitBlock.ToBytes()
	var x uint8 = bytes[0]
	return x
}

// BitBlockToUint16 converts a 16-bit BitBlock to a 16-bit unsigned integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 16.
func BitBlockToUint16(bitBlock *BitBlock) uint16 {
	if bitBlock.Size() != 16 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("uint16", bitBlock.Size()))
	}
	var x uint16 = 0
	bytes := bitBlock.ToBytes()
	for i := 0; i < len(bytes); i++ {
		x = x | (uint16(bytes[i]) << (8 * i))
	}
	return x
}

// BitBlockToUint32 converts a 32-bit BitBlock to a 32-bit unsigned integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 32.
func BitBlockToUint32(bitBlock *BitBlock) uint32 {
	if bitBlock.Size() != 32 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("uint32", bitBlock.Size()))
	}
	var x uint32 = 0
	bytes := bitBlock.ToBytes()
	for i := 0; i < len(bytes); i++ {
		x = x | (uint32(bytes[i]) << (8 * i))
	}
	return x
}

// BitBlockToUint64 converts a 64-bit BitBlock to a 64-bit unsigned integer.
// The BitBlock is supposed to be in little endian format.
// This function panics if the size of the passed BitBlock is
// different from 64.
func BitBlockToUint64(bitBlock *BitBlock) uint64 {
	if bitBlock.Size() != 64 {
		panic(panicMessageInvalidBitBlockSizeToConvertToInteger("uint64", bitBlock.Size()))
	}
	var x uint64 = 0
	bytes := bitBlock.ToBytes()
	for i := 0; i < len(bytes); i++ {
		x = x | (uint64(bytes[i]) << (8 * i))
	}
	return x
}