package encoder

import "testing"

func TestTextEncoderString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := encT.String([]byte{}, tt.in)
		if got, want := string(b), tt.in; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.in, got, want)
		}
		c := encT.String([]byte{}, tt.out)
		if got, want := string(c), tt.out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.out, got, want)
		}
	}
}

func BenchmarkTextEncoderString(b *testing.B) {
	tests := map[string]string{
		"NoEncoding":       `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingFirst":    `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingMiddle":   `aaaaaaaaaaaaaaaaaaaaaaaaa"aaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingLast":     `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`,
		"MultiBytesFirst":  `❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"MultiBytesMiddle": `aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa`,
		"MultiBytesLast":   `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️`,
	}
	for name, str := range tests {
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = encT.String(buf, str)
			}
		})
	}
}
