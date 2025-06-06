package playlist

// Define PlayList Tags
const (
	// Basic Tags
	EXTM3U        Tag = "#EXTM3U"        // RFC 8216, 4.3.1.1
	EXT_X_VERSION Tag = "#EXT-X-VERSION" // RFC 8216, 4.3.1.2

	// Media Segment Tags
	EXTINF                  Tag = "#EXTINF"                  // RFC 8216, 4.3.2.1
	EXT_X_BYTERANGE         Tag = "#EXT-X-BYTERANGE"         // RFC 8216, 4.3.2.2
	EXT_X_DISCONTINUITY     Tag = "#EXT-X-DISCONTINUITY"     // RFC 8216, 4.3.2.3
	EXT_X_KEY               Tag = "#EXT-X-KEY"               // RFC 8216, 4.3.2.4
	EXT_X_MAP               Tag = "#EXT-X-MAP"               // RFC 8216, 4.3.2.5
	EXT_X_PROGRAM_DATE_TIME Tag = "#EXT-X-PROGRAM-DATE-TIME" // RFC 8216, 4.3.2.6
	EXT_X_DATERANGE         Tag = "#EXT-X-DATERANGE"         // RFC 8216, 4.3.2.7

	// Media Playlist Tags
	EXT_X_TARGETDURATION         Tag = "#EXT-X-TARGETDURATION"         // RFC 8216, 4.3.3.1
	EXT_X_MEDIA_SEQUENCE         Tag = "#EXT-X-MEDIA-SEQUENCE"         // RFC 8216, 4.3.3.2
	EXT_X_DISCONTINUITY_SEQUENCE Tag = "#EXT-X-DISCONTINUITY-SEQUENCE" // RFC 8216, 4.3.3.3
	EXT_X_ENDLIST                Tag = "#EXT-X-ENDLIST"                // RFC 8216, 4.3.3.4
	EXT_X_PLAYLIST_TYPE          Tag = "#EXT-X-PLAYLIST-TYPE"          // RFC 8216, 4.3.3.5
	EXT_X_I_FRAMES_ONLY          Tag = "#EXT-X-I-FRAMES-ONLY"          // RFC 8216, 4.3.3.6

	// Master Playlist Tags
	EXT_X_MEDIA              Tag = "#EXT-X-MEDIA"              // RFC 8216, 4.3.4.1
	EXT_X_STREAM_INF         Tag = "#EXT-X-STREAM-INF"         // RFC 8216, 4.3.4.2
	EXT_X_I_FRAME_STREAM_INF Tag = "#EXT-X-I-FRAME-STREAM-INF" // RFC 8216, 4.3.4.3
	EXT_X_SESSION_DATA       Tag = "#EXT-X-SESSION-DATA"       // RFC 8216, 4.3.4.4
	EXT_X_SESSION_KEY        Tag = "#EXT-X-SESSION-KEY"        // RFC 8216, 4.3.4.5

	// Media or Master Playlist Tags
	EXT_X_INDEPENDENT_SEGMENTS Tag = "#EXT-X-INDEPENDENT-SEGMENTS" // RFC 8216, 4.3.5.1
	EXT_X_START                Tag = "#EXT-X-START"                // RFC 8216, 4.3.5.2
)

// Tag is used to define a playlist tag.
type Tag string
