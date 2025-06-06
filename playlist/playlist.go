package playlist

import "io"

// Define the playlist types.
const (
	PlayListTypeMaster = "Master"
	PlayListTypeMedia  = "Media"
)

// PlayList is a playlist interface.
type PlayList interface {
	Type() string // Master or Media
	MinVersion() uint64
}

// Parse reads the data from r and decodes it as a master or media playlist.
func Parse(r io.Reader) (PlayList, error) {
	var p _Parser
	return p.Parse(r)
}
