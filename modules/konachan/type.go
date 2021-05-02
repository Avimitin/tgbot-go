package konachan

type KonachanResponse struct {
	ID                  int           `json:"id"`
	Tags                string        `json:"tags"`
	CreatedAt           int           `json:"created_at"`
	CreatorID           int           `json:"creator_id"`
	Author              string        `json:"author"`
	Change              int           `json:"change"`
	Source              string        `json:"source"`
	Score               int           `json:"score"`
	Md5                 string        `json:"md5"`
	FileSize            int           `json:"file_size"`
	FileURL             string        `json:"file_url"`
	IsShownInIndex      bool          `json:"is_shown_in_index"`
	PreviewURL          string        `json:"preview_url"`
	PreviewWidth        int           `json:"preview_width"`
	PreviewHeight       int           `json:"preview_height"`
	ActualPreviewWidth  int           `json:"actual_preview_width"`
	ActualPreviewHeight int           `json:"actual_preview_height"`
	SampleURL           string        `json:"sample_url"`
	SampleWidth         int           `json:"sample_width"`
	SampleHeight        int           `json:"sample_height"`
	SampleFileSize      int           `json:"sample_file_size"`
	JpegURL             string        `json:"jpeg_url"`
	JpegWidth           int           `json:"jpeg_width"`
	JpegHeight          int           `json:"jpeg_height"`
	JpegFileSize        int           `json:"jpeg_file_size"`
	Rating              string        `json:"rating"`
	HasChildren         bool          `json:"has_children"`
	ParentID            interface{}   `json:"parent_id"`
	Status              string        `json:"status"`
	Width               int           `json:"width"`
	Height              int           `json:"height"`
	IsHeld              bool          `json:"is_held"`
	FramesPendingString string        `json:"frames_pending_string"`
	FramesPending       []interface{} `json:"frames_pending"`
	FramesString        string        `json:"frames_string"`
	Frames              []interface{} `json:"frames"`
}
