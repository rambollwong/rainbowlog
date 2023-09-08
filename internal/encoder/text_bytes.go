package encoder

func (j TextEncoder) Bytes(dst []byte, s []byte) []byte {
	return append(dst, s...)
}

// Hex encodes the input bytes to a hex string and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte and encodes it as hex using
// the hex lookup table.
func (j TextEncoder) Hex(dst []byte, s []byte) []byte {
	for _, v := range s {
		dst = append(dst, hex[v>>4], hex[v&0x0f])
	}
	return dst
}
