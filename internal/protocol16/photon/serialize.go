package photon

func Serialize(val int, offset *int, target []byte) {
	target[*offset] = byte(val >> 24)
	*offset++
	target[*offset] = byte(val >> 16)
	*offset++
	target[*offset] = byte(val >> 8)
	*offset++
	target[*offset] = byte(val)
	*offset++
}
