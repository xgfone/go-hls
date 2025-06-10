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

import "io"

func (pl MasterPlayList) encode(w io.Writer) (err error) {
	// Basic Tags
	err = tryWriteString(w, nil, string(EXTM3U)+"\n")
	if version := pl.MinVersion(); version > 1 {
		err = tryWriteTag(w, err, EXT_X_VERSION, _DecimalInteger(version))
	}

	// Master/Media PlayList Tags
	err = tryWriteTag(w, err, EXT_X_INDEPENDENT_SEGMENTS, _Bool(pl.IndependentSegments))
	err = tryWriteTag(w, err, EXT_X_START, pl.Start)

	for _, s := range pl.Streams {
		if err != nil {
			break
		}

		err = tryWriteMasterTags(w, err, EXT_X_SESSION_KEY, s.SessionKeys)
		err = tryWriteMasterTags(w, err, EXT_X_SESSION_DATA, s.SessionDatas)
		err = tryWriteMasterTags(w, err, EXT_X_MEDIA, s.Medias)
		err = tryWriteMasterTags(w, err, EXT_X_I_FRAME_STREAM_INF, s.IFrameStreams)

		err = tryWriteTag(w, err, EXT_X_STREAM_INF, s.Stream)
	}

	return
}

func tryWriteMasterTags[T _Value](w io.Writer, err error, tag Tag, attrs []T) error {
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		err = tryWriteTag(w, err, tag, attr)
		if err != nil {
			return err
		}
	}

	return err
}
