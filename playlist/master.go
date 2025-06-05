package playlist

var _ PlayList = MasterPlayList{}

type MasterPlayList struct {
	Version uint64
	Start   XStart
	Medias  []XMedia
}

func (pl MasterPlayList) Type() string {
	return PlayListTypeMaster
}

func (pl MasterPlayList) MinVersion() uint64 {
	if pl.Version > 0 {
		return pl.Version
	}
	return 1
}
