package server

import (
	"bytes"
	"testing"
)

// TestAppendFrame checks the frame header of each length class (RFC 6455,
// section 5.2): a short text frame, a 16-bit and a 64-bit length.
func TestAppendFrame(t *testing.T) {
	cases := []struct {
		n      int
		header []byte
	}{
		{0, []byte{0x81, 0}},
		{125, []byte{0x81, 125}},
		{126, []byte{0x81, 126, 0, 126}},
		{65535, []byte{0x81, 126, 0xff, 0xff}},
		{65536, []byte{0x81, 127, 0, 0, 0, 0, 0, 1, 0, 0}},
	}
	for _, c := range cases {
		msg := bytes.Repeat([]byte{'x'}, c.n)
		got := appendFrame(nil, msg)
		want := append(append([]byte{}, c.header...), msg...)
		if !bytes.Equal(got, want) {
			t.Errorf("length %d: got % x, want the header % x and the payload", c.n, got[:min(len(got), len(c.header))], c.header)
		}
	}
}
