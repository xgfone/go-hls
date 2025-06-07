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
	"io"
	"slices"
)

func (pl MediaPlayList) encode(w io.Writer) (err error) {
	// Basic Tags
	err = tryWriteString(w, nil, string(EXTM3U)+"\n")
	if version := pl.MinVersion(); version > 1 {
		err = tryWriteTag(w, err, EXT_X_VERSION, _DecimalInteger(version))
	}

	// Master/Media PlayList Tags
	err = tryWriteTag(w, err, EXT_X_INDEPENDENT_SEGMENTS, _Bool(pl.IndependentSegments))
	err = tryWriteTag(w, err, EXT_X_START, pl.Start)

	// Media PlayList Tags
	err = tryWriteTag(w, err, EXT_X_PLAYLIST_TYPE, newEnum(pl.PlayListType))
	err = tryWriteTag(w, err, EXT_X_TARGETDURATION, _DecimalInteger(pl.TargetDuration))
	err = tryWriteTag(w, err, EXT_X_I_FRAMES_ONLY, _Bool(pl.IFrameOnly))
	err = tryWriteTag(w, err, EXT_X_MEDIA_SEQUENCE, _DecimalInteger(pl.MediaSequence))
	err = tryWriteTag(w, err, EXT_X_DISCONTINUITY_SEQUENCE, _DecimalInteger(pl.DiscontinuitySequence))

	// Media Segment Tags
	lastkeys := make([]XKey, 0, 4)
	curkeys := make([]XKey, 0, 4)

	var xmap XMap
	for _, seg := range pl.Segments {
		if err != nil {
			break
		}

		switch {
		case seg.URI == "":
			panic("missing URI in MediaSegment")
		case seg.Duration <= 0:
			panic("missing Duration in MediaSegment")
		}

		curkeys = append(curkeys[:0], seg.Keys...)
		sortKeys(curkeys)
		switch {
		case len(curkeys) == 0:
		case len(lastkeys) == 0:
			lastkeys = append(lastkeys[:0], curkeys...)

		default:
			if !equalKeys(lastkeys, curkeys) {
				lastkeys = append(lastkeys[:0], curkeys...)
			} else {
				seg.Keys = nil
			}
		}

		switch {
		case xmap.IsZero(), xmap != seg.Map:
			xmap = seg.Map

		case seg.Map.IsZero():
			seg.Map = xmap

		case seg.Map == xmap:
			seg.Map = XMap{}
		}

		for _, key := range seg.Keys {
			err = tryWriteTag(w, err, EXT_X_KEY, key)
		}
		err = tryWriteTag(w, err, EXT_X_MAP, seg.Map)

		err = tryWriteTag(w, err, EXT_X_DISCONTINUITY, _Bool(seg.Discontinuity))
		err = tryWriteTag(w, err, EXT_X_PROGRAM_DATE_TIME, _Time(seg.ProgramDateTime))
		err = tryWriteTag(w, err, EXT_X_BYTERANGE, seg.ByteRange)
		err = tryWriteAny(w, err, string(EXTINF+":"), _DecimalFloat(seg.Duration), ",", _UnquotedString(seg.Title), "\n")
		err = tryWrite(w, err, _UnquotedString(seg.URI))
		err = tryWriteString(w, err, "\n")
	}

	err = tryWriteTag(w, err, EXT_X_ENDLIST, _Bool(pl.EndList))
	return
}

func equalKeys(keys1, keys2 []XKey) bool {
	if len(keys1) != len(keys2) {
		return false
	}

	for i := range keys1 {
		if keys1[i] != keys2[i] {
			return false
		}
	}

	return true
}

func sortKeys(keys []XKey) {
	if len(keys) < 2 {
		return
	}

	slices.SortStableFunc(keys, func(a, b XKey) int {
		switch {
		case a == b:
			return 0

		case a.Format < b.Format || a.URI < b.URI:
			return -1

		default:
			return 1
		}
	})
}
