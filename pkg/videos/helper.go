package videos

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/rs-anantmishra/streamsphere/api/presenter"
	"github.com/rs-anantmishra/streamsphere/config"
	_ "github.com/rs-anantmishra/streamsphere/config"
	"github.com/rs-anantmishra/streamsphere/pkg/entities"
)

func getVideosPageInfo(videos []entities.Videos) []presenter.CardsInfoResponse {

	var lstCardsInfo []presenter.CardsInfoResponse
	for _, elem := range videos {
		var cardsInfo presenter.CardsInfoResponse

		cardsInfo.VideoId = elem.Id
		cardsInfo.Title = elem.Title
		cardsInfo.Description = elem.Description
		cardsInfo.Duration = elem.DurationSeconds
		cardsInfo.WebpageURL = elem.WebpageURL
		cardsInfo.Channel = elem.Channel.Name
		cardsInfo.Domain = elem.Domain.Domain
		cardsInfo.VideoFormat = elem.Format.Format
		cardsInfo.WatchCount = elem.WatchCount
		cardsInfo.ViewsCount = elem.ViewsCount
		cardsInfo.LikesCount = elem.LikesCount
		cardsInfo.UploadDate = elem.UploadDate
		cardsInfo.PlaylistVideoIndex = elem.PlaylistVideoIndex

		//tags and categories - name only
		cardsInfo.Tags, cardsInfo.Categories = getTagsCategories(elem.Tags, elem.Categories)

		//files - limited info only
		cardsInfo.FileSize, cardsInfo.MediaURL, cardsInfo.Thumbnail, cardsInfo.Extension = getFilesInfo(elem.Files)

		//additional transforms
		cardsInfo.MediaURL = urlTransforms(cardsInfo.MediaURL)
		lstCardsInfo = append(lstCardsInfo, cardsInfo)
	}

	return lstCardsInfo
}

func getFilesInfo(files []entities.Files) (int, string, string, string) {

	filesize := 0
	contentFilepath := ``
	thumbnail := ``
	extension := ``

	for idx := range files {
		if files[idx].FileType == "Video" {
			filesize = files[idx].FileSize
			contentFilepath = files[idx].FilePath + string(os.PathSeparator) + files[idx].FileName
			extension = files[idx].Extension
		} else if files[idx].FileType == "Thumbnail" {
			thumbnailURL := files[idx].FilePath + string(os.PathSeparator) + files[idx].FileName
			thumbnail = thumbnailURL

			thumbnail = urlTransforms(thumbnailURL)
		}
	}

	return filesize, contentFilepath, thumbnail, extension
}

func getTagsCategories(tags []entities.Tags, categories []entities.Categories) ([]string, []string) {
	var resultTags []string
	var resultCategories []string

	for _, elem := range tags {
		resultTags = append(resultTags, elem.Name)
	}

	for _, elem := range categories {
		resultCategories = append(resultCategories, elem.Name)
	}

	return resultTags, resultCategories
}

func getImagesFromURL(filepath string) string {
	var base64EncodedImage string
	splitter := "."

	// Read the entire file into a byte slice
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Info(err)
	}

	if len(bytes) == 0 {
		var elems []string
		elems = append(elems, "..")
		elems = append(elems, "utils")
		elems = append(elems, "noimage.png")
		filepath = strings.Join(elems, string(os.PathSeparator))

		// filepath = `..\utils\noimage.png`
		bytes, err = os.ReadFile(filepath)
		if err != nil {
			log.Info(err)
		}
	}

	splits := strings.SplitN(filepath, splitter, -1)
	extension := splits[len(splits)-1]

	switch extension {
	case "jpeg":
		base64EncodedImage += "data:image/jpeg;base64,"
	case "jpg":
		base64EncodedImage += "data:image/jpg;base64,"
	case "png":
		base64EncodedImage += "data:image/png;base64,"
	case "webp":
		base64EncodedImage += "data:image/webp;base64,"
	}

	base64EncodedImage += base64.StdEncoding.EncodeToString(bytes)
	return base64EncodedImage
}

func getContentSearchResponse(list []entities.ContentSearch) []presenter.ContentSearchResponse {
	var result []presenter.ContentSearchResponse
	for idx := range list {
		result = append(result, presenter.ContentSearchResponse{
			VideoId: list[idx].VideoId,
			Title:   list[idx].Title,
			Channel: list[idx].Channel,
		})
	}
	return result
}

func getPlaylistsPageInfo(playlists []entities.Playlist) []presenter.PlaylistsInfoResponse {

	var lstPlaylistsInfo []presenter.PlaylistsInfoResponse
	for _, elem := range playlists {
		var playlistInfo presenter.PlaylistsInfoResponse

		playlistInfo.PlaylistId = elem.Id
		playlistInfo.PlaylistTitle = elem.Title
		playlistInfo.PlaylistUploader = elem.PlaylistUploader
		playlistInfo.ItemCount = elem.ItemCount
		playlistInfo.YoutubePlaylistId = elem.YoutubePlaylistId
		playlistInfo.Thumbnail = urlTransforms(elem.ThumbnailURL)

		lstPlaylistsInfo = append(lstPlaylistsInfo, playlistInfo)
	}

	return lstPlaylistsInfo
}

func urlTransforms(url string) string {
	filesHost := config.Config("FILE_HOSTING", false)
	defaultFilesPath := config.Config("MEDIA_PATH", true)

	//change media directory to public url
	url = strings.ReplaceAll(url, defaultFilesPath, filesHost)

	//direction of slashes for url formation
	url = strings.ReplaceAll(url, "\\", "/")

	//these chars will not be handled by webserver in file names
	//chars: # = %23
	url = strings.ReplaceAll(url, "#", "%23")
	return url
}
