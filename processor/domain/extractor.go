package domain

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

// Is it needed?
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

// #region [CATEGORIES]
type Categories struct {
	Id          int
	Name        string
	IsUsed      int
	CreatedDate int
}

type VideoFileCategories struct {
	Id          int
	CategoryId  int
	VideoId     int
	FileId      int
	CreatedDate int
}

//endregion

// #region [TAGS]
type Tags struct {
	Id          int
	Name        string
	IsUsed      int
	CreatedDate int
}

type VideoFileTags struct {
	Id          int
	TagId       int
	VideoId     int
	FileId      int
	CreatedDate int
}

//endregion

// #region [VIDEOS]

// not in use in extractor
type Videos struct {
	Id                 int
	Title              string
	Description        string
	DurationSeconds    int
	OriginalURL        string
	WebpageURL         string
	LiveStatus         string
	Availability       string
	ViewsCount         int
	LikesCount         int
	DislikeCount       int
	License            string
	AgeLimit           int
	PlayableInEmbed    string
	UploadDate         string
	ReleaseTimestamp   int
	ModifiedTimestamp  int
	IsFileDownloaded   int
	Files              []Files
	Channel            Channel
	Domain             Domain
	Format             Format
	Playlist           Playlist
	Tags               []Tags
	Categories         []Categories
	WatchCount         int
	YoutubeVideoId     string
	IsDeleted          int
	CreatedDate        int64
	PlaylistVideoIndex int
	//ThumbnailFilePath string
	//VideoFilePath string
}

// not in use in extractor
type Channel struct {
	Id                   int
	Name                 string
	ChannelFollowerCount string
	ChannelURL           string
	YoutubeChannelId     int
	CreatedDate          int64
}

// not in use in extractor
type Playlist struct {
	Id                 int
	Title              string `json:"title"`
	ItemCount          int    `json:"playlist_count"`
	PlaylistChannel    string `json:"playlist_channel"`
	PlaylistChannelId  string `json:"playlist_channel_id"`
	PlaylistUploader   string `json:"playlist_uploader"`
	PlaylistUploaderId string `json:"playlist_uploader_id"`
	ThumbnailFileId    int
	ThumbnailURL       string
	YoutubePlaylistId  string `json:"id"`
	CreatedDate        int64
}

// not in use in extractor
type Domain struct {
	Id          int
	Domain      string
	CreatedDate int64
}

// not in use in extractor
type Format struct {
	Id          int
	Format      string
	FormatNote  string
	Resolution  string
	StreamType  string //Audio or Video
	CreatedDate int64
}

// not in use in extractor
type ContentSearch struct {
	VideoId int
	Channel string
	Title   string
}

// AvailabilityType Constants
const (
	Private        = iota
	PremiumOnly    = iota
	SubscriberOnly = iota
	NeedsAuth      = iota
	Unlisted       = iota
	Public         = iota
)

// LiveStatusType Constants
const (
	NotLive    = iota
	IsLive     = iota
	IsUpcoming = iota
	WasLive    = iota
	PostLive   = iota
)

// endregion

// #region [FILES]

type Files struct {
	Id           int
	VideoId      int
	PlaylistId   int
	FileType     string
	SourceId     int
	FilePath     string
	FileName     string
	Extension    string
	FileSize     int
	FileSizeUnit string
	NetworkPath  string
	IsDeleted    int
	CreatedDate  int64
}

type Filepath struct {
	Domain        string
	Channel       string
	PlaylistTitle string
}

type StorageStatus struct {
	StorageUsedDB int64 //`json:"storage_used_db"`
	StorageUsedFS int64 //`json:"storage_used_fs"`
}

// FileType Constants
const (
	Audio     = iota
	Video     = iota
	Thumbnail = iota
	Subtitles = iota
)

// SourceType Constants
const (
	Downloaded   = iota
	Uploaded     = iota
	Local        = iota
	MetadataOnly = iota
)

// endregion
