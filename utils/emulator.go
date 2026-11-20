package utils

func DecodeInstruction(i1, i2 byte) (uint8, uint8, uint8, uint8, uint8, uint16) {
	O := i1 >> 4 // Recorrer 4 bits para mantener los 4 bits altos
	X := i1 & 0b00001111
	Y := i2 >> 4
	N := i2 & 0b00001111
	NN := i2
	NNN := uint16(X) << 8
	NNN = NNN | uint16(i2)

	return O, X, Y, N, NN, NNN
}
