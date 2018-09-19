package btmgmt

func copyReverse(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	for l, r := 0, len(dst)-1; l < r; l, r = l+1, r-1 {
		dst[l], dst[r] = dst[r], dst[l]
	}
	return dst
}

func testBit(in byte, n uint8) bool {
	return in&(1<<n) > 0
}

