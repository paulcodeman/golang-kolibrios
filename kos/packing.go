package kos

func littleEndianUint16(buffer []byte, offset int) uint16 {
	return uint16(buffer[offset]) | uint16(buffer[offset+1])<<8
}

func littleEndianUint32(buffer []byte, offset int) uint32 {
	return uint32(buffer[offset]) |
		uint32(buffer[offset+1])<<8 |
		uint32(buffer[offset+2])<<16 |
		uint32(buffer[offset+3])<<24
}

func unpackUnsignedPoint(packed uint32) Point {
	return Point{
		X: int(packed >> 16),
		Y: int(packed & 0xFFFF),
	}
}

func unpackSignedPackedPoint(packed uint32) Point {
	y := int(int16(uint16(packed)))
	x := int(int32(packed) >> 16)
	if y < 0 {
		x++
	}

	return Point{
		X: x,
		Y: y,
	}
}

func trimASCIIField(field []byte) string {
	end := len(field)

	for end > 0 {
		last := field[end-1]
		if last != 0 && last != ' ' {
			break
		}
		end--
	}

	return string(field[:end])
}
