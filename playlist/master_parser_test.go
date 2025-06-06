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
	"strings"
	"testing"
)

func TestMasterPlayListParser(t *testing.T) {
	const s = `
#EXTM3U
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Main",DEFAULT=YES,URI="low/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Centerfield",URI="low/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Dugout",URI="low/dugout/audio-video.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=1280000,CODECS="mp4a.40.5",VIDEO="low"
low/main/audio-video.m3u8
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Main",DEFAULT=YES,URI="mid/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Centerfield",URI="mid/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Dugout",URI="mid/dugout/audio-video.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=2560000,CODECS="mp4a.40.5",VIDEO="mid"
mid/main/audio-video.m3u8
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Main",DEFAULT=YES,URI="hi/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Centerfield",URI="hi/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Dugout",URI="hi/dugout/audio-video.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=7680000,CODECS="mp4a.40.5",VIDEO="hi"
hi/main/audio-video.m3u8
`

	pl, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	} else if pl.Type() != PlayListTypeMaster {
		t.Fatalf("expect playlist type '%s', but got '%s'", PlayListTypeMaster, pl.Type())
	}

	master, ok := pl.(MasterPlayList)
	if !ok {
		t.Fatal("expect playlist to be MasterPlayList")
	}
	if master.MinVersion() != 1 {
		t.Errorf("expect min version %d, but got %d", 1, master.MinVersion())
	}

	for i, seg := range master.Segments {
		switch i {
		case 0:
			testMasterSegment(t, seg, MasterSegment{
				Stream: XStreamInf{Bandwidth: 1280000, Codecs: []string{"mp4a.40.5"}, Video: "low", URI: "low/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: XMediaTypeVideo, GroupId: "low", Name: "Main", Default: true, URI: "low/main/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "low", Name: "Centerfield", URI: "low/centerfield/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "low", Name: "Dugout", URI: "low/dugout/audio-video.m3u8"},
				},
			})

		case 1:
			testMasterSegment(t, seg, MasterSegment{
				Stream: XStreamInf{Bandwidth: 2560000, Codecs: []string{"mp4a.40.5"}, Video: "mid", URI: "mid/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: XMediaTypeVideo, GroupId: "mid", Name: "Main", Default: true, URI: "mid/main/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "mid", Name: "Centerfield", URI: "mid/centerfield/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "mid", Name: "Dugout", URI: "mid/dugout/audio-video.m3u8"},
				},
			})

		case 2:
			testMasterSegment(t, seg, MasterSegment{
				Stream: XStreamInf{Bandwidth: 7680000, Codecs: []string{"mp4a.40.5"}, Video: "hi", URI: "hi/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: XMediaTypeVideo, GroupId: "hi", Name: "Main", Default: true, URI: "hi/main/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "hi", Name: "Centerfield", URI: "hi/centerfield/audio-video.m3u8"},
					{Type: XMediaTypeVideo, GroupId: "hi", Name: "Dugout", URI: "hi/dugout/audio-video.m3u8"},
				},
			})

		default:
			t.Errorf("unexpected %dth master segment", i)
		}
	}
}

func testMasterSegment(t *testing.T, value, expect MasterSegment) {
	if !reflect.DeepEqual(value, expect) {
		t.Errorf("expect %+v, but got %+v", expect, value)
	}
}
