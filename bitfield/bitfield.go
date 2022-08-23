package bitfield

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	byteOffset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}

	return bf[byteIndex]>>uint(7-byteOffset)&1 == 1
}

func (bf BitField) SetPiece(index int) {
	byteIndex := index / 8
	byteOffset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}

	bf[byteIndex] |= 1 << uint(7-byteOffset)
}
