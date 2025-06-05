package playlist

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/textproto"
	"strings"
)

var (
	errDuplicatedTag    = errors.New("duplicated tag")
	errTooLowerVersion  = errors.New("too lower version")
	errNotMasterOrMedia = errors.New("not master or media playlist")
	errMixedMasterMedia = errors.New("mixed master and media playlist tags")
)

type _URISetter interface {
	setURI(uri string)
}

type _Parser struct {
	lineno  int
	version uint64

	uri      _URISetter
	start    XStart
	masterpl *_MasterPlayList
	mediapl  *_MediaPlayList

	// Master/Media PlayList
	independentSegments bool

	strict bool
}

func (p *_Parser) readline(r *textproto.Reader) (line string, err error) {
	for {
		p.lineno++
		if line, err = r.ReadLine(); err != nil {
			return
		}

		switch line = strings.TrimSpace(line); {
		case line == "": // Blank Line
		case line[0] == '#' && !strings.HasPrefix(line, "#EXT"): // Comment Line
		default:
			return
		}
	}
}

func (p *_Parser) Parse(r io.Reader) (pl PlayList, err error) {
	if err = p.parse(r); err != nil {
		err = ParseError{Line: p.lineno, Err: err}
		return
	}

	if err = p.checkForMaster(); err != nil {
		return
	}

	if err = p.checkForMedia(); err != nil {
		return
	}

	switch {
	case p.masterpl != nil:
		pl = p.masterpl.PlayList()

	case p.mediapl != nil:
		pl = p.mediapl.PlayList()

	default:
		err = errNotMasterOrMedia
	}

	return
}

func (p *_Parser) parse(r io.Reader) (err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	reader := textproto.NewReader(br)
	line, err := p.readline(reader)
	if err != nil {
		return
	} else if line != string(EXTM3U) {
		return errors.New("not start with " + string(EXTM3U))
	}

	for {
		if line, err = p.readline(reader); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		if line[0] == '#' {
			err = p.parseLineForTag(line)
		} else {
			err = p.parseLineForURI(line)
		}

		if err != nil || p.mediapl.end() {
			return
		}
	}
}

func (p *_Parser) parseLineForURI(line string) (err error) {
	if p.uri == nil {
		return errInvalidURI
	}

	var url _UnquotedString
	if err = url.decode(line); err != nil {
		return fmt.Errorf("invalid URI: %w", err)
	}

	p.uri.setURI(url.get())
	return
}

func (p *_Parser) parseLineForTag(line string) (err error) {
	var attr string
	tag := Tag(line)
	if index := strings.IndexByte(line, ':'); index > 0 {
		tag = Tag(line[:index])
		attr = line[index+1:]
	}

	var ok bool
	switch tag {

	////// Basic Tags
	case EXTM3U:
		if p.strict {
			err = errDuplicatedTag
		}

	case EXT_X_VERSION:
		if p.version > 0 && p.strict {
			err = errDuplicatedTag
		} else {
			var version _DecimalInteger
			if err = version.decode(attr, 1); err == nil {
				if p.version = version.get(); p.version == 0 {
					err = errors.New("version must be equal to 0")
				}
			}
		}

	////// Media or Master Playlist Tags
	case EXT_X_INDEPENDENT_SEGMENTS:
		// RFC 8216, 4.3.5.1:
		// It applies to every Media Segment in the Playlist.
		// If appears in a Master Playlist, it applies to every Media Segment
		// in every Media Playlist in the Master Playlist.
		p.independentSegments = true

	case EXT_X_START:
		// RFC 8216, 4.3.5.2:
		err = p.start.decode(attr)

	default:
		if ok, err = p.parseTagForMaster(tag, attr); err == nil && !ok {
			ok, err = p.parseTagForMedia(tag, attr)
		}
	}

	if err != nil {
		err = errors.New(strings.Join([]string{string(tag), err.Error()}, ": "))
	} else if !ok {
		slog.Debug("unknown tag", "tag", tag, "attr", attr)
	}

	return
}
