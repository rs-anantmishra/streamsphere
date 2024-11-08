package entities

// Field names should start with an uppercase letter
type IncomingRequest struct {
	Indicator    string `json:"Indicator"`
	SubtitlesReq bool   `json:"SubtitlesReq"`
	IsAudioOnly  bool   `json:"IsAudioOnly"`
}

type QueueDownloads struct {
	DownloadVideos []DownloadMedia `json:"DownloadMedia"` //Since metadata will always download first, then UI will send local VideoId and VideoURL
}

type DownloadMedia struct {
	VideoId         int    `json:"VideoId"`
	VideoURL        string `json:"VideoURL"`
	IsPlaylistVideo bool   `json:"IsPlaylistVideo"`
}
