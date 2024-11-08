package entities

type NetworkActivityLog struct {
	Id             int
	ActivityTypeId int
	InputURL       string
	Command        string
	Result         string
	CreatedDate    int
}

type LocalActivityLog struct {
	Id             int
	ActivityTypeId int
	VideoId        int
	FileId         int
	Remarks        string
	CreatedDate    int
}
