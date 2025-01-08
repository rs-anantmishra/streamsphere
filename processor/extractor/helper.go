package extractor

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs-anantmishra/streamsphere/utils/processor/config"
	e "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
)

func createMetadataResponse(lstSavedInfo []e.SavedInfo, subtitles []e.Files, subtitlesReq bool, thumbnails []e.Files) []e.MetadataResponse {
	//bind here to presenter entity
	var cardMetaDataInfoList []e.MetadataResponse
	const _blank string = ""

	for _, elem := range lstSavedInfo {
		var cardMetaDataInfo e.MetadataResponse

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
		//cardMetaDataInfoList[i].Thumbnail = getImagesFromURL(thumbnails[i])
		cardMetaDataInfoList[i].Thumbnail = urlTransforms(thumbnails[i].FilePath + string(os.PathSeparator) + thumbnails[i].FileName)
	}

	return cardMetaDataInfoList
}

func getImagesFromURL(file e.Files) string {
	var base64EncodedImage string

	filepath := file.FilePath + string(os.PathSeparator) + file.FileName
	// Read the entire file into a byte slice
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
			fmt.Println(err)
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
func removeForbiddenChars(metadata e.MediaInformation) e.MediaInformation {

	// / = 10744

	forbiddenChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	emptyString := ""
	singleSpace := " "
	doubleSpaces := "  "

	// for i := 0; i < len(metadata); i++ {
	for _, elem := range forbiddenChars {

		//handle Domain
		if strings.Contains(metadata.Domain, elem) {
			metadata.Domain = strings.ReplaceAll(metadata.Domain, elem, emptyString)
			metadata.Domain = strings.TrimSpace(metadata.Domain)                             //Trim leading and trailing spaces
			metadata.Domain = strings.TrimRight(metadata.Domain, ".")                        //Trim trailing period
			metadata.Domain = strings.ReplaceAll(metadata.Domain, doubleSpaces, singleSpace) //Replace any double spaces that may have occurred as a result of removing characters
		}

		//handle Video Channel
		if strings.Contains(metadata.Channel, elem) {
			metadata.Channel = strings.ReplaceAll(metadata.Channel, elem, emptyString)
			metadata.Channel = strings.TrimSpace(metadata.Channel)
			metadata.Channel = strings.TrimRight(metadata.Channel, ".")
			metadata.Channel = strings.ReplaceAll(metadata.Channel, doubleSpaces, singleSpace)
		}

		//handle Video Title
		if strings.Contains(metadata.Title, elem) {
			metadata.Title = strings.ReplaceAll(metadata.Title, elem, emptyString)
			metadata.Title = strings.TrimSpace(metadata.Title)
			metadata.Title = strings.TrimRight(metadata.Title, ".")
			metadata.Title = strings.ReplaceAll(metadata.Title, doubleSpaces, singleSpace)
		}

		if strings.Contains(metadata.PlaylistTitle, elem) {
			metadata.PlaylistTitle = strings.ReplaceAll(metadata.PlaylistTitle, elem, emptyString)
			metadata.PlaylistTitle = strings.TrimSpace(metadata.PlaylistTitle)
			metadata.PlaylistTitle = strings.TrimRight(metadata.PlaylistTitle, ".")
			metadata.PlaylistTitle = strings.ReplaceAll(metadata.PlaylistTitle, doubleSpaces, singleSpace)
		}
	}
	// }

	return metadata
}

func removeForbiddenCharsGeneric(input string) string {
	forbiddenChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	emptyString := ""
	singleSpace := " "
	doubleSpaces := "  "

	for _, elem := range forbiddenChars {
		input = strings.ReplaceAll(input, elem, emptyString)
		input = strings.TrimSpace(input)
		input = strings.TrimRight(input, ".")
		input = strings.ReplaceAll(input, doubleSpaces, singleSpace)
	}

	return input
}

func getFilepaths(fPath e.Filepath, pathType int) string {
	fp := GetVideoFilepath(fPath, pathType)
	return fp
}

func cleanDirectoryStructureFields(mediaInfo e.MediaInformation) e.MediaInformation {

	// for k := 0; k < len(mediaInfo); k++ {
	mediaInfo.Domain = strings.TrimSpace(mediaInfo.Domain)
	mediaInfo.Channel = strings.TrimSpace(mediaInfo.Channel)
	// }

	return mediaInfo
}

// temporarily placed to only accept yt
func checkContentDomain(meta e.MediaInformation) bool {
	return strings.TrimSpace(meta.Domain) == "youtube.com"
}

func urlTransforms(url string) string {
	filesHost := config.Config("FILE_HOSTING", false)
	defaultFilesPath := config.Config("MEDIA_PATH", true)

	//change media directory to public url
	url = strings.ReplaceAll(url, defaultFilesPath, filesHost)

	//direction of slashes for url formation
	url = strings.ReplaceAll(url, "\\", "/")

	//these chars will not be handled by webserver in file names
	//chars: # = %23 //chars: % = %25 //handling this on ui
	//url = strings.ReplaceAll(url, "#", "%23")
	//url = strings.ReplaceAll(url, "%", "%25")
	return url
}
