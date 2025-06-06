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

func TestXByteRange(t *testing.T) {
	notiframe := XByteRange{Length: 123456, Offset: 789008}
	iframe := XByteRange{Length: 123472, Offset: 788992}

	r := XByteRange{Length: 123456, Offset: 789012}
	if r = r.Align16(); r != notiframe {
		t.Errorf("expect %+v, but got %+v", notiframe, r)
	}

	if r, _ = r.AdjustForIFrame(); r != iframe {
		t.Errorf("expect %+v, but got %+v", iframe, r)
	}
}

func TestFormatIV(t *testing.T) {
	iv := FormatIV([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, true)
	expect := "0x0102030405060708090A0B0C0D0E0F10"
	if expect != iv {
		t.Errorf("expect IV '%s', but got '%s'", expect, iv)
	}
}
