// Copyright 2025 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Keys      []XKey
	Map       XMap

	ProgramDateTime time.Time

	MediaSequence         uint64 // Cannot be encoded
	DiscontinuitySequence uint64 // Cannot be encoded

	Discontinuity bool
}

// IV try to decode the iv from a hexadecimal-sequence string to a 16-octet bytes.
func (s MediaSegment) IV() (data []byte, err error) {
	if len(s.Keys) > 0 && s.Keys[0].IV != "" {
		var seq _HexSequence
		err = seq.decode(s.Keys[0].IV)
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

// GetSegmentIndexByMediaSequence returns the index of the media segment
// whose the media sequence is equal to the given seq.
//
// Return -1 if not found.
//
// NOTE: It is only valid for the same media playlist.
func (pl MediaPlayList) GetSegmentIndexByMediaSequence(seq uint64) (index int) {
	index = int(seq - pl.MediaSequence)
	if index < 0 || index >= len(pl.Segments) {
		return -1
	}

	switch seg := &pl.Segments[index]; {
	case seq < seg.MediaSequence:
		segments := pl.Segments[:index+1]
		for i := len(segments) - 1; i >= 0; i-- {
			if segments[i].MediaSequence == seq {
				index = i
				break
			}
		}

	case seg.MediaSequence < seq:
		segments := pl.Segments[index+1:]
		for i := range segments {
			if segments[i].MediaSequence == seq {
				index += 1 + i
				break
			}
		}
	}

	return
}
