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
	"fmt"
	"io"
)

// Media PlayList Types.
const (
	MediaPlayListTypeVOD   = "VOD"
	MediaPlayListTypeEvent = "EVENT"
)

// MediaPlayList represents a media playlist, which implemented the PlayList interface.
type MediaPlayList struct {
	Version uint64 `json:",omitempty,omitzero"`

	Start    XStart         `json:",omitzero"`
	Segments []MediaSegment `json:",omitempty,omitzero"`

	TargetDuration        uint64 `json:",omitempty,omitzero"` // Unit: second
	MediaSequence         uint64 `json:",omitempty,omitzero"`
	DiscontinuitySequence uint64 `json:",omitempty,omitzero"`
	PlayListType          string `json:",omitempty,omitzero"`
	IndependentSegments   bool   `json:",omitempty,omitzero"`
	IFrameOnly            bool   `json:",omitempty,omitzero"`
	EndList               bool   `json:",omitempty,omitzero"`
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

func (pl MediaPlayList) minVersion() (minVersion uint64) {
	setVersion := func(version uint64) {
		minVersion = max(minVersion, version)
	}

	for _, seg := range pl.Segments {
		if !isIntegerFloat64(seg.Duration) {
			setVersion(3)
		}

		setVersion(seg.ByteRange.minVersion())
		for i := range seg.Keys {
			setVersion(seg.Keys[i].minVersion())
		}
		if seg.Map.valid() {
			if pl.IFrameOnly {
				setVersion(5)
			} else {
				setVersion(6)
			}
		}
		if pl.IFrameOnly {
			setVersion(4)
		}
	}

	return
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
	index := -1
	lastdseq := pl.DiscontinuitySequence
	lastmseq := pl.MediaSequence
	for i := range pl.Segments {
		s := &pl.Segments[i]
		s.MediaSequence = lastmseq
		lastmseq++

		if s.Discontinuity {
			lastdseq++
		}
		s.DiscontinuitySequence = lastdseq

		if index == -1 && !s.ProgramDateTime.IsZero() {
			index = i
		}
	}

	// Recover the Media Sequence Number parsed by #EXT-X-MEDIA-SEQUENCE.
	if len(pl.Segments) > 0 {
		pl.MediaSequence = pl.Segments[0].MediaSequence
	}

	if index > -1 {
		pl._updateProgramDateTime(index)
	}
}

func (pl *MediaPlayList) _updateProgramDateTime(index int) {
	// 1. Backward
	{
		lastseg := &pl.Segments[index]
		for i := index - 1; i >= 0; i-- {
			seg := &pl.Segments[i]
			seg.ProgramDateTime = lastseg.nextProgramDateTime(-seg.Duration)
			lastseg = seg
		}
	}

	// 2. Forward
	{
		lastseg := &pl.Segments[index]
		for i, _len := index+1, len(pl.Segments); i < _len; i++ {
			seg := &pl.Segments[i]
			if seg.ProgramDateTime.IsZero() {
				seg.ProgramDateTime = lastseg.nextProgramDateTime(lastseg.Duration)
			}
			lastseg = seg
		}
	}
}
