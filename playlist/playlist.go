package playlist

import "io"

const (
	PlayListTypeMaster = "Master"
	PlayListTypeMedia  = "Media"
)

type PlayList interface {
	Type() string // Master or Media
	MinVersion() uint64
}

func Parse(r io.Reader) (PlayList, error) {
	var p _Parser
	return p.Parse(r)
}
