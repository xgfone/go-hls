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
	"strings"
	"testing"
)

func TestMasterPlayListEncoderSimple(t *testing.T) {
	const expect = `
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000,AVERAGE-BANDWIDTH=1000000
http://example.com/low.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,AVERAGE-BANDWIDTH=2000000
http://example.com/mid.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,AVERAGE-BANDWIDTH=6000000
http://example.com/hi.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
http://example.com/audio-only.m3u8
`

	pl := MasterPlayList{
		Segments: []MasterSegment{
			{Stream: XStreamInf{Bandwidth: 1280000, AverageBandwidth: 1000000, URI: "http://example.com/low.m3u8"}},
			{Stream: XStreamInf{Bandwidth: 2560000, AverageBandwidth: 2000000, URI: "http://example.com/mid.m3u8"}},
			{Stream: XStreamInf{Bandwidth: 7680000, AverageBandwidth: 6000000, URI: "http://example.com/hi.m3u8"}},
			{Stream: XStreamInf{Bandwidth: 65000, Codecs: []string{"mp4a.40.5"}, URI: "http://example.com/audio-only.m3u8"}},
		},
	}

	var buf strings.Builder
	buf.Grow(512)
	if err := pl.encode(&buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}

func TestMasterPlayListEncoderWithIFrames(t *testing.T) {
	const expect = `
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000
low/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=86000,URI="low/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=2560000
mid/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=150000,URI="mid/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=7680000
hi/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=550000,URI="hi/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
audio-only.m3u8
`

	pl := MasterPlayList{
		Segments: []MasterSegment{
			{
				Stream: XStreamInf{Bandwidth: 1280000, URI: "low/audio-video.m3u8"},
			},
			{
				Stream: XStreamInf{Bandwidth: 2560000, URI: "mid/audio-video.m3u8"},
				IFrameStreams: []XIFrameStreamInf{
					{Bandwidth: 86000, URI: "low/iframe.m3u8"},
				},
			},
			{
				Stream: XStreamInf{Bandwidth: 7680000, URI: "hi/audio-video.m3u8"},
				IFrameStreams: []XIFrameStreamInf{
					{Bandwidth: 150000, URI: "mid/iframe.m3u8"},
				},
			},
			{
				Stream: XStreamInf{Bandwidth: 65000, Codecs: []string{"mp4a.40.5"}, URI: "audio-only.m3u8"},
				IFrameStreams: []XIFrameStreamInf{
					{Bandwidth: 550000, URI: "hi/iframe.m3u8"},
				},
			},
		},
	}

	var buf strings.Builder
	buf.Grow(512)
	if err := pl.encode(&buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}

func TestMasterPlayListEncoderWithAlternativeAudio(t *testing.T) {
	const expect = `
#EXTM3U
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="English",LANGUAGE="en",DEFAULT=YES,AUTOSELECT=YES,URI="main/english-audio.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="Deutsch",LANGUAGE="de",AUTOSELECT=YES,URI="main/german-audio.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="Commentary",LANGUAGE="en",URI="commentary/audio-only.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=1280000,CODECS="mp4a.40.5",AUDIO="aac"
low/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,CODECS="mp4a.40.5",AUDIO="aac"
mid/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,CODECS="mp4a.40.5",AUDIO="aac"
hi/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5",AUDIO="aac"
main/english-audio.m3u8
`

	pl := MasterPlayList{
		Segments: []MasterSegment{
			{
				Stream: XStreamInf{Bandwidth: 1280000, Codecs: []string{"mp4a.40.5"}, Audio: "aac", URI: "low/video-only.m3u8"},
				Medias: []XMedia{
					{Type: "AUDIO", GroupId: "aac", Name: "English", Default: true, AutoSelect: true, Language: "en", URI: "main/english-audio.m3u8"},
					{Type: "AUDIO", GroupId: "aac", Name: "Deutsch", Default: false, AutoSelect: true, Language: "de", URI: "main/german-audio.m3u8"},
					{Type: "AUDIO", GroupId: "aac", Name: "Commentary", Default: false, AutoSelect: false, Language: "en", URI: "commentary/audio-only.m3u8"},
				},
			},
			{Stream: XStreamInf{Bandwidth: 2560000, Codecs: []string{"mp4a.40.5"}, Audio: "aac", URI: "mid/video-only.m3u8"}},
			{Stream: XStreamInf{Bandwidth: 7680000, Codecs: []string{"mp4a.40.5"}, Audio: "aac", URI: "hi/video-only.m3u8"}},
			{Stream: XStreamInf{Bandwidth: 65000, Codecs: []string{"mp4a.40.5"}, Audio: "aac", URI: "main/english-audio.m3u8"}},
		},
	}

	var buf strings.Builder
	buf.Grow(512)
	if err := pl.encode(&buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}

func TestMasterPlayListEncoderWithAlternativeVideo(t *testing.T) {
	const expect = `
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

	pl := MasterPlayList{
		Segments: []MasterSegment{
			{
				Stream: XStreamInf{Bandwidth: 1280000, Codecs: []string{"mp4a.40.5"}, Video: "low", URI: "low/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: "VIDEO", GroupId: "low", Name: "Main", Default: true, URI: "low/main/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "low", Name: "Centerfield", Default: false, URI: "low/centerfield/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "low", Name: "Dugout", Default: false, URI: "low/dugout/audio-video.m3u8"},
				},
			},
			{
				Stream: XStreamInf{Bandwidth: 2560000, Codecs: []string{"mp4a.40.5"}, Video: "mid", URI: "mid/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: "VIDEO", GroupId: "mid", Name: "Main", Default: true, URI: "mid/main/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "mid", Name: "Centerfield", Default: false, URI: "mid/centerfield/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "mid", Name: "Dugout", Default: false, URI: "mid/dugout/audio-video.m3u8"},
				},
			},
			{
				Stream: XStreamInf{Bandwidth: 7680000, Codecs: []string{"mp4a.40.5"}, Video: "hi", URI: "hi/main/audio-video.m3u8"},
				Medias: []XMedia{
					{Type: "VIDEO", GroupId: "hi", Name: "Main", Default: true, URI: "hi/main/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "hi", Name: "Centerfield", Default: false, URI: "hi/centerfield/audio-video.m3u8"},
					{Type: "VIDEO", GroupId: "hi", Name: "Dugout", Default: false, URI: "hi/dugout/audio-video.m3u8"},
				},
			},
		},
	}

	var buf strings.Builder
	buf.Grow(512)
	if err := pl.encode(&buf); err != nil {
		t.Fatal(err)
	} else if s := buf.String(); s != expect[1:] {
		t.Errorf("expected:\n%s\ngot:\n%s", expect[1:], s)
	}
}
