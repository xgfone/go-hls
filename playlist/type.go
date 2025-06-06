package playlist

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	errInvalidURI        = errors.New("invalid URI")
	errInvalidKeyMethod  = errors.New("invalid key method")
	errInvalidByteRange  = errors.New("invalid byte range")
	errInvalidResolution = errors.New("invalid resolution")
	errInvalidHDCPLevel  = errors.New("invalid HDCP level")
)

/// ----------------------------------------------------------------------- ///

// XByteRange represents a byte range for download.
type XByteRange struct {
	Length uint64
	Offset uint64
}

func (v XByteRange) minVersion() (version uint64) {
	version = 1
	if v.valid() {
		version = 4
	}
	return
}

func (v XByteRange) IsZero() bool { return v.Length == 0 }

func (v XByteRange) valid() bool { return v.Length > 0 }

func (v XByteRange) encode(w io.Writer) (err error) {
	if !v.valid() {
		return errInvalidByteRange
	}

	if v.Offset > 0 {
		_, err = fmt.Fprintf(w, "%d@%d", v.Length, v.Offset)
	} else {
		_, err = fmt.Fprintf(w, "%d", v.Length)
	}
	return
}

func (v *XByteRange) decode(s string) (err error) {
	if index := strings.IndexByte(s, '@'); index > 0 {
		if v.Offset, err = strconv.ParseUint(s[index+1:], 10, 64); err != nil {
			return errInvalidByteRange
		}
		s = s[:index]
	}

	v.Length, err = strconv.ParseUint(s, 10, 64)
	if err != nil || !v.valid() {
		return errInvalidByteRange
	}

	return
}

/// ----------------------------------------------------------------------- ///

// XResolution represents a resolution ratio, that's, WidthxHeight.
type XResolution struct {
	Width  uint64 // Required. Unit: byte
	Height uint64 // Optional. Unit: byte
}

func (v XResolution) IsZero() bool {
	return v.Width == 0 || v.Height == 0
}

func (v XResolution) valid() bool { return v.Width > 0 && v.Height > 0 }

func (v XResolution) encode(w io.Writer) error {
	if !v.valid() {
		return errInvalidResolution
	}

	_, err := fmt.Fprintf(w, "%dx%d", v.Width, v.Height)
	return err
}

func (v *XResolution) decode(s string) (err error) {
	index := strings.IndexByte(s, 'x')
	if index < 0 {
		return errInvalidResolution
	}

	if v.Width, err = strconv.ParseUint(s[:index], 10, 64); err != nil {
		return errInvalidResolution
	}

	if v.Height, err = strconv.ParseUint(s[index+1:], 10, 64); err != nil {
		return errInvalidResolution
	}

	if !v.valid() {
		return errInvalidResolution
	}

	return
}

/// ----------------------------------------------------------------------- ///

const (
	XKeyMethodNone      XKeyMethod = "NONE"
	XKeyMethodAES128    XKeyMethod = "AES-128"
	XKeyMethodSampleAES XKeyMethod = "SAMPLE-AES"
)

type XKeyMethod string

func (m XKeyMethod) validate() error {
	switch m {
	case XKeyMethodNone, XKeyMethodAES128, XKeyMethodSampleAES:
		return nil
	default:
		return errInvalidKeyMethod
	}
}

// FormatIV formats the 16-octet bytes to a hexadecimal-sequence string.
func FormatIV(iv []byte, strict bool) string {
	if len(iv) != 16 && strict {
		panic(errors.New("IV is not 16-octet bytes"))
	}

	var buf strings.Builder
	buf.Grow(34)

	_ = _HexSequence(iv).encode(&buf)
	return buf.String()
}

type XKey struct {
	Method  XKeyMethod // Required
	URI     string
	IV      string
	Format  string
	Version string
}

func (x XKey) minVersion() (version uint64) {
	version = 1
	if x.Method != "" {
		if x.IV != "" && version < 2 {
			version = 2
		}
		if x.Format != "" && version < 5 {
			version = 5
		}
		if x.Version != "" && version < 5 {
			version = 5
		}
	}
	return
}

func (x XKey) IsZero() bool { return x.Method == "" }

func (x XKey) encode(w io.Writer) (err error) {
	// Method
	err = _Value(_NewAttr("METHOD", newEnum(x.Method))).encode(w)
	if err != nil || x.Method == XKeyMethodNone {
		return
	}

	if x.URI == "" {
		return errInvalidURI
	}

	var iv _HexSequence
	if x.IV != "" {
		if err = iv.decode(x.IV); err != nil {
			return fmt.Errorf("invalid IV: %w", err)
		}
	}

	err = tryWriteAttrs(w, err, false,
		_NewAttr("IV", iv),
		_NewAttr("URI", _QuotedString(x.URI)),
		_NewAttr("KEYFORMAT", _QuotedString(x.Format)),
		_NewAttr("KEYFORMATVERSIONS", _QuotedString(x.Version)),
	)

	return
}

func (x *XKey) decode(s string) (err error) {
	items := splitAttributes(s, -1)
	for _, item := range items {
		var key, value string
		if err = parseAttribute(item, &key, &value); err != nil {
			break
		}

		switch key {
		case "METHOD":
			var method _Enum[XKeyMethod]
			if err = method.decode(value); err != nil {
				err = fmt.Errorf("invalid METHOD: %w", err)
			} else {
				x.Method = method.get()
			}

		case "URI":
			var uri _QuotedString
			if err = uri.decode(value); err != nil {
				err = fmt.Errorf("invalid URI: %w", err)
			} else {
				x.URI = uri.get()
			}

		case "IV":
			var iv _HexSequence
			if err = iv.decode(value); err != nil {
				err = fmt.Errorf("invalid IV: %w", err)
			} else {
				x.IV = iv.String()
			}

		case "KEYFORMAT":
			var format _QuotedString
			if err = format.decode(value); err != nil {
				err = fmt.Errorf("invalid KEYFORMAT: %w", err)
			} else {
				x.Format = format.get()
			}

		case "KEYFORMATVERSIONS":
			var version _QuotedString
			if err = version.decode(value); err != nil {
				err = fmt.Errorf("invalid KEYFORMATVERSIONS: %w", err)
			} else {
				x.Version = version.get()
			}
		}

		if err != nil {
			return
		}
	}

	err = x._check()
	return
}

func (x *XKey) _check() (err error) {
	switch XKeyMethod(x.Method) {
	case "":
		return errors.New("missing METHOD")

	case XKeyMethodNone:
		return
	}

	if x.URI == "" {
		return errors.New("missing URI")
	}

	return
}

/// ----------------------------------------------------------------------- ///

type XMap struct {
	URI string // Required

	ByteRange XByteRange
}

func (x XMap) IsZero() bool { return x.URI == "" }

func (x XMap) valid() bool { return x.URI != "" }

func (x XMap) encode(w io.Writer) (err error) {
	return tryWriteAttrs(w, nil, true,
		_NewAttr("URI", _QuotedString(x.URI)),
		_NewAttr("BYTERANGE", x.ByteRange),
	)
}

func (x *XMap) decode(s string) (err error) {
	items := splitAttributes(s, -1)
	for _, item := range items {
		var name, value string
		if err = parseAttribute(item, &name, &value); err != nil {
			return
		}

		switch name {
		case "URI":
			var uri _QuotedString
			if err = uri.decode(value); err != nil {
				err = fmt.Errorf("invalid URI: %w", err)
			} else {
				x.URI = uri.get()
			}

		case "BYTERANGE":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid BYTERANGE: %w", err)
			} else {
				err = x.ByteRange.decode(s.get())
			}
		}

		if err != nil {
			return
		}
	}

	err = x.check()
	return
}

func (x *XMap) check() (err error) {
	if x.URI == "" {
		return errors.New("missing URI")
	}
	return
}

/// ----------------------------------------------------------------------- ///

const (
	XMediaTypeAudio          XMediaType = "AUDIO"
	XMediaTypeVideo          XMediaType = "VIDEO"
	XMediaTypeSubtitles      XMediaType = "SUBTITLES"
	XMediaTypeClosedCaptions XMediaType = "CLOSED-CAPTIONS"
)

type XMediaType string

func (t XMediaType) validate() error {
	switch t {
	case XMediaTypeAudio, XMediaTypeVideo, XMediaTypeSubtitles, XMediaTypeClosedCaptions:
		return nil
	default:
		return fmt.Errorf("invalid media type: %s", t)
	}
}

type XMedia struct {
	Type            XMediaType // Required
	Name            string     // Required
	GroupId         string     // Required
	Language        string
	AssocLanguage   string
	InstreamId      string
	Characteristics string
	Channels        string
	URI             string

	AutoSelect bool
	Default    bool
	Forced     bool
}

func (x XMedia) IsZero() bool {
	return x.Type == ""
}

func (x XMedia) encode(w io.Writer) (err error) {
	if err = x.check(); err != nil {
		return
	}

	return tryWriteAttrs(w, nil, true,
		_NewAttr("TYPE", newEnum(x.Type)),
		_NewAttr("GROUP-ID", _QuotedString(x.GroupId)),
		_NewAttr("NAME", _QuotedString(x.Name)),
		_NewAttr("LANGUAGE", _QuotedString(x.Language)),
		_NewAttr("ASSOC-LANGUAGE", _QuotedString(x.AssocLanguage)),
		_NewAttr("DEFAULT", _Bool(x.Default)),
		_NewAttr("FORCED", _Bool(x.Forced)),
		_NewAttr("AUTOSELECT", _Bool(x.AutoSelect)),
		_NewAttr("INSTREAM-ID", _QuotedString(x.InstreamId)),
		_NewAttr("CHARACTERISTICS", _QuotedString(x.Characteristics)),
		_NewAttr("CHANNELS", _QuotedString(x.Channels)),
		_NewAttr("URI", _QuotedString(x.URI)),
	)
}

func (x *XMedia) decode(s string) (err error) {
	items := splitAttributes(s, -1)
	for _, item := range items {
		var name, value string
		if err = parseAttribute(item, &name, &value); err != nil {
			return
		}

		switch name {
		case "TYPE":
			var _type _Enum[XMediaType]
			if err = _type.decode(value); err != nil {
				err = fmt.Errorf("invalid TYPE: %w", err)
			} else {
				x.Type = _type.get()
			}

		case "NAME":
			var name _QuotedString
			if err = name.decode(value); err != nil {
				err = fmt.Errorf("invalid NAME: %w", err)
			} else {
				x.Name = name.get()
			}

		case "GROUP-ID":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid GROUP-ID: %w", err)
			} else {
				x.GroupId = s.get()
			}

		case "URI":
			var uri _QuotedString
			if err = uri.decode(value); err != nil {
				err = fmt.Errorf("invalid URI: %w", err)
			} else {
				x.URI = uri.get()
			}

		case "ASSOC-LANGUAGE":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid ASSOC-LANGUAGE: %w", err)
			} else {
				x.AssocLanguage = s.get()
			}

		case "DEFAULT":
			var v _Bool
			if err = v.decode(value); err != nil {
				err = fmt.Errorf("invalid DEFAULT: %w", err)
			} else {
				x.Default = v.get()
			}

		case "FORCED":
			var v _Bool
			if err = v.decode(value); err != nil {
				err = fmt.Errorf("invalid FORCED: %w", err)
			} else {
				x.Forced = v.get()
			}

		case "AUTOSELECT":
			var v _Bool
			if err = v.decode(value); err != nil {
				err = fmt.Errorf("invalid AUTOSELECT: %w", err)
			} else {
				x.AutoSelect = v.get()
			}

		case "INSTREAM-ID":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid INSTREAM-ID: %w", err)
			} else {
				x.InstreamId = s.get()
			}

		case "CHARACTERISTICS":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid CHARACTERISTICS: %w", err)
			} else {
				x.Characteristics = s.get()
			}

		case "CHANNELS":
			var s _QuotedString
			if err = s.decode(value); err != nil {
				err = fmt.Errorf("invalid CHANNELS: %w", err)
			} else {
				x.Channels = s.get()
			}
		}

		if err != nil {
			return
		}
	}

	err = x.check()
	return
}

func (x *XMedia) check() (err error) {
	if x.Type == "" {
		return errors.New("missing TYPE")
	}
	if x.Name == "" {
		return errors.New("missing NAME")
	}
	if x.GroupId == "" {
		return errors.New("missing GROUP-ID")
	}

	if x.Type == XMediaTypeClosedCaptions {
		if x.URI != "" {
			return errors.New("CLOSED-CAPTIONS media type cannot have URI")
		}

		switch x.InstreamId {
		case "CC1", "CC2", "CC3", "CC4":
		case "SERVICE1", "SERVICE2", "SERVICE3", "SERVICE4":
		default:
			if !strings.HasPrefix(x.InstreamId, "SERVICE") {
				return errors.New("invalid INSTREAM-ID")
			}

			var v _DecimalInteger
			s := strings.TrimPrefix(x.InstreamId, "SERVICE")
			if err = v.decode(s, 1); err != nil {
				err = fmt.Errorf("invalid INSTREAM-ID: %w", err)
			} else if v > 63 {
				err = errors.New("invalid INSTREAM-ID")
			}
		}
	} else if x.InstreamId != "" {
		return errors.New("INSTREAM-ID is only valid for CLOSED-CAPTIONS media type")
	}

	return
}

func checkXMedias(medias []XMedia) (err error) {
	for _, m := range medias {
		if err = m.check(); err != nil {
			return
		}
	}

	if len(medias) < 2 {
		return
	}

	type Group struct {
		Names   map[string]struct{}
		Default bool
	}

	_mediam := make(map[string]*Group, len(medias))
	for _, m := range medias {
		group, ok := _mediam[m.GroupId]
		if !ok {
			names := make(map[string]struct{}, 2)
			names[m.Name] = struct{}{}
			_mediam[m.GroupId] = &Group{Names: names, Default: m.Default}
			continue
		}

		if _, exists := group.Names[m.Name]; exists {
			return fmt.Errorf("duplicate media name %q in group %q", m.Name, m.GroupId)
		}

		if m.Default {
			if group.Default {
				return fmt.Errorf("multiple default media in group %q", m.GroupId)
			}
			group.Default = true
		}

		group.Names[m.Name] = struct{}{}
	}

	return
}

/// ----------------------------------------------------------------------- ///

const (
	HDCPLevelNone  HDCPLevel = "NONE"
	HDCPLevelType0 HDCPLevel = "TYPE-0"
)

type HDCPLevel string

func (l HDCPLevel) validate() (err error) {
	switch l {
	case HDCPLevelNone, HDCPLevelType0:
	default:
		err = errInvalidHDCPLevel
	}
	return
}

type XStreamInf struct {
	URI string // Required

	Bandwidth        uint64 // Required. Unit: bit/s
	AverageBandwidth uint64 // Optional. Unit: bit/s

	Codecs     []string
	FrameRate  float64
	HdcpLevel  HDCPLevel
	Resolution XResolution

	Audio          string
	Video          string
	Subtitles      string
	ClosedCaptions string
}

func (x XStreamInf) IsZero() bool {
	return x.URI == ""
}

func (x XStreamInf) encode(w io.Writer) (err error) {
	if err = x.check(true); err != nil {
		return
	}

	closedCaptions := _Value(_QuotedString(x.ClosedCaptions))
	if x.ClosedCaptions == "" || x.ClosedCaptions == "NONE" {
		closedCaptions = _UnquotedString(x.ClosedCaptions)
	}

	err = tryWriteAttrs(w, nil, true,
		_NewAttr("BANDWIDTH", _DecimalInteger(x.Bandwidth)),
		_NewAttr("AVERAGE-BANDWIDTH", _DecimalInteger(x.AverageBandwidth)),
		_NewAttr("CODECS", _QuotedString(strings.Join(x.Codecs, ","))),
		_NewAttr("FRAME-RATE", _DecimalFloat(x.FrameRate)),
		_NewAttr("HDCP-LEVEL", newEnum(x.HdcpLevel)),
		_NewAttr("RESOLUTION", x.Resolution),

		_NewAttr("AUDIO", _QuotedString(x.Audio)),
		_NewAttr("VIDEO", _QuotedString(x.Video)),
		_NewAttr("SUBTITLES", _QuotedString(x.Subtitles)),
		_NewAttr("CLOSED-CAPTIONS", closedCaptions),
	)

	err = tryWriteAny(w, err, "\n", _UnquotedString(x.URI))
	return
}

func (x *XStreamInf) decode(s string) (err error) {
	err = iterAttributes(s, -1, func(name, value string) (err error) {
		switch name {
		case "BANDWIDTH":
			var v _DecimalInteger
			if err = v.decode(value, 1); err == nil {
				x.Bandwidth = v.get()
			}

		case "AVERAGE-BANDWIDTH":
			var v _DecimalInteger
			if err = v.decode(value, 0); err == nil {
				x.AverageBandwidth = v.get()
			}

		case "CODECS":
			var v _QuotedString
			if err = v.decode(value); err == nil {
				x.Codecs = strings.Split(v.get(), ",")
			}

		case "RESOLUTION":
			var v XResolution
			if err = v.decode(value); err == nil {
				x.Resolution = v
			}

		case "FRAME-RATE":
			var v _DecimalFloat
			if err = v.decode(value); err == nil {
				x.FrameRate = v.get()
			}

		case "HDCP-LEVEL":
			var v _Enum[HDCPLevel]
			if err = v.decode(value); err == nil {
				x.HdcpLevel = v.get()
			}

		case "AUDIO":
			var s _QuotedString
			if err = s.decode(value); err == nil {
				x.Audio = s.get()
			}

		case "VIDEO":
			var s _QuotedString
			if err = s.decode(value); err == nil {
				x.Video = s.get()
			}

		case "SUBTITLES":
			var s _QuotedString
			if err = s.decode(value); err == nil {
				x.Subtitles = s.get()
			}

		case "CLOSED-CAPTIONS":
			if value == "NONE" {
				x.ClosedCaptions = "NONE"
			} else {
				var s _QuotedString
				if err = s.decode(value); err == nil {
					x.ClosedCaptions = s.get()
				}
			}
		}
		return
	})

	if err == nil {
		err = x.check(false)
	}
	return
}

func (x XStreamInf) check(uri bool) (err error) {
	switch {
	case uri && x.URI == "":
		return errors.New("missing URI")
	case x.Bandwidth == 0:
		return errors.New("missing BANDWIDTH")
	}
	return
}

/// ----------------------------------------------------------------------- ///

type XIFrameStreamInf struct {
	URI string // Required

	Bandwidth        uint64 // Required. Unit: bit/s
	AverageBandwidth uint64 // Optional. Unit: bit/s

	Codecs     []string
	HdcpLevel  HDCPLevel
	Resolution XResolution

	Video string
}

func (x XIFrameStreamInf) IsZero() bool {
	return x.URI == ""
}

func (x XIFrameStreamInf) encode(w io.Writer) (err error) {
	if err = x.check(); err != nil {
		return
	}

	return tryWriteAttrs(w, nil, true,
		_NewAttr("BANDWIDTH", _DecimalInteger(x.Bandwidth)),
		_NewAttr("AVERAGE-BANDWIDTH", _DecimalInteger(x.AverageBandwidth)),
		_NewAttr("CODECS", _QuotedString(strings.Join(x.Codecs, ","))),
		_NewAttr("HDCP-LEVEL", newEnum(x.HdcpLevel)),
		_NewAttr("RESOLUTION", x.Resolution),
		_NewAttr("VIDEO", _QuotedString(x.Video)),
		_NewAttr("URI", _QuotedString(x.URI)),
	)
}

func (x *XIFrameStreamInf) decode(s string) (err error) {
	err = iterAttributes(s, -1, func(name, value string) (err error) {
		switch name {
		case "URI":
			var s _QuotedString
			if err = s.decode(value); err == nil {
				x.URI = s.get()
			}

		case "BANDWIDTH":
			var v _DecimalInteger
			if err = v.decode(value, 1); err == nil {
				x.Bandwidth = v.get()
			}

		case "AVERAGE-BANDWIDTH":
			var v _DecimalInteger
			if err = v.decode(value, 0); err == nil {
				x.AverageBandwidth = v.get()
			}

		case "CODECS":
			var v _QuotedString
			if err = v.decode(value); err == nil {
				x.Codecs = strings.Split(v.get(), ",")
			}

		case "RESOLUTION":
			var v XResolution
			if err = v.decode(value); err == nil {
				x.Resolution = v
			}

		case "HDCP-LEVEL":
			var v _Enum[HDCPLevel]
			if err = v.decode(value); err == nil {
				x.HdcpLevel = v.get()
			}

		case "VIDEO":
			var s _QuotedString
			if err = s.decode(value); err == nil {
				x.Video = s.get()
			}
		}
		return
	})

	if err == nil {
		err = x.check()
	}
	return
}

func (x XIFrameStreamInf) check() (err error) {
	switch {
	case x.URI == "":
		return errors.New("missing URI")
	case x.Bandwidth == 0:
		return errors.New("missing BANDWIDTH")
	}
	return
}

/// ----------------------------------------------------------------------- ///

type XSessionData struct {
	DataId   string
	Value    string
	URI      string
	Language string
}

func (x XSessionData) IsZero() bool {
	return x.DataId == ""
}

func (x XSessionData) encode(w io.Writer) (err error) {
	if err = x.check(); err != nil {
		return
	}

	return tryWriteAttrs(w, nil, true,
		_NewAttr("DATA-ID", _QuotedString(x.DataId)),
		_NewAttr("VALUE", _QuotedString(x.Value)),
		_NewAttr("LANGUAGE", _QuotedString(x.Language)),
		_NewAttr("URI", _QuotedString(x.URI)),
	)
}

func (x *XSessionData) decode(s string) (err error) {
	err = iterAttributes(s, -1, func(name, value string) (err error) {
		switch name {
		case "DATA-ID":
			var v _QuotedString
			if err = v.decode(s); err == nil {
				x.DataId = v.get()
			}

		case "VALUE":
			var v _QuotedString
			if err = v.decode(s); err == nil {
				x.Value = v.get()
			}

		case "LANGUAGE":
			var v _QuotedString
			if err = v.decode(s); err == nil {
				x.Language = v.get()
			}

		case "URI":
			var v _QuotedString
			if err = v.decode(s); err == nil {
				x.URI = v.get()
			}
		}
		return
	})

	if err == nil {
		err = x.check()
	}
	return
}

func (x XSessionData) check() (err error) {
	if x.DataId == "" {
		return errors.New("missing DATA-ID")
	}
	return
}

/// ----------------------------------------------------------------------- ///

type XStart struct {
	TimeOffset float64 // Required. Unit: Second
	Precise    bool
}

func (x XStart) IsZero() bool { return x.TimeOffset == 0 }

func (x XStart) encode(w io.Writer) (err error) {
	return tryWriteAttrs(w, nil, true,
		_NewAttr("TIME-OFFSET", _SignDecimalFloat(x.TimeOffset)),
		_NewAttr("PRECISE", _Bool(x.Precise)),
	)
}

func (x *XStart) decode(s string) (err error) {
	items := splitAttributes(s, -1)
	for _, item := range items {
		var name, value string
		if err = parseAttribute(item, &name, &value); err != nil {
			return
		}

		switch name {
		case "TIME-OFFSET":
			var offset _SignDecimalFloat
			if err = offset.decode(value); err != nil {
				err = fmt.Errorf("invalid TIME-OFFSET: %w", err)
			} else {
				x.TimeOffset = offset.get()
			}

		case "PRECISE":
			var name _Bool
			if err = name.decode(value); err != nil {
				err = fmt.Errorf("invalid PRECISE: %w", err)
			} else {
				x.Precise = name.get()
			}
		}

		if err != nil {
			return
		}
	}

	err = x.check()
	return
}

func (x *XStart) check() (err error) {
	if x.TimeOffset == 0 {
		return errors.New("missing TIME-OFFSET")
	}

	return
}
