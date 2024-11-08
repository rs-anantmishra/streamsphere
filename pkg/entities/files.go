package entities

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

//operation include - UPLOAD/INSERT, READ, DELETE, METADATA INSERT (for files existing locally)
//Sources - UI, yt-dlp, local, ffmpeg (local file thumbnails), static files - default thumbnails for audio/video files

//operations - socket connection to UI for updates
