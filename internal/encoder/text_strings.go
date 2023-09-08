package encoder

func (j TextEncoder) DoubleQuote(dst []byte) []byte {
	return append(dst, '"')
}

func (j TextEncoder) String(dst []byte, s string) []byte {
	return append(dst, s...)
}

func (j TextEncoder) Strings(dst []byte, s ...string) []byte {
	if len(s) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.String(dst, s[0])
	if len(s) > 1 {
		for _, v := range s[1:] {
			dst = j.Comma(dst)
			dst = j.String(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
