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
	"bytes"
	"testing"
)

func TestMediaPlayListEncoderSimple(t *testing.T) {
	const expect = `
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXTINF:9.009,
http://media.example.com/first.ts
#EXTINF:9.009,
http://media.example.com/second.ts
#EXTINF:3.003,
http://media.example.com/third.ts
#EXT-X-ENDLIST
`

	pl := MediaPlayList{
		Version: 3,

		EndList:        true,
		TargetDuration: 10,

		Segments: []MediaSegment{
			{
				URI:      "http://media.example.com/first.ts",
				Duration: 9.009,
			},
			{
				URI:      "http://media.example.com/second.ts",
				Duration: 9.009,
			},
			{
				URI:      "http://media.example.com/third.ts",
				Duration: 3.003,
			},
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	if err := pl.Output(buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}

func TestMediaPlayListEncoderMore(t *testing.T) {
	const expect = `
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:15
#EXT-X-MEDIA-SEQUENCE:7794
#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=52"
#EXTINF:2.833,
http://media.example.com/fileSequence52-A.ts
#EXTINF:15,
http://media.example.com/fileSequence52-B.ts
#EXTINF:13.333,
http://media.example.com/fileSequence52-C.ts
#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=53"
#EXTINF:15,
http://media.example.com/fileSequence53-A.ts
`

	var (
		keys1 = []XKey{{Method: XKeyMethodAES128, URI: "https://priv.example.com/key.php?r=52"}}
		keys2 = []XKey{{Method: XKeyMethodAES128, URI: "https://priv.example.com/key.php?r=53"}}
	)

	pl := MediaPlayList{
		Version: 3,

		TargetDuration: 15,
		MediaSequence:  7794,

		Segments: []MediaSegment{
			{
				Keys:     keys1,
				URI:      "http://media.example.com/fileSequence52-A.ts",
				Duration: 2.833,
			},
			{
				Keys:     keys1,
				URI:      "http://media.example.com/fileSequence52-B.ts",
				Duration: 15,
			},
			{
				Keys:     keys1,
				URI:      "http://media.example.com/fileSequence52-C.ts",
				Duration: 13.333,
			},
			{
				Keys:     keys2,
				URI:      "http://media.example.com/fileSequence53-A.ts",
				Duration: 15,
			},
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	if err := pl.Output(buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}
