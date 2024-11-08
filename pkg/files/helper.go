package files

import (
	"os"
	"path/filepath"

	"github.com/rs-anantmishra/streamsphere/config"
)

// get media directory size from fs
func FSStorageStatusInfo(path string) (int64, error) {
	var size int64
	adjSize := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	}
	err := filepath.Walk(path, adjSize)
	return size, err
}

func GetMediaDirectory() string {
	return config.Config("MEDIA_PATH", true)
}
