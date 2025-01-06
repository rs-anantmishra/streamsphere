package extractor

import (
	"errors"
	"fmt"
	"time"

	// p "github.com/rs-anantmishra/streamsphere/api/presenter"
	// g "github.com/rs-anantmishra/streamsphere/pkg/global"
	"github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	e "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	r "github.com/rs-anantmishra/streamsphere/utils/processor/requests"
)

type IService interface {
	ExtractIngestMetadata() ([]e.MetadataResponse, error) // here we have an option to dl subs as well, when the metadata is available.
	ExtractIngestMedia()                                  //in case it was a metadata only files, youre free to dl video at a later time.
	ExtractSubtitlesOnly(string) bool                     // here we are navigating to a Video and downloading subs for it.

	/****** UPDATES *********************************************/
	ProcessRequests()
}

type service struct {
	repository IRepository
	download   IDownload
	request    r.IRequestRepository
}

func NewDownloadService(r IRepository, rq r.IRequestRepository, d IDownload) IService {
	return &service{
		repository: r,
		download:   d,
		request:    rq,
	}
}

// func (s *service) ExtractIngestMetadata(params e.IncomingRequest) ([]e.MetadataResponse, error) {
func (s *service) ExtractIngestMetadata() ([]e.MetadataResponse, error) {
	var response []e.MetadataResponse

	metadata, fp := s.download.ExtractMetadata()
	if len(metadata) > 0 {
		domainCheck := checkContentDomain(metadata) //temporary check placed
		if !domainCheck {
			return response, errors.New("failed: domain constraint")
		}

		lstSavedInfo := s.repository.SaveMetadata(metadata, fp)
		//error check here before continuing exec for thumbs and subs

		var thumbnails []e.Files
		var subtitles []e.Files

		thumbnails = s.download.ExtractThumbnail(fp, lstSavedInfo)
		s.repository.SaveThumbnail(thumbnails)

		// if params.SubtitlesReq {
		if true {
			subtitles = s.download.ExtractSubtitles(fp, lstSavedInfo)
			s.repository.SaveSubtitles(subtitles)
		}

		// response = createMetadataResponse(lstSavedInfo, subtitles, params.SubtitlesReq, thumbnails)
		response = createMetadataResponse(lstSavedInfo, subtitles, true, thumbnails)
		return response, nil
	}
	return response, nil
}

func (s *service) ExtractIngestMedia() {

	// defer falsifyQueueAlive()

	//cleanup of processed
	// s.download.Cleanup()

	// lstDownloads := g.NewDownloadStatus()
	// activeItem := g.NewActiveItem()

	// if len(lstDownloads) > 0 {
	if 1 > 0 {
		// for i := 0; i < len(lstDownloads); i++ {
		for i := 0; i < 5; i++ {

			//skip empties
			// if lstDownloads[i].State == 1 || lstDownloads[i].VideoURL == "" {
			if true {
				continue
			}

			//copy to active-item
			//activeItem[0] = lstDownloads[i]
			//lstDownloads[i].State = g.Downloading

			//download file
			// smi, fp, err := s.repository.GetVideoFileInfo(activeItem[0].VideoId)
			smi, fp, err := s.repository.GetVideoFileInfo(5)
			// lstDownloads[i].State = s.download.ExtractMediaContent(smi)
			state := s.download.ExtractMediaContent(smi)
			_ = state

			if err != nil {
				fmt.Println(err)
			}

			fileInfo := s.download.GetDownloadedMediaFileInfo(smi, fp)
			dbResult := s.repository.SaveMediaContent(fileInfo)

			//wait before next
			duration := time.Second
			time.Sleep(duration)

			// activeItem = clearActiveItem(activeItem)

			_ = dbResult
		}
	}
}

func (s *service) ExtractSubtitlesOnly(videoId string) bool {
	return false
}

/****** UPDATES *********************************************/

func (s *service) ProcessRequests() {

	//read incoming requests
	requests, err := s.request.GetRequestsByRequestType(e.RequestStatus_Recieved)
	if err != nil {
		fmt.Println("error is", err)
	}

	for i := range requests {

		if i != 2 {
			continue
		}

		//Fetch Channel Playlists
		plResult := s.download.GetChannelPlaylists(requests[i])

		//Single Video
		if len(plResult) == 1 && plResult[0].PlaylistChannelId == "" {
			//Execute Metadata and so on
			processContent(s, requests[i].RequestUrl)
		}
		//Single Playlist
		if plResult == nil {
			//get playlist details by uploader id
			plResult, contentResult := getPlaylistResult(s, requests[i])
			processPlaylist(s, plResult[0], contentResult)
		}
		//Channel
		if len(plResult) > 0 && plResult[0].PlaylistChannelId != "" {
			processChannel(s, plResult)
		}
	}

}

func processContent(s *service, requestUrl string) {

}

func processPlaylist(s *service, playlistResult e.Playlist, contentResult e.PlaylistContent) {

	//update tblPlaylists
	//update tblPlaylistVideoFiles

	//call process content
}

func processChannel(s *service, plResult []e.Playlist) {
	//handle playlists in channel
	for k := range plResult {
		playlistContentInfo := domain.PlaylistContent{PlaylistId: plResult[k].YoutubePlaylistId, Content: []e.PlaylistContentMeta{}}
		contentResult := s.download.GetPlaylistContents(playlistContentInfo)
		fmt.Println(contentResult)

		//handle playlists
		processPlaylist(s, plResult[k], contentResult)
	}
}

//region [helper methods]

func getPlaylistResult(s *service, request e.RequestWithStatusId) ([]e.Playlist, e.PlaylistContent) {
	playlistContentInfo := domain.PlaylistContent{PlaylistId: request.RequestUrl, Content: []e.PlaylistContentMeta{}}
	contentResult := s.download.GetPlaylistContents(playlistContentInfo)

	//pick 1st item get playlist_uploader_id
	var uploaderId domain.PlaylistUploader
	if len(contentResult.Content) > 0 {
		uploaderId = s.download.GetPlaylistUploader(request.RequestUrl)
	}
	//get all playlists details
	modifiedReq := request
	modifiedReq.RequestUrl = ytUrl + uploaderId.PlaylistUploaderId
	playlistInfoResult := s.download.GetChannelPlaylists(modifiedReq)

	//filter for this playlist
	var result []e.Playlist
	for i := range playlistInfoResult {
		if playlistInfoResult[i].YoutubePlaylistId == uploaderId.PlaylistId {
			result = []e.Playlist{playlistInfoResult[i]}
		}
	}
	return result, contentResult
}

//endregion
