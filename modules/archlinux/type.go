package archlinux

import "time"

type SearchResponse struct {
	Version int  `json:"version"`
	Limit   int  `json:"limit"`
	Valid   bool `json:"valid"`
	Results []struct {
		Pkgname        string        `json:"pkgname"`
		Pkgbase        string        `json:"pkgbase"`
		Repo           string        `json:"repo"`
		Arch           string        `json:"arch"`
		Pkgver         string        `json:"pkgver"`
		Pkgrel         string        `json:"pkgrel"`
		Epoch          int           `json:"epoch"`
		Pkgdesc        string        `json:"pkgdesc"`
		URL            string        `json:"url"`
		Filename       string        `json:"filename"`
		CompressedSize int           `json:"compressed_size"`
		InstalledSize  int           `json:"installed_size"`
		BuildDate      time.Time     `json:"build_date"`
		LastUpdate     time.Time     `json:"last_update"`
		FlagDate       interface{}   `json:"flag_date"`
		Maintainers    []string      `json:"maintainers"`
		Packager       string        `json:"packager"`
		Groups         []interface{} `json:"groups"`
		Licenses       []string      `json:"licenses"`
		Conflicts      []interface{} `json:"conflicts"`
		Provides       []interface{} `json:"provides"`
		Replaces       []interface{} `json:"replaces"`
		Depends        []string      `json:"depends"`
		Optdepends     []interface{} `json:"optdepends"`
		Makedepends    []string      `json:"makedepends"`
		Checkdepends   []interface{} `json:"checkdepends"`
	} `json:"results"`
	NumPages int `json:"num_pages"`
	Page     int `json:"page"`
}

type AURSearchResponse struct {
	Version     int    `json:"version"`
	Type        string `json:"type"`
	Resultcount int    `json:"resultcount"`
	Results     []struct {
		ID             int         `json:"ID"`
		Name           string      `json:"Name"`
		PackageBaseID  int         `json:"PackageBaseID"`
		PackageBase    string      `json:"PackageBase"`
		Version        string      `json:"Version"`
		Description    string      `json:"Description"`
		URL            string      `json:"URL"`
		NumVotes       int         `json:"NumVotes"`
		Popularity     float64     `json:"Popularity"`
		OutOfDate      interface{} `json:"OutOfDate"`
		Maintainer     string      `json:"Maintainer"`
		FirstSubmitted int         `json:"FirstSubmitted"`
		LastModified   int         `json:"LastModified"`
		URLPath        string      `json:"URLPath"`
	} `json:"results"`
}

type AURInfoResponse struct {
	Version     int    `json:"version"`
	Type        string `json:"type"`
	Resultcount int    `json:"resultcount"`
	Results     []struct {
		ID             int           `json:"ID"`
		Name           string        `json:"Name"`
		PackageBaseID  int           `json:"PackageBaseID"`
		PackageBase    string        `json:"PackageBase"`
		Version        string        `json:"Version"`
		Description    string        `json:"Description"`
		URL            string        `json:"URL"`
		NumVotes       int           `json:"NumVotes"`
		Popularity     float64       `json:"Popularity"`
		OutOfDate      int64         `json:"OutOfDate"`
		Maintainer     string        `json:"Maintainer"`
		FirstSubmitted int           `json:"FirstSubmitted"`
		LastModified   int           `json:"LastModified"`
		URLPath        string        `json:"URLPath"`
		Depends        []string      `json:"Depends"`
		MakeDepends    []string      `json:"MakeDepends"`
		OptDepends     []string      `json:"OptDepends"`
		Conflicts      []string      `json:"Conflicts"`
		Provides       []string      `json:"Provides"`
		License        []string      `json:"License"`
		Keywords       []interface{} `json:"Keywords"`
	} `json:"results"`
}
