package entities

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
