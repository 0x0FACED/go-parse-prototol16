package photon

import (
	"math"
)

// This method calculate CRC
func CalculateCRC(byteArr []byte, len int) uint {
	var res uint
	var key uint
	res = math.MaxUint
	key = 3988292384

	for i := 0; i < len; i++ {
		res ^= uint(byteArr[i])

		for j := 0; j < 8; j++ {
			if (res & 1) > 0 {
				res = res>>1 ^ key
			} else {
				res >>= 1
			}
		}
	}

	return res
}
