package extractor

import (
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	p "github.com/rs-anantmishra/streamsphere/api/presenter"
	e "github.com/rs-anantmishra/streamsphere/pkg/entities"
	g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

func falsifyQueueAlive() {
	qa := g.NewQueueAlive()
	qa[0] = 0
}

func clearActiveItem(activeItem []g.DownloadStatus) []g.DownloadStatus {

	//empty out active item
	activeItem[0].VideoURL = ""
	activeItem[0].VideoId = 0
	activeItem[0].StatusMessage = ""
	activeItem[0].State = 0

	return activeItem
}

func createMetadataResponse(lstSavedInfo []e.SavedInfo, subtitles []e.Files, subtitlesReq bool, thumbnails []e.Files) []p.CardsInfoResponse {
	//bind here to presenter entity
	var cardMetaDataInfoList []p.CardsInfoResponse
	const _blank string = ""

	for _, elem := range lstSavedInfo {
		var cardMetaDataInfo p.CardsInfoResponse

		cardMetaDataInfo.Channel = elem.MediaInfo.Channel
		cardMetaDataInfo.CreatedDate = int(time.Now().Unix())
		cardMetaDataInfo.Domain = elem.MediaInfo.Domain
		cardMetaDataInfo.Duration = elem.MediaInfo.Duration
		cardMetaDataInfo.IsDeleted = false
		cardMetaDataInfo.MediaURL = _blank
		cardMetaDataInfo.WebpageURL = elem.MediaInfo.WebpageURL
		cardMetaDataInfo.Playlist = elem.MediaInfo.PlaylistTitle
		cardMetaDataInfo.PlaylistVideoIndex = elem.MediaInfo.PlaylistVideoIndex
		cardMetaDataInfo.Title = elem.MediaInfo.Title
		cardMetaDataInfo.VideoFormat = elem.MediaInfo.Format
		cardMetaDataInfo.VideoId = elem.VideoId
		cardMetaDataInfo.WatchCount = 0
		cardMetaDataInfo.ViewsCount = elem.MediaInfo.YoutubeViewCount
		cardMetaDataInfo.LikesCount = elem.MediaInfo.LikeCount
		cardMetaDataInfo.FileSize = elem.MediaInfo.Filesize
		cardMetaDataInfo.UploadDate = elem.MediaInfo.UploadDate
		cardMetaDataInfo.Tags = elem.MediaInfo.Tags
		cardMetaDataInfo.Categories = elem.MediaInfo.Categories

		cardMetaDataInfoList = append(cardMetaDataInfoList, cardMetaDataInfo)
	}

	//subtitles
	if subtitlesReq {
		for i, elem := range subtitles {
			cardMetaDataInfoList[i].SubtitlesURL = elem.FilePath + elem.FileName
		}
	}

	//thumbnails // playlist thumbnails can be figured out on the UI side from Video Index
	for i := range thumbnails {
		cardMetaDataInfoList[i].Thumbnail = getImagesFromURL(thumbnails[i])
		// cardMetaDataInfoList[i].Thumbnail = thumbnails[i].FilePath + "\\" + thumbnails[i].FileName
	}

	return cardMetaDataInfoList
}

func getImagesFromURL(file e.Files) string {
	var base64EncodedImage string

	filepath := file.FilePath + string(os.PathSeparator) + file.FileName
	// Read the entire file into a byte slice
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Info(err)
	}

	switch file.Extension {
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

// can be used if want to send base64 image in place of url
func getImagesFromURLString(filepath string) string {
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
		elems = append(elems, "noimage.jpg")
		filepath = strings.Join(elems, string(os.PathSeparator))
		//filepath = `..\utils\noimage.png`
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

// Remove Characters that are not allowed in folder-names
// Only specific fields to be handled which are used for generating folder names
// Folder naming rules:
// Windows (FAT32, NTFS): Any Unicode except NUL, \, /, :, *, ?, ", <, >, |. Also, no space character at the start or end, and no period at the end.
// Mac(HFS, HFS+): Any valid Unicode except : or /
// Linux(ext[2-4]): Any byte except NUL or /
func removeForbiddenChars(metadata []e.MediaInformation) []e.MediaInformation {

	// / = 10744

	forbiddenChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	emptyString := ""
	singleSpace := " "
	doubleSpaces := "  "

	for i := 0; i < len(metadata); i++ {
		for _, elem := range forbiddenChars {

			//handle Domain
			if strings.Contains(metadata[i].Domain, elem) {
				metadata[i].Domain = strings.ReplaceAll(metadata[i].Domain, elem, emptyString)
				metadata[i].Domain = strings.TrimSpace(metadata[i].Domain)                             //Trim leading and trailing spaces
				metadata[i].Domain = strings.TrimRight(metadata[i].Domain, ".")                        //Trim trailing period
				metadata[i].Domain = strings.ReplaceAll(metadata[i].Domain, doubleSpaces, singleSpace) //Replace any double spaces that may have occurred as a result of removing characters
			}

			//handle Video Channel
			if strings.Contains(metadata[i].Channel, elem) {
				metadata[i].Channel = strings.ReplaceAll(metadata[i].Channel, elem, emptyString)
				metadata[i].Channel = strings.TrimSpace(metadata[i].Channel)
				metadata[i].Channel = strings.TrimRight(metadata[i].Channel, ".")
				metadata[i].Channel = strings.ReplaceAll(metadata[i].Channel, doubleSpaces, singleSpace)
			}

			//handle Video Title
			if strings.Contains(metadata[i].Title, elem) {
				metadata[i].Title = strings.ReplaceAll(metadata[i].Title, elem, emptyString)
				metadata[i].Title = strings.TrimSpace(metadata[i].Title)
				metadata[i].Title = strings.TrimRight(metadata[i].Title, ".")
				metadata[i].Title = strings.ReplaceAll(metadata[i].Title, doubleSpaces, singleSpace)
			}

			if strings.Contains(metadata[i].PlaylistTitle, elem) {
				metadata[i].PlaylistTitle = strings.ReplaceAll(metadata[i].PlaylistTitle, elem, emptyString)
				metadata[i].PlaylistTitle = strings.TrimSpace(metadata[i].PlaylistTitle)
				metadata[i].PlaylistTitle = strings.TrimRight(metadata[i].PlaylistTitle, ".")
				metadata[i].PlaylistTitle = strings.ReplaceAll(metadata[i].PlaylistTitle, doubleSpaces, singleSpace)
			}
		}
	}

	return metadata
}

func getFilepaths(playlistId int, fPath e.Filepath, pathType int) string {
	var fp string

	if playlistId < 0 {
		fp = GetVideoFilepath(fPath, pathType)
	} else if playlistId > 0 {
		fp = GetPlaylistFilepath(fPath, pathType)
	}

	return fp
}

func cleanDirectoryStructureFields(mediaInfo []e.MediaInformation) []e.MediaInformation {

	for k := 0; k < len(mediaInfo); k++ {
		mediaInfo[k].Domain = strings.TrimSpace(mediaInfo[k].Domain)
		mediaInfo[k].Channel = strings.TrimSpace(mediaInfo[k].Channel)
	}

	return mediaInfo
}

// temporarily placed to only accept yt
func checkContentDomain(meta []e.MediaInformation) bool {
	for _, elem := range meta {
		if strings.TrimSpace(elem.Domain) != "youtube.com" {
			return false
		}
	}
	return true
}
