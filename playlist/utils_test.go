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
