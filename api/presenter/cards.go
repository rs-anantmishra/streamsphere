package presenter

type CardsInfoResponse struct {
	VideoId            int      `json:"video_id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Duration           int      `json:"duration"`
	WebpageURL         string   `json:"webpage_url"`
	Channel            string   `json:"channel"`
	Domain             string   `json:"domain"`
	VideoFormat        string   `json:"video_format"`
	Extension          string   `json:"extension"`
	WatchCount         int      `json:"watch_count"`
	ViewsCount         int      `json:"views_count"`
	LikesCount         int      `json:"likes_count"`
	FileSize           int      `json:"filesize"`
	UploadDate         string   `json:"upload_date"`
	IsDeleted          bool     `json:"is_deleted"`
	Categories         []string `json:"categories"`
	Tags               []string `json:"tags"`
	Playlist           string   `json:"playlist"`
	PlaylistVideoIndex int      `json:"playlist_video_index"`
	Thumbnail          string   `json:"thumbnail"`
	MediaURL           string   `json:"media_url"`
	SubtitlesURL       string   `json:"subs_url"`
	CreatedDate        int      `json:"created_date"`
}

type DownloadStatusResponse struct {
	Message  string `json:"download"`
	VideoURL string `json:"video_url"`
}

type PlaylistsInfoResponse struct {
	PlaylistId        int    `json:"playlist_id"`
	PlaylistTitle     string `json:"playlist_title"`
	PlaylistUploader  string `json:"playlist_uploader"`
	ItemCount         int    `json:"item_count"`
	YoutubePlaylistId string `json:"yt_playlist_id"`
	Thumbnail         string `json:"thumbnail"`
}
