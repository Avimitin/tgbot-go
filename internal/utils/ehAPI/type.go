package ehAPI

import "errors"

var (
	CantParseNilDataErr = errors.New("can't parse empty byte data")
	NoneUrlErr          = errors.New("invalid url input")
	NilGidErr           = errors.New("invalid gid input")
	WrongUrlErr         = errors.New("invalid url for fetching gid")
)

type GIDListType [][]string

type GRequest struct {
	Method    string      `json:"method"`
	GIDList   GIDListType `json:"gidlist"`
	Namespace int32       `json:"namespace"`
}

func NewGRequest(gl GIDListType, ns int32) *GRequest {
	if ns != 1 {
		return &GRequest{Method: "gdata", GIDList: gl}
	}
	return &GRequest{Method: "gdata", GIDList: gl, Namespace: 1}
}

type TorrentType struct {
	Hash  string `json:"hash"`
	Added string `json:"added"`
	TSize string `json:"tsize"`
	FSize string `json:"fsize"`
}

type GMetaData struct {
	GMD []*GMetaDataType `json:"gmetadata"`
}

type GMetaDataType struct {
	Gid          int64          `json:"gid"`
	Token        string         `json:"token"`
	ArchiveKey   string         `json:"archive_key"`
	Title        string         `json:"title"`
	TitleJpn     string         `json:"title_jpn"`
	Category     string         `json:"category"`
	Thumb        string         `json:"thumb"`
	Uploader     string         `json:"uploader"`
	Posted       string         `json:"posted"`
	FileCount    string         `json:"filecount"`
	Filesize     int64          `json:"filesize"`
	Expunged     bool           `json:"expunged"`
	Rating       string         `json:"rating"`
	TorrentCount string         `json:"torrentcount"`
	Torrents     []*TorrentType `json:"torrents"`
	Tags         []string       `json:"tags"`
	Error        string         `json:"error"`
}
