package entities

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

// can be created through UI or through yt-dlp
// operations include - INSERT, DELETE, READ
