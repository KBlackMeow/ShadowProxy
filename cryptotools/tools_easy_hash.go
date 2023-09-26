package cryptotools

func EasyHash_uint64(data string) uint64 {
	code := uint64(0)
	for i, v := range data {
		code += uint64(i) * uint64(v) % 65537
	}
	return code
}
