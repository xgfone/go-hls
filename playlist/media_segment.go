package playlist

import (
	"encoding/binary"
	"time"

	"github.com/xgfone/go-hls/aes128"
)

// MediaSegment represents a media segment in a media playlist.
//
// See [[RFC 8216, 4.3.2]].
//
// [RFC 8216, 4.3.2]: https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2
type MediaSegment struct {
	URI   string // Required.
	Title string

	Duration  float64 // Required. Unit: Second
	ByteRange XByteRange
	Key       XKey
	Map       XMap

	ProgramDateTime time.Time

	MediaSequence         uint64 // Cannot be encoded
	DiscontinuitySequence uint64 // Cannot be encoded

	Discontinuity bool
}

// IV try to decode the iv from a hexadecimal-sequence string to a 16-octet bytes.
func (s MediaSegment) IV() (data []byte, err error) {
	if s.Key.IV != "" {
		var seq _HexSequence
		err = seq.decode(s.Key.IV)
		data = []byte(seq)
	} else {
		data = make([]byte, 16)
		binary.BigEndian.PutUint64(data[8:], s.MediaSequence)
	}
	return
}

// AES128Decrypt decrypted the encrypted data with key and iv.
//
// If padding is true, remove the PKCS7 padding. Or, do nothing.
func (s MediaSegment) AES128Decrypt(encryptedData, key []byte, removePadding bool) (decryptedData []byte, err error) {
	iv, _ := s.IV()
	return aes128.Decrypt(encryptedData, key, iv, removePadding)
}
