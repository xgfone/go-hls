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
	"errors"
	"io"
)

// MasterStream represents a master stream in a master playlist.
type MasterStream struct {
	Stream XStreamInf `json:",omitzero"`

	Medias        []XMedia           `json:",omitempty,omitzero"`
	IFrameStreams []XIFrameStreamInf `json:",omitempty,omitzero"`
	SessionDatas  []XSessionData     `json:",omitempty,omitzero"`
	SessionKeys   []XKey             `json:",omitempty,omitzero"`
}

// MasterPlayList represents a master playlist, which implemented the PlayList interface.
type MasterPlayList struct {
	Version uint64 `json:",omitempty,omitzero"`
	Start   XStart `json:",omitzero"`

	Streams []MasterStream `json:",omitempty,omitzero"`

	IndependentSegments bool `json:",omitempty,omitzero"`
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
	return pl.minVersion()
}

// Output encodes the master playlist as the M3U8 format to w.
func (pl MasterPlayList) Output(w io.Writer) error {
	if err := pl.validate(); err != nil {
		return err
	}
	return pl.encode(w)
}

func (pl MasterPlayList) minVersion() (minVersion uint64) {
	setVersion := func(version uint64) {
		minVersion = max(minVersion, version)
	}

	for _, s := range pl.Streams {
		for _, m := range s.Medias {
			setVersion(m.minVersion())
		}
	}

	return
}

func (pl MasterPlayList) validate() (err error) {
	for _, s := range pl.Streams {
		if s.Stream.URI == "" {
			return errors.New(string(EXT_X_STREAM_INF) + ": missing URI")
		}
		if err = checkXMedias(s.Medias); err != nil {
			return err
		}

		for _, key := range s.SessionKeys {
			if key.Method == XKeyMethodNone {
				return errors.New(string(EXT_X_SESSION_KEY) + ": METHOD must not be NONE")
			}
		}
	}

	return
}
