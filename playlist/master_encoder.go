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

	for _, seg := range pl.Segments {
		if err != nil {
			break
		}

		err = tryWriteMasterTags(w, err, EXT_X_SESSION_KEY, seg.SessionKeys)
		err = tryWriteMasterTags(w, err, EXT_X_SESSION_DATA, seg.SessionDatas)
		err = tryWriteMasterTags(w, err, EXT_X_MEDIA, seg.Medias)
		err = tryWriteMasterTags(w, err, EXT_X_I_FRAME_STREAM_INF, seg.IFrameStreams)

		err = tryWriteTag(w, err, EXT_X_STREAM_INF, seg.Stream)
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
