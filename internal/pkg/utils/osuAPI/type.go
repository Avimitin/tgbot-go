package osuAPI

type Beatmap struct {
	BeatmapsetID        string `json:"beatmapset_id"`
	BeatmapID           string `json:"beatmap_id"`
	Approved            string `json:"approved"`
	TotalLength         string `json:"total_length"`
	HitLength           string `json:"hit_length"`
	Version             string `json:"version"`
	FileMd5             string `json:"file_md5"`
	DiffSize            string `json:"diff_size"`
	DiffOverall         string `json:"diff_overall"`
	DiffApproach        string `json:"diff_approach"`
	DiffDrain           string `json:"diff_drain"`
	Mode                string `json:"mode"`
	CountNormal         string `json:"count_normal"`
	CountSlider         string `json:"count_slider"`
	CountSpinner        string `json:"count_spinner"`
	SubmitDate          string `json:"submit_date"`
	ApprovedDate        string `json:"approved_date"`
	LastUpdate          string `json:"last_update"`
	Artist              string `json:"artist"`
	ArtistUnicode       string `json:"artist_unicode"`
	Title               string `json:"title"`
	TitleUnicode        string `json:"title_unicode"`
	Creator             string `json:"creator"`
	CreatorID           string `json:"creator_id"`
	Bpm                 string `json:"bpm"`
	Source              string `json:"source"`
	Tags                string `json:"tags"`
	GenreID             string `json:"genre_id"`
	LanguageID          string `json:"language_id"`
	FavouriteCount      string `json:"favourite_count"`
	Rating              string `json:"rating"`
	Storyboard          string `json:"storyboard"`
	Video               string `json:"video"`
	DownloadUnavailable string `json:"download_unavailable"`
	AudioUnavailable    string `json:"audio_unavailable"`
	Playcount           string `json:"playcount"`
	Passcount           string `json:"passcount"`
	Packs               string `json:"packs"`
	MaxCombo            string `json:"max_combo"`
	DiffAim             string `json:"diff_aim"`
	DiffSpeed           string `json:"diff_speed"`
	Difficultyrating    string `json:"difficultyrating"`
}

type User struct {
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	JoinDate           string `json:"join_date"`
	Count300           string `json:"count300"`
	Count100           string `json:"count100"`
	Count50            string `json:"count50"`
	Playcount          string `json:"playcount"`
	RankedScore        string `json:"ranked_score"`
	TotalScore         string `json:"total_score"`
	PpRank             string `json:"pp_rank"`
	Level              string `json:"level"`
	PpRaw              string `json:"pp_raw"`
	Accuracy           string `json:"accuracy"`
	CountRankSs        string `json:"count_rank_ss"`
	CountRankSSH       string `json:"count_rank_ssh"`
	CountRankS         string `json:"count_rank_s"`
	CountRankSh        string `json:"count_rank_sh"`
	CountRankA         string `json:"count_rank_a"`
	Country            string `json:"country"`
	TotalSecondsPlayed string `json:"total_seconds_played"`
	PpCountryRank      string `json:"pp_country_rank"`
	Events             []struct {
		DisplayHTML  string `json:"display_html"`
		BeatmapID    string `json:"beatmap_id"`
		BeatmapsetID string `json:"beatmapset_id"`
		Date         string `json:"date"`
		Epicfactor   string `json:"epicfactor"`
	} `json:"events"`
}
