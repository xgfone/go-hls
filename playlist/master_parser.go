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

// Parse reads the data from r and parses it as the master playlist.
func (pl *MasterPlayList) Parse(r io.Reader) (err error) {
	var p _Parser
	p.initMaster()

	_pl, err := p.Parse(r)
	if err == nil {
		*pl = _pl.(MasterPlayList)
	}

	return
}

type _MasterPlayList struct {
	parser *_Parser

	master MasterPlayList

	curseg   *MasterSegment
	segcache MasterSegment
}

func (p *_MasterPlayList) PlayList() MasterPlayList {
	p.master.IndependentSegments = p.parser.independentSegments
	p.master.Version = p.parser.version
	p.master.Start = p.parser.start
	return p.master
}

func (p *_MasterPlayList) setURI(uri string) {
	if p.curseg != nil {
		p.curseg.Stream.URI = uri
		p.master.Segments = append(p.master.Segments, *p.curseg)
		p.curseg = nil
	}
}

func (p *_MasterPlayList) initCurrentSegment() {
	if p.curseg == nil {
		p.segcache = MasterSegment{}
		p.curseg = &p.segcache
		p.parser.uri = p
	}
}

func (p *_Parser) initMaster() {
	if p.masterpl == nil {
		p.masterpl = &_MasterPlayList{parser: p}
	}
}

func (p *_Parser) checkForMaster() (err error) {
	if p.masterpl == nil {
		return
	}
	return p.masterpl.master.validate()
}

func (p *_Parser) parseTagForMaster(tag Tag, attr string) (ok bool, err error) {
	switch tag {
	case EXT_X_MEDIA,
		EXT_X_STREAM_INF,
		EXT_X_I_FRAME_STREAM_INF,
		EXT_X_SESSION_DATA,
		EXT_X_SESSION_KEY:

	default:
		return
	}

	ok = true
	if p.mediapl != nil {
		err = errMixedMasterMedia
		return
	}

	p.initMaster()
	err = p.masterpl.parseTag(tag, attr)
	return
}

func (p *_MasterPlayList) parseTag(tag Tag, attr string) (err error) {
	switch tag {
	case EXT_X_MEDIA:
		// RFC 8216, 4.3.4.1:
		p.initCurrentSegment()
		var media XMedia
		if err = media.decode(attr); err == nil {
			p.curseg.Medias = append(p.curseg.Medias, media)
		}

	case EXT_X_STREAM_INF:
		// RFC 8216, 4.3.4.2:
		p.initCurrentSegment()
		err = p.curseg.Stream.decode(attr)

	case EXT_X_I_FRAME_STREAM_INF:
		// RFC 8216, 4.3.4.3:
		p.initCurrentSegment()
		var stream XIFrameStreamInf
		if err = stream.decode(attr); err == nil {
			p.curseg.IFrameStreams = append(p.curseg.IFrameStreams, stream)
		}

	case EXT_X_SESSION_DATA:
		// RFC 8216, 4.3.4.4:
		p.initCurrentSegment()
		var sdata XSessionData
		if err = sdata.decode(attr); err == nil {
			p.curseg.SessionDatas = append(p.curseg.SessionDatas, sdata)
		}

	case EXT_X_SESSION_KEY:
		// RFC 8216, 4.3.4.5:
		p.initCurrentSegment()
		var xkey XKey
		if err = xkey.decode(attr); err == nil {
			if xkey.Method == XKeyMethodNone {
				err = errors.New("METHOD must not be NONE")
			} else {
				p.curseg.SessionKeys = append(p.curseg.SessionKeys, xkey)
			}
		}
	}

	return
}
