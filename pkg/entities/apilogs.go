package entities

type APILogs struct {
	Id             int
	APIName        string
	ExecutionStart int
	ExecutionEnd   int
	APIInputs      string
	APIResult      string
	ActivityLogId  int
	VideoId        int
	CreatedDate    int
}
