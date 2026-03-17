package xlog

import "unicode/utf16"

func encodeUTF16LE(s string) []byte {
	encoded := utf16.Encode([]rune(s))
	buf := make([]byte, len(encoded)*2)
	for i, v := range encoded {
		buf[i*2] = byte(v)
		buf[i*2+1] = byte(v >> 8)
	}
	return buf
}
