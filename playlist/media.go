package playlist

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"
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

// MediaSegment represents a media segment in a media playlist.
//
// See [[RFC 8216, 4.3.2]].
//
// [RFC 8216, 4.3.2]: https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2
type MediaSegment struct {
	URI   string // Required.
	Title string

	Duration  float64 // Required. Unit: Second
	ByteRange XByteRange
	Key       XKey
	Map       XMap

	ProgramDateTime time.Time

	MediaSequence         uint64 // Cannot be encoded
	DiscontinuitySequence uint64 // Cannot be encoded

	Discontinuity bool
}

// IV try to decode the iv from a hexadecimal-sequence string to a 16-octet bytes.
func (s MediaSegment) IV() (data []byte, err error) {
	if s.Key.IV != "" {
		var seq _HexSequence
		err = seq.decode(s.Key.IV)
		data = []byte(seq)
	} else {
		data = make([]byte, 16)
		binary.BigEndian.PutUint64(data[8:], s.MediaSequence)
	}
	return
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
	for i := range pl.Segments {
		s := &pl.Segments[i]
		s.MediaSequence = pl.calcMediaSequence()
		s.DiscontinuitySequence = pl.calcDiscontinuitySequence()
	}
}

func (pl *MediaPlayList) calcMediaSequence() uint64 {
	return atomic.AddUint64(&pl.MediaSequence, 1) - 1
}

func (pl *MediaPlayList) calcDiscontinuitySequence() uint64 {
	return atomic.AddUint64(&pl.DiscontinuitySequence, 1) - 1
}
