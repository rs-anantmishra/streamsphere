package extractor

import (
	"fmt"
	"time"

	// p "github.com/rs-anantmishra/streamsphere/api/presenter"
	// g "github.com/rs-anantmishra/streamsphere/pkg/global"
	"github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	e "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	r "github.com/rs-anantmishra/streamsphere/utils/processor/requests"
)

type IService interface {
	ExtractIngestMedia()              //in case it was a metadata only files, youre free to dl video at a later time.
	ExtractSubtitlesOnly(string) bool // here we are navigating to a Video and downloading subs for it.

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
			processContent(s, e.PlaylistContentMeta{ContentId: requests[i].RequestUrl, PlaylistContentIndex: 0}, e.Playlist{Id: -1}, requests[i])
		}
		//Single Playlist
		if plResult == nil {
			//get playlist details by uploader id
			plResult, contentResult := getPlaylistResult(s, requests[i])
			processPlaylist(s, plResult[0], contentResult, requests[i])
		}
		//Channel
		if len(plResult) > 0 && plResult[0].PlaylistChannelId != "" {
			processChannel(s, plResult, requests[i])
		}
	}

}

func processContent(s *service, content e.PlaylistContentMeta, playlist e.Playlist, request e.RequestWithStatusId) bool {

	var contentMeta e.MediaInformation
	if request.Metadata == 1 {
		metadata := s.download.ExtractMetadata(content)
		if len(metadata) > 0 {
			domainCheck := checkContentDomain(metadata) //temporary check placed
			if !domainCheck {
				// log: errors.New("failed: domain constraint")
				return false
			}
		}

		//update playlist attributes
		metadata[0] = mapPlaylistMetadata(playlist, metadata[0], content.PlaylistContentIndex)
		contentMeta = metadata[0]
		//save metadata
		lstSavedInfo := s.repository.SaveMetadata(metadata)
		fmt.Println(lstSavedInfo)
	}

	//savedInfo - Channel, Domain, Title, (PlaylistId - won't be needed anymore)
	if request.Thumbnail == 1 {

		//filenames to be formed require this info
		var filenameInfo e.FilenameInfo
		if request.Metadata == 0 {
			filenameInfo = s.download.ExtractFilenameInfo(content.ContentId)
		} else if request.Metadata == 1 {
			filenameInfo = e.FilenameInfo{Domain: contentMeta.Domain, Channel: contentMeta.Channel, Title: contentMeta.Title}
		}

		thumbnails := s.download.ExtractThumbnail(filenameInfo, content.ContentId)
		s.repository.SaveThumbnail(thumbnails)
	}

	if request.Subtitles == 1 {
		subtitles := s.download.ExtractSubtitles(e.Filepath{}, []e.SavedInfo{})
		s.repository.SaveSubtitles(subtitles)

	}

	if request.Content == 1 {
		//place content downloader here.
	}
	// response = createMetadataResponse(lstSavedInfo, subtitles, true, thumbnails)
	// return response, nil
	return true

}

func processPlaylist(s *service, playlist e.Playlist, contentResult e.PlaylistContent, request e.RequestWithStatusId) {

	//update tblPlaylists
	// playlist.Id = s.repository.SavePlaylist(playlist)

	//call process content - tblPVF will be updated here.
	for i := range contentResult.Content {
		result := processContent(s, contentResult.Content[i], playlist, request)
		fmt.Println(result)
	}
}

func processChannel(s *service, plResult []e.Playlist, request e.RequestWithStatusId) {
	//handle playlists in channel
	for k := range plResult {
		playlistContentInfo := domain.PlaylistContent{PlaylistId: plResult[k].YoutubePlaylistId, Content: []e.PlaylistContentMeta{}}
		contentResult := s.download.GetPlaylistContents(playlistContentInfo)
		fmt.Println(contentResult)

		//handle playlists
		processPlaylist(s, plResult[k], contentResult, request)
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
			//update playlist_count since its not recievable from yt-dlp
			playlistInfoResult[i].ItemCount = len(contentResult.Content)
			result = []e.Playlist{playlistInfoResult[i]}
			break
		}
	}
	return result, contentResult
}

func mapPlaylistMetadata(playlist e.Playlist, metadata e.MediaInformation, contentIndex int) e.MediaInformation {
	//update playlist attributes
	metadata.PlaylistTitle = playlist.Title
	metadata.PlaylistCount = playlist.ItemCount
	metadata.PlaylistChannel = playlist.PlaylistChannel
	metadata.PlaylistChannelId = playlist.PlaylistChannelId
	metadata.PlaylistUploader = playlist.PlaylistUploader
	metadata.PlaylistUploaderId = playlist.PlaylistUploaderId
	metadata.YoutubePlaylistId = playlist.YoutubePlaylistId
	metadata.PlaylistVideoIndex = contentIndex

	return metadata
}

//endregion
