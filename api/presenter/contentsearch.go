package presenter

type ContentSearchResponse struct {
	VideoId int    `json:"video_id"`
	Channel string `json:"channel"`
	Title   string `json:"title"`
}
