package extractor

import (
	"fmt"

	"github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	"github.com/rs-anantmishra/streamsphere/utils/processor/requests"
)

type IService interface {
	ExtractSubtitlesOnly(string) bool // here we are navigating to a Video and downloading subs for it.

	/****** UPDATES *********************************************/
	ProcessRequests()
}

type service struct {
	repository IRepository
	download   IDownload
	request    requests.IRequestRepository
}

func NewDownloadService(r IRepository, rq requests.IRequestRepository, d IDownload) IService {
	return &service{
		repository: r,
		download:   d,
		request:    rq,
	}
}

func (s *service) ExtractSubtitlesOnly(videoId string) bool {
	return false
}

/****** UPDATES *********************************************/

func (s *service) ProcessRequests() {

	//read incoming requests
	requests, err := s.request.GetRequestsByRequestType(domain.RequestStatus_Recieved)
	if err != nil {
		fmt.Println("error is", err)
	}

	for i := range requests {
		//Fetch Channel Playlists
		plResult := s.download.GetChannelPlaylists(requests[i])

		//Single Video
		if len(plResult) == 1 && plResult[0].PlaylistChannelId == "" {
			//Execute Metadata and so on
			processContent(s, domain.PlaylistContentMeta{ContentId: requests[i].RequestUrl, PlaylistContentIndex: 0}, domain.Playlist{Id: -1}, requests[i])
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

func processContent(s *service, content domain.PlaylistContentMeta, playlist domain.Playlist, request domain.RequestWithStatusId) bool {

	var contentMeta domain.SavedInfo
	var domainCheck bool
	if request.Metadata == 1 {
		metadata := s.download.ExtractMetadata(content)
		if metadata.YoutubeVideoId != "" {
			domainCheck = checkContentDomain(metadata) //temporary domain restriction
		}

		if domainCheck {
			metadata = mapPlaylistMetadata(playlist, metadata, content.PlaylistContentIndex)
			contentMeta = s.repository.SaveMetadata(metadata)
		}
	}

	//pass the domain check
	if domainCheck {
		//get filename-info details from db or network
		filenameInfo := getFilenameInfo(s, contentMeta, content, request)

		if request.Thumbnail == 1 {
			thumbnail := s.download.ExtractThumbnail(filenameInfo)
			s.repository.SaveThumbnail(thumbnail)
		}

		if request.Subtitles == 1 {
			subtitles := s.download.ExtractSubtitles(filenameInfo)
			s.repository.SaveSubtitles(subtitles)
		}

		if request.Content == 1 {
			//download file
			state := s.download.ExtractMediaContent(filenameInfo)
			_ = state

			fileInfo := s.download.GetDownloadedMediaFileInfo(contentMeta, filenameInfo)
			dbResult := s.repository.SaveMediaContent(fileInfo)
			_ = dbResult
		}
	}
	return true

}

func processPlaylist(s *service, playlist domain.Playlist, contentResult domain.PlaylistContent, request domain.RequestWithStatusId) {

	//call process content - tblPVF will be updated here.
	for i := range contentResult.Content {
		playlist.ItemCount = len(contentResult.Content)
		result := processContent(s, contentResult.Content[i], playlist, request)
		fmt.Println(result)
	}
}

func processChannel(s *service, plResult []domain.Playlist, request domain.RequestWithStatusId) {
	//handle playlists in channel
	for k := range plResult {
		playlistContentInfo := domain.PlaylistContent{PlaylistId: plResult[k].YoutubePlaylistId, Content: []domain.PlaylistContentMeta{}}
		contentResult := s.download.GetPlaylistContents(playlistContentInfo)
		fmt.Println(contentResult)

		//handle playlists
		processPlaylist(s, plResult[k], contentResult, request)
	}
}

//region [helper methods]

func getPlaylistResult(s *service, request domain.RequestWithStatusId) ([]domain.Playlist, domain.PlaylistContent) {
	playlistContentInfo := domain.PlaylistContent{PlaylistId: request.RequestUrl, Content: []domain.PlaylistContentMeta{}}
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
	var result []domain.Playlist
	for i := range playlistInfoResult {
		if playlistInfoResult[i].YoutubePlaylistId == uploaderId.PlaylistId {
			//update playlist_count since its not recievable from yt-dlp
			playlistInfoResult[i].ItemCount = len(contentResult.Content)
			result = []domain.Playlist{playlistInfoResult[i]}
			break
		}
	}
	return result, contentResult
}

func mapPlaylistMetadata(playlist domain.Playlist, metadata domain.MediaInformation, contentIndex int) domain.MediaInformation {
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

func getFilenameInfo(s *service, contentMeta domain.SavedInfo, content domain.PlaylistContentMeta, request domain.RequestWithStatusId) domain.FilenameInfo {
	var filenameInfo domain.FilenameInfo
	if request.Metadata == 0 {
		filenameInfo = s.download.ExtractFilenameInfo(content.ContentId)
		filenameInfo.Id, filenameInfo.PlaylistId, _ = s.repository.GetVideoIdByContentId(content.ContentId)
	} else if request.Metadata == 1 {
		filenameInfo = domain.FilenameInfo{
			Domain:       contentMeta.MediaInfo.Domain,
			Channel:      contentMeta.MediaInfo.Channel,
			Title:        contentMeta.MediaInfo.Title,
			ContentId:    contentMeta.YoutubeVideoId,
			Id:           contentMeta.VideoId,
			PlaylistId:   contentMeta.ChannelId,
			ThumbnailUrl: contentMeta.MediaInfo.ThumbnailURL,
		}
	}
	return filenameInfo
}

//endregion
