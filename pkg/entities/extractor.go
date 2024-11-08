package entities

type MediaInformation struct {
	YoutubeVideoId       string   `json:"id"`
	Channel              string   `json:"channel"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	Extension            string   `json:"ext"`
	Duration             int      `json:"duration"`
	Domain               string   `json:"webpage_url_domain"`
	OriginalURL          string   `json:"original_url"`
	WebpageURL           string   `json:"webpage_url"` //added later
	Tags                 []string `json:"tags"`
	Format               string   `json:"format"`
	Filesize             int      `json:"filesize_approx"`
	FormatNote           string   `json:"format_note"`
	Resolution           string   `json:"resolution"`
	Categories           []string `json:"categories"`
	ChannelId            string   `json:"channel_id"`
	ChannelURL           string   `json:"channel_url"`
	Availability         string   `json:"availability"`
	LiveStatus           string   `json:"live_status"`
	YoutubePlaylistId    string   `json:"playlist_id"`
	PlaylistTitle        string   `json:"playlist_title"`
	PlaylistCount        int      `json:"playlist_count"`
	PlaylistVideoIndex   int      `json:"playlist_index"`
	ThumbnailURL         string   `json:"thumbnail"`              //added later
	License              string   `json:"license"`                //added later
	ChannelFollowerCount int      `json:"channel_follower_count"` //added later
	UploadDate           string   `json:"upload_date"`            //added later
	ReleaseTimestamp     int64    `json:"release_timestamp"`      //added later
	ModifiedTimestamp    int64    `json:"modified_timestamp"`     //added later
	YoutubeViewCount     int      `json:"view_count"`             //added later
	LikeCount            int      `json:"like_count"`             //added later
	DislikeCount         int      `json:"dislike_count"`          //added later
	AgeLimit             int      `json:"age_limit"`              //added later
	PlayableInEmbed      Bool     `json:"playable_in_embed"`      //added later
	PlaylistChannel      string   `json:"playlist_channel"`       //added later
	PlaylistChannelId    string   `json:"playlist_channel_id"`    //added later
	PlaylistUploader     string   `json:"playlist_uploader"`      //added later
	PlaylistUploaderId   string   `json:"playlist_uploader_id"`   //added later
}

// helps unmarshalling unquoted true/false as bools in json
type Bool bool

type SavedInfo struct {
	VideoId        int
	YoutubeVideoId string
	PlaylistId     int
	ChannelId      int
	DomainId       int
	FormatId       int
	MediaInfo      MediaInformation
}

type MinimalCardsInfo struct {
	VideoId       int
	Title         string
	Description   string
	Duration      int
	WebpageURL    string
	Thumbnail     string
	VideoFilepath string
	Channel       string
}
