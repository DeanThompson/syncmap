package syncmap

const hashSeed uint32 = 131

func BkdrHash(str string) uint32 {
	var result uint32

	for _, c := range str {
		result = result*hashSeed + uint32(c)
	}

	return result
}
