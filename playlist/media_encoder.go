package playlist

import "io"

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
	var xkey XKey
	var xmap XMap
	for _, seg := range pl.Segments {
		if err != nil {
			break
		}

		switch {
		case xkey.IsZero(), xkey != seg.Key:
			xkey = seg.Key

		case seg.Key.IsZero():
			seg.Key = xkey

		case seg.Key == xkey:
			seg.Key = XKey{}
		}

		switch {
		case xmap.IsZero(), xmap != seg.Map:
			xmap = seg.Map

		case seg.Map.IsZero():
			seg.Map = xmap

		case seg.Map == xmap:
			seg.Map = XMap{}
		}

		err = tryWriteTag(w, err, EXT_X_KEY, seg.Key)
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
