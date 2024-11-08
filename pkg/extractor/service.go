package extractor

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	p "github.com/rs-anantmishra/streamsphere/api/presenter"
	e "github.com/rs-anantmishra/streamsphere/pkg/entities"
	g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

type IService interface {
	ExtractIngestMetadata(p e.IncomingRequest) ([]p.CardsInfoResponse, error) // here we have an option to dl subs as well, when the metadata is available.
	ExtractIngestMedia()                                                      //in case it was a metadata only files, youre free to dl video at a later time.
	ExtractSubtitlesOnly(string) bool                                         // here we are navigating to a Video and downloading subs for it.
}

type service struct {
	repository IRepository
	download   IDownload
}

func NewDownloadService(r IRepository, d IDownload) IService {
	return &service{
		repository: r,
		download:   d,
	}
}

func (s *service) ExtractIngestMetadata(params e.IncomingRequest) ([]p.CardsInfoResponse, error) {
	var response []p.CardsInfoResponse

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

		if params.SubtitlesReq {
			subtitles = s.download.ExtractSubtitles(fp, lstSavedInfo)
			s.repository.SaveSubtitles(subtitles)
		}

		response = createMetadataResponse(lstSavedInfo, subtitles, params.SubtitlesReq, thumbnails)
		return response, nil
	}
	return response, nil
}

func (s *service) ExtractIngestMedia() {

	defer falsifyQueueAlive()

	//cleanup of processed
	s.download.Cleanup()

	lstDownloads := g.NewDownloadStatus()
	activeItem := g.NewActiveItem()

	if len(lstDownloads) > 0 {
		for i := 0; i < len(lstDownloads); i++ {

			//skip empties
			if lstDownloads[i].State == g.Completed || lstDownloads[i].VideoURL == "" {
				continue
			}

			//copy to active-item
			activeItem[0] = lstDownloads[i]
			lstDownloads[i].State = g.Downloading

			//download file
			smi, fp, err := s.repository.GetVideoFileInfo(activeItem[0].VideoId)
			lstDownloads[i].State = s.download.ExtractMediaContent(smi)

			if err != nil {
				log.Info(err)
			}

			fileInfo := s.download.GetDownloadedMediaFileInfo(smi, fp)
			dbResult := s.repository.SaveMediaContent(fileInfo)

			//wait before next
			duration := time.Second
			time.Sleep(duration)

			activeItem = clearActiveItem(activeItem)

			_ = dbResult
		}
	}
}

func (s *service) ExtractSubtitlesOnly(videoId string) bool {
	return false
}
