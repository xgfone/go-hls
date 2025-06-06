package playlist

import (
	"errors"
	"fmt"
	"io"
	"sync/atomic"
)

// Media PlayList Types.
const (
	MediaPlayListTypeVOD   XMediaPlayListType = "VOD"
	MediaPlayListTypeEvent XMediaPlayListType = "EVENT"
)

// XMediaPlayListType is used to define the type of the media playlist.
type XMediaPlayListType string

func (t XMediaPlayListType) validate() error {
	switch t {
	case MediaPlayListTypeVOD, MediaPlayListTypeEvent:
		return nil
	default:
		return errors.New("invalid media playlist type")
	}
}

// MediaPlayList represents a media playlist, which implemented the PlayList interface.
type MediaPlayList struct {
	Version uint64

	Start    XStart
	Segments []MediaSegment

	TargetDuration        uint64 // Unit: second
	MediaSequence         uint64
	DiscontinuitySequence uint64
	PlayListType          XMediaPlayListType
	IndependentSegments   bool
	IFrameOnly            bool
	EndList               bool
}

// Type returns the fixed "Media".
func (pl MediaPlayList) Type() string {
	return PlayListTypeMedia
}

// MinVersion returns the mininal version.
func (pl MediaPlayList) MinVersion() uint64 {
	if pl.Version > 0 {
		return pl.Version
	}
	return pl.minVersion()
}

// TotalDuration calculates and returns the total duration of media segments.
func (pl MediaPlayList) TotalDuration() float64 {
	var total float64
	for _, seg := range pl.Segments {
		total += seg.Duration
	}
	return total
}

// Output encodes the media playlist as the M3U8 format to w.
func (pl MediaPlayList) Output(w io.Writer) error {
	if err := pl.validate(pl.Version); err != nil {
		return err
	}
	return pl.encode(w)
}

func (pl MediaPlayList) minVersion() uint64 {
	for _, seg := range pl.Segments {
		if !isIntegerFloat64(seg.Duration) {
			pl.setVersion(3)
		}

		pl.setVersion(seg.ByteRange.minVersion())
		pl.setVersion(seg.Key.minVersion())
		if seg.Map.valid() {
			if pl.IFrameOnly {
				pl.setVersion(5)
			} else {
				pl.setVersion(6)
			}
		}
		if pl.IFrameOnly {
			pl.setVersion(4)
		}
	}
	return pl.Version
}

func (pl *MediaPlayList) setVersion(version uint64) {
	if pl.Version < version {
		pl.Version = version
	}
}

func (pl MediaPlayList) validate(minVersion uint64) (err error) {
	version := pl.minVersion()
	if version > 1 && minVersion < version {
		return errTooLowerVersion
	}

	if len(pl.Segments) == 0 {
		return errMissingMediaSegments
	}

	for i, seg := range pl.Segments {
		if uint64(seg.Duration+0.5) > pl.TargetDuration {
			return fmt.Errorf("media segment duration exceeds target duration at %d", i)
		}
	}

	return
}

func (pl *MediaPlayList) update() {
	lastdseq := pl.DiscontinuitySequence
	for i := range pl.Segments {
		s := &pl.Segments[i]
		s.MediaSequence = pl.calcMediaSequence()

		if s.Discontinuity {
			lastdseq++
		}
		s.DiscontinuitySequence = lastdseq
	}

	// Recover the Media Sequence Number parsed by #EXT-X-MEDIA-SEQUENCE.
	if len(pl.Segments) > 0 {
		pl.MediaSequence = pl.Segments[0].MediaSequence
	}
}

func (pl *MediaPlayList) calcMediaSequence() uint64 {
	return atomic.AddUint64(&pl.MediaSequence, 1) - 1
}
