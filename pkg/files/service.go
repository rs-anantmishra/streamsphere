package files

import (
	"fmt"

	"github.com/rs-anantmishra/streamsphere/api/presenter"
)

type IService interface {
	StorageStatusInfo() (presenter.StorageStatus, error)
}

type service struct {
	repository IRepository
}

func NewFilesService(r IRepository) IService {
	return &service{
		repository: r,
	}
}

func (s *service) StorageStatusInfo() (presenter.StorageStatus, error) {

	var storageStatus presenter.StorageStatus
	dbStorageStatus, err := s.repository.DBStorageStatusInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dbStorageStatus)
	storageStatus.StorageUsedDB = dbStorageStatus

	mediaDir := GetMediaDirectory()
	fsStorageStatus, err := FSStorageStatusInfo(mediaDir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsStorageStatus)
	storageStatus.StorageUsedFS = fsStorageStatus

	return storageStatus, nil
}
