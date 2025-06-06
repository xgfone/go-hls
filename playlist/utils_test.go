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
	"testing"
)

func expectAttrs(t *testing.T, expects, values []string) {
	if !reflect.DeepEqual(expects, values) {
		t.Errorf("expect %v, but got %v", expects, values)
	}
}

func TestSplitAttributes(t *testing.T) {
	expectAttrs(t, []string{"10"}, splitAttributes(`10`, 2))
	expectAttrs(t, []string{"10"}, splitAttributes(`10,`, 2))
	expectAttrs(t, []string{"10", "title"}, splitAttributes(`10,title`, 2))
	expectAttrs(t, []string{"10", "title,"}, splitAttributes(`10,title,`, 2))

	expectAttrs(t, []string{"A=1"}, splitAttributes(`A=1`, -1))
	expectAttrs(t, []string{"A=1", "B=2"}, splitAttributes(`A=1,B=2,`, -1))
	expectAttrs(t, []string{"A=1", `CODECS="mp4a"`}, splitAttributes(`A=1,CODECS="mp4a"`, -1))
	expectAttrs(t, []string{"A=1", `CODECS="mp4a,mp4b"`}, splitAttributes(`A=1,CODECS="mp4a,mp4b"`, -1))
}

func TestFloat64ToDuration(t *testing.T) {
	duration := float64ToDuration(1.2).String()
	if duration != "1.2s" {
		t.Errorf("expect duration %s, but got %s", "1.2s", duration)
	}
}
