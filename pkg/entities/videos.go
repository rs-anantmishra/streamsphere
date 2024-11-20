package entities

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

type Channel struct {
	Id                   int
	Name                 string
	ChannelFollowerCount string
	ChannelURL           string
	YoutubeChannelId     int
	CreatedDate          int64
}

type Playlist struct {
	Id                 int
	Title              string
	ItemCount          int
	PlaylistChannel    string
	PlaylistChannelId  string
	PlaylistUploader   string
	PlaylistUploaderId string
	ThumbnailFileId    int
	ThumbnailURL       string
	YoutubePlaylistId  string
	CreatedDate        int64
}

type Domain struct {
	Id          int
	Domain      string
	CreatedDate int64
}

type Format struct {
	Id          int
	Format      string
	FormatNote  string
	Resolution  string
	StreamType  string //Audio or Video
	CreatedDate int64
}

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

//Source - yt-dlp only
//Operations - DOWNLOAD, READ, DELETE
