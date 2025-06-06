package playlist

import "testing"

func TestFormatIV(t *testing.T) {
	iv := FormatIV([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, true)
	expect := "0x0102030405060708090A0B0C0D0E0F10"
	if expect != iv {
		t.Errorf("expect IV '%s', but got '%s'", expect, iv)
	}
}
