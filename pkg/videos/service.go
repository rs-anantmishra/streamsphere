package videos

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs-anantmishra/streamsphere/api/presenter"
)

// Service is an interface from which our api module can access our repository of all our models
type IService interface {
	GetVideos() ([]presenter.CardsInfoResponse, error)
	GetContentById(int) ([]presenter.CardsInfoResponse, error)
	DeleteContentById(int) (bool, error)
	GetPlaylistVideos(int) ([]presenter.CardsInfoResponse, error)
	GetPlaylists() ([]presenter.PlaylistsInfoResponse, error)
	GetVideoSearchData() ([]presenter.ContentSearchResponse, error)
	//TODO:
	//InsertVideo(videos *entities.Videos) (*entities.Videos, error)
	//UpdateVideo(videos *entities.Videos) (*entities.Videos, error)
	//RemoveVideo(ID string) error
}

type service struct {
	repository IRepository
}

func NewVideoService(r IRepository) IService {
	return &service{
		repository: r,
	}
}

// GetVideos implements IService.
func (s *service) GetVideos() ([]presenter.CardsInfoResponse, error) {
	allVideos, err := s.repository.GetAllVideos()
	if err != nil {
		return nil, err
	}
	result := getVideosPageInfo(allVideos)

	return result, nil
}

// GetVideos implements IService.
func (s *service) GetPlaylists() ([]presenter.PlaylistsInfoResponse, error) {
	allPlaylists, err := s.repository.GetAllPlaylists()
	if err != nil {
		return nil, err
	}
	result := getPlaylistsPageInfo(allPlaylists)
	return result, nil
}

func (s *service) GetPlaylistVideos(playlistId int) ([]presenter.CardsInfoResponse, error) {
	allVideos, err := s.repository.GetPlaylistVideos(playlistId)
	if err != nil {
		return nil, err
	}
	result := getVideosPageInfo(allVideos)

	return result, nil
}

// GetVideos implements IService.
func (s *service) GetVideoSearchData() ([]presenter.ContentSearchResponse, error) {
	searchInfo, err := s.repository.GetVideoSearchInfo()
	if err != nil {
		return nil, err
	}

	result := getContentSearchResponse(searchInfo)
	return result, nil
}

func (s *service) GetContentById(contentId int) ([]presenter.CardsInfoResponse, error) {
	content, err := s.repository.GetContentById(contentId)
	if err != nil {
		return nil, err
	}
	result := getVideosPageInfo(content)

	return result, nil
}

func (s *service) DeleteContentById(contentId int) (bool, error) {
	isDeleteSuccessful, contentFilePath, thumbnailFilePath, err := s.repository.DeleteContentById(contentId)

	var result bool
	//if delete is successful in db remove files from filessystem
	if isDeleteSuccessful {
		//remove content file
		err := os.Remove(contentFilePath)
		if err != nil {
			return result, errors.New(fmt.Sprintln("Db cleared. Removal of file: %v failed.", contentFilePath))
		}
		//remove thumbnail file
		err = os.Remove(thumbnailFilePath)
		if err != nil {
			return result, errors.New(fmt.Sprintln("Db cleared. Removal of file: %v failed.", thumbnailFilePath))
		}
		result = true
	} else {
		return isDeleteSuccessful, err
	}

	return result, nil
}
