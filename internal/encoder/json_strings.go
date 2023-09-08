package encoder

func (j JsonEncoder) DoubleQuote(dst []byte) []byte {
	return append(dst, '"')
}

// String encodes the input string to json and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte in the string looking
// for characters that need json or utf8 encoding. If the string
// does not need encoding, then the string is appended in its
// entirety to the byte slice.
// If we encounter a byte that does need encoding, switch up
// the operation and perform a byte-by-byte read-encode-append.
func (j JsonEncoder) String(dst []byte, s string) []byte {
	// Start with a double quote.
	dst = j.DoubleQuote(dst)
	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			dst = appendStringComplex(dst, s, i)
			dst = j.DoubleQuote(dst)
			return dst
		}
	}
	// The string has no need for encoding and therefore is directly
	// appended to the byte slice.
	dst = append(dst, s...)
	// End with a double quote
	return j.DoubleQuote(dst)
}

func (j JsonEncoder) Strings(dst []byte, s ...string) []byte {
	if len(s) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.String(dst, s[0])
	if len(s) > 1 {
		for _, v := range s[1:] {
			dst = j.Delim(dst)
			dst = j.String(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
