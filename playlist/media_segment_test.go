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
	"testing"
	"time"
)

func TestMediaSegmentIndexByDuration(t *testing.T) {
	pl := MediaPlayList{
		TargetDuration: 10,
		Segments: []MediaSegment{
			{URI: "segment-1.ts", Duration: 9},
			{URI: "segment-2.ts", Duration: 8},
			{URI: "segment-3.ts", Duration: 7},
			{URI: "segment-4.ts", Duration: 6},
			{URI: "segment-5.ts", Duration: 5},
			{URI: "segment-6.ts", Duration: 4},
			{URI: "segment-7.ts", Duration: 3},
			{URI: "segment-8.ts", Duration: 2},
			{URI: "segment-9.ts", Duration: 1},
		},
	}

	if index := pl.GetSegmentIndexByDuration(1); index != 0 {
		t.Errorf("expect segment index %d, but got %d", 0, index)
	}
	if index := pl.GetSegmentIndexByDuration(30); index != 4 {
		t.Errorf("expect segment index %d, but got %d", 4, index)
	}
	if index := pl.GetSegmentIndexByDuration(44); index != 8 {
		t.Errorf("expect segment index %d, but got %d", 8, index)
	}
	if index := pl.GetSegmentIndexByDuration(45); index != -1 {
		t.Errorf("expect segment index %d, but got %d", -1, index)
	}
}

func TestMediaSegmentIndexByMediaSequence(t *testing.T) {
	pl := MediaPlayList{
		MediaSequence: 100,
		Segments:      make([]MediaSegment, 10),
	}
	pl.update()

	if index := pl.GetSegmentIndexByMediaSequence(99); index != -1 {
		t.Errorf("expect segment index %d, but got %d", -1, index)
	}
	if index := pl.GetSegmentIndexByMediaSequence(120); index != -1 {
		t.Errorf("expect segment index %d, but got %d", -1, index)
	}

	if index := pl.GetSegmentIndexByMediaSequence(100); index != 0 {
		t.Errorf("expect segment index %d, but got %d", 0, index)
	}
	if index := pl.GetSegmentIndexByMediaSequence(101); index != 1 {
		t.Errorf("expect segment index %d, but got %d", 1, index)
	}
	if index := pl.GetSegmentIndexByMediaSequence(109); index != 9 {
		t.Errorf("expect segment index %d, but got %d", 0, index)
	}
	if index := pl.GetSegmentIndexByMediaSequence(110); index != -1 {
		t.Errorf("expect segment index %d, but got %d", -1, index)
	}
}

func TestMediaSegmentProgramDateTime(t *testing.T) {
	musttime := func(s string) time.Time {
		t, err := time.ParseInLocation(time.DateTime, s, time.Local)
		if err != nil {
			panic(err)
		}
		return t
	}

	pl := MediaPlayList{
		Segments: []MediaSegment{
			{Duration: 1},
			{Duration: 2},
			{Duration: 3, ProgramDateTime: musttime("2025-06-07 00:01:10")},
			{Duration: 4},
			{Duration: 5},
			{Duration: 6, ProgramDateTime: musttime("2025-06-07 00:01:30")},
			{Duration: 7},
		},
	}
	pl.update()

	for i, s := range pl.Segments {
		dt := s.ProgramDateTime.Format(time.DateTime)
		switch j := i + 1; {
		case j == 1 && dt == "2025-06-07 00:01:07":
		case j == 2 && dt == "2025-06-07 00:01:08":
		case j == 3 && dt == "2025-06-07 00:01:10":
		case j == 4 && dt == "2025-06-07 00:01:13":
		case j == 5 && dt == "2025-06-07 00:01:17":
		case j == 6 && dt == "2025-06-07 00:01:30":
		case j == 7 && dt == "2025-06-07 00:01:36":
		default:
			t.Errorf("%d: unexpect time '%s'", i, dt)
		}
	}
}
