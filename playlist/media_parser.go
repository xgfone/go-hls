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

var (
	errMissingMediaSegments = errors.New("missing media segments")
)

// Parse reads the data from r and parses it as the media playlist.
func (pl *MediaPlayList) Parse(r io.Reader) (err error) {
	var p _Parser
	p.initMedia()

	_pl, err := p.Parse(r)
	if err == nil {
		*pl = _pl.(MediaPlayList)
	}

	return
}

type _MediaPlayList struct {
	parser *_Parser
	media  MediaPlayList

	// Media Segment
	segcache MediaSegment
	curseg   *MediaSegment
}

func (p *_MediaPlayList) PlayList() MediaPlayList {
	p.media.IndependentSegments = p.parser.independentSegments
	p.media.Version = p.parser.version
	p.media.Start = p.parser.start
	p.media.update()
	return p.media
}

func (p *_MediaPlayList) end() bool {
	return p != nil && p.media.EndList
}

func (p *_MediaPlayList) setURI(uri string) {
	if p.curseg != nil {
		if len(p.curseg.Keys) == 0 && len(p.media.Segments) > 0 {
			p.curseg.Keys = p.media.Segments[len(p.media.Segments)-1].Keys
		}
		p.curseg.URI = uri
		p.media.Segments = append(p.media.Segments, *p.curseg)
		p.curseg = nil
	}
}

func (p *_MediaPlayList) initCurrentMediaSegment() {
	if p.curseg == nil {
		p.segcache = MediaSegment{}
		p.curseg = &p.segcache
		p.parser.uri = p
	}
}

func (p *_Parser) initMedia() {
	if p.mediapl == nil {
		p.mediapl = &_MediaPlayList{parser: p}
	}
}

func (p *_Parser) checkForMedia() (err error) {
	if p.mediapl == nil {
		return
	}
	return p.mediapl.media.validate(p.version)
}

func (p *_Parser) parseTagForMedia(tag Tag, attr string) (ok bool, err error) {
	switch tag {
	case
		////// Media Segment Tags
		EXTINF,
		EXT_X_KEY,
		EXT_X_MAP,
		EXT_X_BYTERANGE,
		EXT_X_DISCONTINUITY,
		EXT_X_PROGRAM_DATE_TIME,
		EXT_X_DATERANGE,

		////// Media Playlist Tags
		EXT_X_TARGETDURATION,
		EXT_X_MEDIA_SEQUENCE,
		EXT_X_DISCONTINUITY_SEQUENCE,
		EXT_X_PLAYLIST_TYPE,
		EXT_X_I_FRAMES_ONLY,
		EXT_X_ENDLIST:

	default:
		return
	}

	ok = true
	if p.masterpl != nil {
		err = errMixedMasterMedia
		return
	}

	p.initMedia()
	err = p.mediapl.parseTag(tag, attr)
	return
}

func (p *_MediaPlayList) parseTag(tag Tag, attr string) (err error) {
	parser := p.parser
	switch tag {
	////// Media Segment Tags
	case EXTINF:
		// RFC 8216, 4.3.2.1:
		// It applies only to the next Media Segment.
		// This tag is REQUIRED for each Media Segment.
		p.initCurrentMediaSegment()

		items := splitAttributes(attr, 2)
		if len(items) == 2 {
			p.curseg.Title = items[1]
		}

		var duration _DecimalFloat
		if err = duration.decode(items[0]); err == nil {
			if duration <= 0 {
				err = errInvalidDecimalFloat
			} else {
				p.curseg.Duration = duration.get()
			}
		}

	case EXT_X_BYTERANGE:
		// RFC 8216, 4.3.2.2:
		// It applies only to the next URI line that follows it in the Playlist.
		p.initCurrentMediaSegment()
		err = p.curseg.ByteRange.decode(attr)

	case EXT_X_DISCONTINUITY:
		// RFC 8216, 4.3.2.3:
		// 1. Must be set when changed:
		//    - file format,
		//    - number, type, and identifiers of tracks
		//    - timestamp sequence
		// 2. SHOULD be set when changed:
		//    - encoding parameters
		//    - encoding sequence
		p.initCurrentMediaSegment()
		p.curseg.Discontinuity = true

	case EXT_X_KEY:
		// RFC 8216, 4.3.2.4:
		// It applies to every Media Segment and to every Media Initialization Section
		// declared by an EXT-X-MAP tag that appears between it and the next EXT-X-KEY tag
		// in the Playlist file with the same KEYFORMAT attribute (or the end of the Playlist file).
		//
		// Two or more EXT-X-KEY tags with different KEYFORMAT attributes MAY apply
		// to the same Media Segment if they ultimately produce the same decryption key.
		p.initCurrentMediaSegment()
		var key XKey
		if err = key.decode(attr); err == nil {
			p.curseg.Keys = append(p.curseg.Keys, key)
		}

	case EXT_X_MAP:
		// RFC 8216, 4.3.2.5:
		// It applies to every Media Segment that appears after it in the Playlist
		// until the next EXT-X-MAP tag or until the end of the Playlist.
		p.initCurrentMediaSegment()
		err = p.curseg.Map.decode(attr)

	case EXT_X_PROGRAM_DATE_TIME:
		// RFC 8216, 4.3.2.6:
		// It applies only to the next Media Segment.
		p.initCurrentMediaSegment()
		var time _Time
		if err = time.decode(attr); err == nil {
			p.curseg.ProgramDateTime = time.get()
		}

	// case EXT_X_DATERANGE: // RFC 8216, 4.3.2.7:
	// p.initCurrentMediaSegment()

	////// Media Playlist Tags
	case EXT_X_TARGETDURATION:
		// RFC 8216, 4.3.3.1:
		// It applies to the entire Playlist file.
		if p.media.TargetDuration > 0 && parser.strict {
			err = errDuplicatedTag
		} else {
			var duration _DecimalInteger
			if err = duration.decode(attr, 1); err == nil {
				p.media.TargetDuration = duration.get()
			}
		}

	case EXT_X_MEDIA_SEQUENCE:
		// RFC 8216, 4.3.3.2:
		// It indicates the Media Sequence Number of the first Media Segment
		// that appears in a Playlist file.
		if p.media.MediaSequence > 0 && parser.strict {
			err = errDuplicatedTag
		} else {
			var seq _DecimalInteger
			if err = seq.decode(attr, 1); err == nil {
				p.media.MediaSequence = seq.get()
			}
		}

	case EXT_X_DISCONTINUITY_SEQUENCE:
		// RFC 8216, 4.3.3.3:
		// It MUST appear before the first Media Segment in the Playlist.
		// It MUST appear before any EXT-X-DISCONTINUITY tag.
		if p.media.DiscontinuitySequence > 0 && parser.strict {
			err = errDuplicatedTag
		} else {
			var seq _DecimalInteger
			if err = seq.decode(attr, 1); err == nil {
				if len(p.media.Segments) > 0 || p.curseg != nil {
					err = errors.New("must appear before any media segments")
				} else {
					p.media.DiscontinuitySequence = seq.get()
				}
			}
		}

	case EXT_X_ENDLIST:
		// RFC 8216, 4.3.3.4:
		// It indicates that no more Media Segments will be added to the Media Playlist file.
		// It MAY occur anywhere in the Media Playlist file.
		if p.media.EndList && parser.strict {
			err = errDuplicatedTag
		} else {
			p.media.EndList = true
		}

	case EXT_X_PLAYLIST_TYPE:
		// RFC 8216, 4.3.3.5:
		// It applies to the entire Media Playlist file.
		if p.media.PlayListType != "" && parser.strict {
			err = errDuplicatedTag
		} else {
			var _type _Enum
			if err = _type.decode(attr); err == nil {
				p.media.PlayListType = _type.get()
			}
		}

	case EXT_X_I_FRAMES_ONLY:
		// RFC 8216, 4.3.3.6:
		// It applies to the entire Playlist.
		if p.media.IFrameOnly && parser.strict {
			err = errDuplicatedTag
		} else {
			p.media.IFrameOnly = true
		}
	}

	return
}
