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
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func expectDuration(duration float64, expect string) bool {
	return strconv.FormatFloat(duration, 'f', -1, 64) == expect
}

func TestMediaPlayListParser(t *testing.T) {
	const s = `
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:2680
#EXT-X-DISCONTINUITY-SEQUENCE:1

#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=52"

#EXTINF:9.009,
http://media.example.com/first.ts
#EXTINF:9.009,
http://media.example.com/second.ts

#EXT-X-DISCONTINUITY
#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=53"

#EXTINF:3.003,
http://media.example.com/third.ts
#EXT-X-ENDLIST
`

	pl, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}

	if pl.Type() != PlayListTypeMedia {
		t.Fatalf("expect playlist type '%s', but got '%s'", PlayListTypeMedia, pl.Type())
	}

	media, ok := pl.(MediaPlayList)
	if !ok {
		t.Fatal("expect playlist to be MediaPlayList")
	}

	if media.TargetDuration != 10 {
		t.Errorf("expect target duration %d, but got %d", 10, media.TargetDuration)
	}
	if media.MinVersion() != 3 {
		t.Errorf("expect min version %d, but got %d", 3, media.MinVersion())
	}

	var (
		keys1 = []XKey{{Method: XKeyMethodAES128, URI: "https://priv.example.com/key.php?r=52"}}
		keys2 = []XKey{{Method: XKeyMethodAES128, URI: "https://priv.example.com/key.php?r=53"}}
	)

	if len(media.Segments) != 3 {
		t.Fatalf("expect %d media segments, but got %d", 3, len(media.Segments))
	}

	if total := media.TotalDuration(); !expectDuration(total, "21.021") {
		t.Errorf("expect total duration %s, but got %f", "21.021", total)
	}

	if media.MediaSequence != 2680 {
		t.Errorf("expect media sequence %d, but got %d", 2680, media.MediaSequence)
	}

	if media.DiscontinuitySequence != 1 {
		t.Errorf("expect discontinuity sequence %d, but got %d", 1, media.DiscontinuitySequence)
	}

	for i, seg := range media.Segments {
		switch {
		case i == 0 && seg.MediaSequence == uint64(2680+i) && seg.DiscontinuitySequence == 1 &&
			reflect.DeepEqual(seg.Keys, keys1) && expectDuration(seg.Duration, "9.009") &&
			seg.URI == "http://media.example.com/first.ts":
		case i == 1 && seg.MediaSequence == uint64(2680+i) && seg.DiscontinuitySequence == 1 &&
			reflect.DeepEqual(seg.Keys, keys1) && expectDuration(seg.Duration, "9.009") &&
			seg.URI == "http://media.example.com/second.ts":
		case i == 2 && seg.MediaSequence == uint64(2680+i) && seg.DiscontinuitySequence == 2 &&
			reflect.DeepEqual(seg.Keys, keys2) && expectDuration(seg.Duration, "3.003") &&
			seg.URI == "http://media.example.com/third.ts":
		default:
			t.Errorf("unexpect media segment: index=%d, segment=%+v", i, seg)
		}
	}
}
