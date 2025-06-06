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

import "testing"

func TestMediaSegmentIndex(t *testing.T) {
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
