package photon

// alias for int16 (2^16)
type Short int16

func DeserializeInt(src []byte, offset *int) int {
	val := int(src[*offset]) << 24
	*offset++
	val = val | int(src[*offset])<<16
	*offset++
	val = val | int(src[*offset])<<8
	*offset++
	res := val | int(src[*offset])
	*offset++

	return res
}

func DeserializeShort(src []byte, offset *int) Short {
	val := Short(src[*offset]) << 8
	*offset++
	res := val | Short(src[*offset])
	*offset++

	return res
}
