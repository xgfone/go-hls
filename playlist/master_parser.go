package playlist

import "errors"

type _MasterPlayList struct {
	parser  *_Parser
	Version uint64

	master MasterPlayList
}

func (p *_MasterPlayList) PlayList() MasterPlayList {
	p.master.Version = p.parser.version
	p.master.Start = p.parser.start
	// TODO:
	return p.master
}

func (p *_Parser) checkForMaster() (err error) {
	if p.masterpl == nil {
		return
	}

	if err = checkXMedias(p.masterpl.master.Medias); err != nil {
		return
	}

	// TODO:

	return
}

func (p *_Parser) parseTagForMaster(tag Tag, _ string) (ok bool, err error) {
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

	err = errors.New("not implemented for master playlist")
	return
}
