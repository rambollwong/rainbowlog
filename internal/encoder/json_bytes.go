package encoder

func (j JsonEncoder) Bytes(dst []byte, s []byte) []byte {
	dst = j.DoubleQuote(dst)
	for i := 0; i < len(s); i++ {
		if !noEscapeTable[s[i]] {
			dst = appendBytesComplex(dst, s, i)
			dst = j.DoubleQuote(dst)
			return dst
		}
	}
	dst = append(dst, s...)
	return j.DoubleQuote(dst)
}

// Hex encodes the input bytes to a hex string and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte and encodes it as hex using
// the hex lookup table.
func (j JsonEncoder) Hex(dst []byte, s []byte) []byte {
	dst = append(dst, '"')
	for _, v := range s {
		dst = append(dst, hex[v>>4], hex[v&0x0f])
	}
	return append(dst, '"')
}
