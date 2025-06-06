package playlist

import (
	"errors"
	"io"
)

type MasterSegment struct {
	Stream XStreamInf

	Medias        []XMedia
	IFrameStreams []XIFrameStreamInf
	SessionDatas  []XSessionData
	SessionKeys   []XKey
}

type MasterPlayList struct {
	Version uint64
	Start   XStart

	Segments []MasterSegment

	IndependentSegments bool
}

func (pl MasterPlayList) Type() string {
	return PlayListTypeMaster
}

func (pl MasterPlayList) MinVersion() uint64 {
	if pl.Version > 0 {
		return pl.Version
	}
	return 1
}

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
