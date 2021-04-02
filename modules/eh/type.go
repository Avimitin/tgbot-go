package eh

type TorrentType struct {
	Hash  string `json:"hash"`
	Added string `json:"added"`
	TSize string `json:"tsize"`
	FSize string `json:"fsize"`
}

type GMetaData struct {
	Medas []*MetaData `json:"gmetadata"`
}

type MetaData struct {
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
