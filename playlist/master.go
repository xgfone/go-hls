package playlist

import (
	"errors"
	"io"
)

// MasterSegment represents a master segment in a master playlist.
type MasterSegment struct {
	Stream XStreamInf

	Medias        []XMedia
	IFrameStreams []XIFrameStreamInf
	SessionDatas  []XSessionData
	SessionKeys   []XKey
}

// MasterPlayList represents a master playlist, which implemented the PlayList interface.
type MasterPlayList struct {
	Version uint64
	Start   XStart

	Segments []MasterSegment

	IndependentSegments bool
}

// Type returns the fixed "Master".
func (pl MasterPlayList) Type() string {
	return PlayListTypeMaster
}

// MinVersion returns the mininal version.
func (pl MasterPlayList) MinVersion() uint64 {
	if pl.Version > 0 {
		return pl.Version
	}
	return 1
}

// Output encodes the master playlist as the M3U8 format to w.
func (pl MasterPlayList) Output(w io.Writer) error {
	if err := pl.validate(); err != nil {
		return err
	}
	return pl.encode(w)
}

func (pl MasterPlayList) validate() (err error) {
	for _, seg := range pl.Segments {
		if seg.Stream.URI == "" {
			return errors.New(string(EXT_X_STREAM_INF) + ": missing URI")
		}
		if err = checkXMedias(seg.Medias); err != nil {
			return err
		}

		for _, key := range seg.SessionKeys {
			if key.Method == XKeyMethodNone {
				return errors.New(string(EXT_X_SESSION_KEY) + ": METHOD must not be NONE")
			}
		}
	}

	return
}
