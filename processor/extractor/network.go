package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	c "github.com/rs-anantmishra/streamsphere/utils/processor/config"
	"github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	e "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
	// g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

type IDownload interface {
	ExtractMetadata(input e.PlaylistContentMeta) e.MediaInformation
	ExtractMediaContent(filenameInfo e.FilenameInfo) int
	ExtractSubtitles(filenameInfo e.FilenameInfo) e.Files
	ExtractThumbnail(lstSavedInfo e.FilenameInfo) e.Files
	GetDownloadedMediaFileInfo(smi e.SavedInfo, fp e.FilenameInfo) e.Files
	// Cleanup()
	GetChannelPlaylists(request domain.Request) []domain.Playlist
	GetPlaylistContents(request domain.PlaylistContent) domain.PlaylistContent
	GetPlaylistUploader(playlistUrl string) domain.PlaylistUploader
	ExtractFilenameInfo(contentId string) e.FilenameInfo
}

type download struct {
	// p e.IncomingRequest
	// lstDownloads []g.DownloadStatus
}

func NewDownload() IDownload {
	return &download{}
}

// Queue cleanup not needed.
// func (d *download) Cleanup() {
// 	// var updated []g.DownloadStatus
// 	for i := 0; i < len(d.lstDownloads); i++ {
// 		ds := d.lstDownloads
// 		if ds[i].State != g.Completed {
// 			updated = append(updated, ds[i])
// 		}
// 	}
// 	d.lstDownloads = updated
// }

func (d *download) GetChannelPlaylists(request domain.Request) []domain.Playlist {

	args, command, itemsCount := cmdBuilderChannelPlaylists(request.RequestUrl)
	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "GetChannelPlaylists - Cmd.Start")

	pResult := executeProcess(stdout, false)
	playlistInfo := parsePlaylistInfo(pResult, itemsCount)

	// helper method fix
	for i := range playlistInfo {
		playlistInfo[i].Title = removeForbiddenCharsGeneric(playlistInfo[i].Title)
	}

	return playlistInfo
}

func (d *download) GetPlaylistContents(contentInfo domain.PlaylistContent) domain.PlaylistContent {

	args, command, itemsCount := cmdBuilderPlaylistContents(contentInfo.PlaylistId)
	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "GetPlaylistContents - Cmd.Start")

	pResult := executeProcess(stdout, false)
	contentInfo.Content = parseContentInfo(pResult, itemsCount)

	return contentInfo
}

func (d *download) GetPlaylistUploader(playlistUrl string) domain.PlaylistUploader {

	args, command, itemsCount := cmdBuilderGetPlaylistDetails(playlistUrl)
	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "GetPlaylistDetails - Cmd.Start")

	pResult := executeProcess(stdout, false)
	uploaderInfo := parseUploaderInfo(pResult, itemsCount)

	fmt.Println(uploaderInfo)

	return uploaderInfo
}

// #region [public methods]

func (d *download) ExtractFilenameInfo(contentId string) e.FilenameInfo {

	// args, command, totalItems := cmdBuilderMetadata(d.p.Indicator)
	videoUrl := contentId
	args, command, totalItems := cmdBuilderFilenameInfo(videoUrl)

	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "ExtractFilenameInfo - Cmd.Start")

	var pResult []string
	var filenameInfo e.FilenameInfo

	pResult = executeProcess(stdout, false)
	//there are 22 properties that are to be claimed
	//this condition handles various errors from yt-dlp
	if len(pResult) >= (totalItems - 1) {
		filenameInfo = parseFilenameInfo(pResult, totalItems)

		// helper method fix
		filenameInfo.Channel = removeForbiddenCharsGeneric(filenameInfo.Channel)
		filenameInfo.Domain = removeForbiddenCharsGeneric(filenameInfo.Domain)
		filenameInfo.Title = removeForbiddenCharsGeneric(filenameInfo.Title)
	}
	filenameInfo.ContentId = contentId
	return filenameInfo
}

func (d *download) ExtractMetadata(input domain.PlaylistContentMeta) e.MediaInformation {

	// args, command, totalItems := cmdBuilderMetadata(d.p.Indicator)
	videoUrl := input.ContentId
	args, command, totalItems := cmdBuilderMetadata(videoUrl)

	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "Metadata - Cmd.Start")

	var pResult []string
	var mediaInfo e.MediaInformation

	pResult = executeProcess(stdout, false)
	//there are 22 properties that are to be claimed
	//this condition handles various errors from yt-dlp
	if len(pResult) >= (totalItems - 1) {
		mediaInfo = parseResults(pResult, VideoMetadata)

		// helper method fix
		mediaInfo = removeForbiddenChars(mediaInfo)
		mediaInfo = cleanDirectoryStructureFields(mediaInfo)

		return mediaInfo
	}
	return e.MediaInformation{}
}

func (d *download) ExtractMediaContent(filenameInfo domain.FilenameInfo) int {

	var (
		args    string
		command string
	)

	//get download command
	args, command = cmdBuilderDownload(filenameInfo)
	logCommand := command + Space + args

	//log executed command - in activity log later
	_ = logCommand
	cmd, stdout := buildProcess(args, GetCommandString())

	err := cmd.Start()
	handleErrors(err, "Download - Cmd.Start")

	pResult := executeDownloadProcess(stdout)
	_, errors, results := stripResultSections(pResult)

	//results are not really needed - except maybe to check for errors.
	_ = errors
	_ = results

	return 0
}

func (d *download) ExtractThumbnail(filenameInfo e.FilenameInfo) e.Files {

	args, command := cmdBuilderThumbnails(filenameInfo.ContentId, filenameInfo)
	logCommand := command + Space + args

	//log executed command - in activity log later
	_ = logCommand
	cmd, stdout := buildProcess(args, GetCommandString())

	err := cmd.Start()
	handleErrors(err, "ExtractThumbnail - Cmd.Start")

	pResult := executeProcess(stdout, true)
	_, lstErrors, results := stripResultSections(pResult)
	_ = results

	var filepath string
	if lstErrors == nil {
		_, filepath = buildDownloadPath(filenameInfo, e.Thumbnail)
	} else if len(lstErrors) > 0 {
		return e.Files{}
	}

	//sort out below method
	result := getVideoThumbnailFiles(filenameInfo.ContentId, filepath, filenameInfo.Id, filenameInfo.PlaylistId, filenameInfo.ThumbnailUrl)
	return result
}

func (d *download) ExtractSubtitles(filenameInfo e.FilenameInfo) e.Files {

	//#region [download thumbnails]

	args, command := cmdBuilderSubtitles(filenameInfo)
	logCommand := command + Space + args

	//log executed command - in activity log later
	_ = logCommand
	cmd, stdout := buildProcess(args, GetCommandString())

	err := cmd.Start()
	handleErrors(err, "Subtitles - Cmd.Start")

	pResult := executeProcess(stdout, true)
	_, errors, results := stripResultSections(pResult)

	//results are not really needed - except maybe to check for errors.
	_ = errors
	_ = results

	if len(errors) > 0 {
		//Show error on UI
		fmt.Println(errors)
		return e.Files{}
	}

	//#endregion

	//Get Subtitles FilePaths
	fp := getFilepaths(e.Filepath{Domain: filenameInfo.Domain, Channel: filenameInfo.Channel, PlaylistTitle: ""}, e.Subtitles)

	c, err := os.ReadDir(fp)
	handleErrors(err, "network - ExtractSubtitles")

	smiIndex := 0
	var file e.Files
	//Todo: re-write this to compare filenames after removing all special characters.
	//if there is a match then do the assignment.
	for _, entry := range c {
		info, _ := entry.Info()
		splits := strings.SplitN(info.Name(), ".", -1)
		fs_filename := info.Name()

		if strings.Contains(fs_filename, filenameInfo.ContentId) {
			file = e.Files{
				VideoId:      filenameInfo.Id,
				PlaylistId:   filenameInfo.PlaylistId,
				FileType:     "Subtitles",
				SourceId:     e.Downloaded,
				FilePath:     fp,
				FileName:     info.Name(),
				Extension:    splits[len(splits)-1],
				FileSize:     int(info.Size()),
				FileSizeUnit: "bytes",
				NetworkPath:  "",
				IsDeleted:    0,
				CreatedDate:  info.ModTime().Unix(),
			}
			smiIndex++
		}

	}
	return file
}

func (d *download) GetDownloadedMediaFileInfo(smi e.SavedInfo, filenameInfo e.FilenameInfo) e.Files {

	//Get Video FilePaths
	_, dirPath := buildDownloadPath(filenameInfo, e.Video)
	c, err := os.ReadDir(dirPath)
	handleErrors(err, "network - ExtractMediaContent")

	smiIndex := 0
	var files e.Files
	for _, entry := range c {
		info, _ := entry.Info()
		splits := strings.SplitN(info.Name(), ".", -1)
		fs_filename := info.Name()

		if strings.Contains(fs_filename, smi.YoutubeVideoId) {
			files = e.Files{
				VideoId:      smi.VideoId,
				PlaylistId:   filenameInfo.PlaylistId,
				FileType:     "Video",
				SourceId:     e.Downloaded,
				FilePath:     dirPath,
				FileName:     info.Name(),
				Extension:    splits[len(splits)-1],
				FileSize:     int(info.Size()),
				FileSizeUnit: "bytes",
				NetworkPath:  smi.MediaInfo.WebpageURL,
				IsDeleted:    0,
				CreatedDate:  info.ModTime().Unix(),
			}
			smiIndex++
		}
	}

	return files
}

// #endregion

// #region [private methods]

func executeProcess(stdout io.ReadCloser, isFileDownload bool) []string {
	// var result string
	// activeItem := g.NewActiveItem()
	// activeItem[0].VideoURL = "Metadata"
	var b bytes.Buffer
	for {
		//Read data from pipe into temp
		temp := make([]byte, 4096)
		n, e := stdout.Read(temp)
		b.WriteString(string(temp[:n]))

		// result := b.String()
		// results := strings.Split(result, "\n")
		// activeItem[0].StatusMessage = results[len(results)-2]

		//terminate loop at eof
		if e != nil {
			fmt.Println("Error Reading:", e)
			break
		}
	}

	//sanitize and return
	var results []string
	if !isFileDownload {
		results = sanitizeResults(b)
	}
	return results
}

// func executeDownloadProcess(stdout io.ReadCloser, activeItem []g.DownloadStatus) []string {
func executeDownloadProcess(stdout io.ReadCloser) []string {

	//Update State
	// activeItem[0].State = g.Downloading
	arrayOffset := 2

	// var result string
	var b bytes.Buffer
	for {
		//Read data from pipe into temp
		temp := make([]byte, 2048)
		n, err := stdout.Read(temp)
		b.WriteString(string(temp[:n]))
		result := b.String()
		results := strings.Split(result, "\n")

		//handle empty strings at the end
		if len(results)-arrayOffset >= 0 {
			// activeItem[0].StatusMessage = results[len(results)-arrayOffset]
			// fmt.Println("MESSAGE VALUE: ", activeItem[0].StatusMessage)
		}

		//terminate loop at eof
		if err != nil {
			fmt.Println("Error Reading:", err)
			if err == io.EOF {
				//activeItem[0].StatusMessage = "Download completed successfully."
			} else {
				//activeItem[0].StatusMessage = "Error:" + err.Error()
			}
			break
		}
	}
	//change state to stop reprocessing.
	// activeItem[0].State = g.Completed

	result := b.String()
	results := strings.Split(result, "\n")

	return results
}

func sanitizeResults(b bytes.Buffer) []string {

	//split result by newlines
	result := b.String()
	results := strings.Split(result, "\n")

	for i := 0; i < len(results); i++ {
		//valid json require keys and values to be enclosed in double quotes, not single quotes
		results[i] = proximityQuoteReplacement(results[i])

		//handle description differently
		if strings.Contains(results[i], "description") {
			results[i] = strings.ReplaceAll(results[i], "\\n", "<br />")
		}
	}

	//remove newlines from the end
	if results[len(results)-1] == "" {
		results = results[:len(results)-1]
	}
	return results
}

// valid json require keys and values to be enclosed in double quotes, not single quotes
// escaped single quotes within data replaced with escaped double quotes
func proximityQuoteReplacement(data string) string {

	//check for array
	seqArraryCheck1 := strings.Index(data, ": ['")
	seqArraryCheck2 := strings.LastIndex(data, "']")

	isArray := false
	if seqArraryCheck1 >= 0 && seqArraryCheck2 >= 0 {
		isArray = true
	}
	//replace double-quotes with escaped-double-quotes
	//if condition because in some case double quotes is text qualifier
	//for the value field while the key is still using single quotes
	if !strings.Contains(data, ": \"") && !isArray {
		data = strings.ReplaceAll(data, "\"", "\\\"")
	}

	//replace boolean values not enclosed in any quotes to be of lower case atleast.
	if strings.Contains(data, ": True}") || strings.Contains(data, ": False}") {
		data = strings.ToLower(data)
	}

	dQ := []byte("\"")[0]
	b := []byte(data)

	if seqArraryCheck1 >= 0 && seqArraryCheck2 >= 0 {
		//handle this case for arrays also
		data = strings.ReplaceAll(data, "\\'", "'")

		startIndex := strings.Index(data, ": [")
		endIndex := strings.Index(data, "]}")

		//replace inside arrays
		strData := strings.Split(data[startIndex+3:endIndex], ", ")
		var resultant []string
		for _, elem := range strData {
			if strings.Index(elem, "'") == 0 && strings.LastIndex(elem, "'") == len(elem)-1 {
				elem = strings.ReplaceAll(elem, "'", "\"")
			}
			resultant = append(resultant, elem)
		}

		ans := strings.Join(resultant, ", ")
		//value fix
		data = strings.ReplaceAll(data, data[startIndex+3:endIndex], ans)

		//key fix
		keyEndIdx := strings.Index(data, ":")
		replaced := strings.ReplaceAll(data[0:keyEndIdx], "'", "\"")
		data = strings.ReplaceAll(data, data[0:keyEndIdx], replaced)
		return data
	}

	if seq1 := strings.Index(data, "{'"); seq1 >= 0 {
		b[seq1+1] = dQ
	}

	if seq2 := strings.Index(data, "':"); seq2 >= 0 {
		b[seq2] = dQ
	}

	if seq3 := strings.Index(data, ": '"); seq3 >= 0 {
		b[seq3+2] = dQ
	}

	if seq4 := strings.LastIndex(data, "'}"); seq4 >= 0 {
		b[seq4] = dQ
	}

	data = string(b)

	// Case 1: replace escaped-single-quotes that may appear inside data with single-quotes since it is no longer the qualifier
	// Case 2: placed at the end since, a string like -> {'channel': '/ Mad Moose Media \\\\'} <- here a backslash is escaped
	// right before the	single quote but the single quote is not inside data and that does not need to be escaped.
	data = strings.ReplaceAll(data, "\\'", "'")

	return data
}

// region [Result Parsers]
func parseResults(pResult []string, metadataType int) e.MediaInformation {

	_, _, results := stripResultSections(pResult)
	metaItemsCount := 0
	for _, elem := range BuilderOptions() {
		if metadataType == VideoMetadata && elem.Group.Video.Metadata && elem.DataField {
			metaItemsCount++
		} else if metadataType == PlaylistMetadata && elem.Group.Playlist.Metadata && elem.DataField {
			metaItemsCount++
		}
	}

	mediaInfo := e.MediaInformation{}
	for i := 0; i < metaItemsCount; i++ {
		//Unmarshall is unreliable since the json coming from yt-dlp is invalid.
		//case statement for handling each field is required here because unmarshal is shit.
		if results[i][0] == '{' && results[i][len(results[i])-1] == '}' {
			err := json.Unmarshal([]byte(results[i]), &mediaInfo)
			handleErrors(err, "JSON Parser")
		}
	}

	if c.Config("PATCHING", false) == "enabled" {
		mediaInfo = patchDataField(mediaInfo)
	}

	return mediaInfo
}

func parsePlaylistInfo(pResult []string, itemsCount int) []e.Playlist {

	_, _, results := stripResultSections(pResult)
	titleKey := "{\"id\":" //Id starts with this in results json
	var itemCount int      //totalItems based on start key being Id

	//count Id keys
	for _, elem := range results {
		//if the string starts with a title key, it is a new playlist
		if strings.Index(elem, titleKey) == 0 {
			itemCount = itemCount + 1
		}
	}

	var lstplaylistInfo []e.Playlist
	for k := 0; k < itemCount; k++ {
		playlistInfo := e.Playlist{}
		for i := (0 + k*itemsCount); i < (k+1)*itemsCount; i++ {
			if results[i][0] == '{' && results[i][len(results[i])-1] == '}' {
				err := json.Unmarshal([]byte(results[i]), &playlistInfo)
				handleErrors(err, "JSON Parser")
			}
		}
		lstplaylistInfo = append(lstplaylistInfo, playlistInfo)
	}
	return lstplaylistInfo
}

func parseUploaderInfo(pResult []string, itemsCount int) e.PlaylistUploader {

	_, _, results := stripResultSections(pResult)
	titleKey := "{\"playlist_id\":" //Id starts with this in results json
	var itemCount int               //totalItems based on start key being Id

	//count Id keys
	for _, elem := range results {
		//if the string starts with a title key, it is a new playlist
		if strings.Index(elem, titleKey) == 0 {
			itemCount = itemCount + 1
		}
	}

	var lstplaylistInfo []e.PlaylistUploader
	for k := 0; k < itemCount; k++ {
		playlistInfo := e.PlaylistUploader{}
		for i := (0 + k*itemsCount); i < (k+1)*itemsCount; i++ {
			if results[i][0] == '{' && results[i][len(results[i])-1] == '}' {
				err := json.Unmarshal([]byte(results[i]), &playlistInfo)
				handleErrors(err, "JSON Parser")
			}
		}
		lstplaylistInfo = append(lstplaylistInfo, playlistInfo)
	}

	var uploaderInfo domain.PlaylistUploader
	if len(lstplaylistInfo) > 0 {
		uploaderInfo = lstplaylistInfo[0]
	}

	return uploaderInfo
}

func parseContentInfo(pResult []string, itemsCount int) []domain.PlaylistContentMeta {

	_, _, results := stripResultSections(pResult)
	titleKey := "{\"id\":" //Id starts with this in results json
	var itemCount int      //totalItems based on start key being Id

	//count Id keys
	for _, elem := range results {
		//if the string starts with a title key, it is a new playlist
		if strings.Index(elem, titleKey) == 0 {
			itemCount = itemCount + 1
		}
	}

	var lstContentInfo []e.PlaylistContentMeta
	for k := 0; k < itemCount; k++ {
		contentInfo := e.PlaylistContentMeta{}
		for i := (0 + k*itemsCount); i < (k+1)*itemsCount; i++ {
			if results[i][0] == '{' && results[i][len(results[i])-1] == '}' {
				err := json.Unmarshal([]byte(results[i]), &contentInfo)
				handleErrors(err, "JSON Parser")
			}
		}
		lstContentInfo = append(lstContentInfo, contentInfo)
	}
	return lstContentInfo
}

func parseFilenameInfo(pResult []string, itemsCount int) domain.FilenameInfo {

	_, _, results := stripResultSections(pResult)
	titleKey := "{\"title\":" //Id starts with this in results json
	var itemCount int         //totalItems based on start key being Id

	//count Id keys
	for _, elem := range results {
		//if the string starts with a title key, it is a new playlist
		if strings.Index(elem, titleKey) == 0 {
			itemCount = itemCount + 1
		}
	}

	var lstContentInfo []e.FilenameInfo
	for k := 0; k < itemCount; k++ {
		contentInfo := e.FilenameInfo{}
		for i := (0 + k*itemsCount); i < (k+1)*itemsCount; i++ {
			if results[i][0] == '{' && results[i][len(results[i])-1] == '}' {
				err := json.Unmarshal([]byte(results[i]), &contentInfo)
				handleErrors(err, "JSON Parser")
			}
		}
		lstContentInfo = append(lstContentInfo, contentInfo)
	}
	return lstContentInfo[0]
}

//endregion [Result Parsers]

func stripResultSections(pResult []string) ([]string, []string, []string) {

	var warnings []string
	var errors []string
	var results []string
	var previous string

	for _, elem := range pResult {
		if val := strings.Index(elem, WARNING); val == 0 {
			warnings = append(warnings, elem)
			previous = WARNING
		} else if val := strings.Index(elem, ERROR); val == 0 {
			errors = append(errors, elem)
			previous = ERROR
		} else if val := strings.Index(elem, ANSWER_START); val == 0 {
			results = append(results, elem)
		} else {
			//append to previous entry if nothing matches -- most tested and stable solution
			if previous == WARNING {
				warnings = append(warnings, elem)
			} else if previous == ERROR {
				errors = append(errors, elem)
			}
		}
	}

	return warnings, errors, results
}

func stripResultSectionsPlainText(pResult []string, answer_prefix string) ([]string, []string, []string) {

	var warnings []string
	var errors []string
	var results []string
	var previous string

	for _, elem := range pResult {
		if val := strings.Index(elem, WARNING); val == 0 {
			warnings = append(warnings, elem)
			previous = WARNING
		} else if val := strings.Index(elem, ERROR); val == 0 {
			errors = append(errors, elem)
			previous = ERROR
		} else if val := strings.Index(elem, answer_prefix); val == 0 {
			results = append(results, elem)
		} else {
			//append to previous entry if nothing matches -- most tested and stable solution
			if previous == WARNING {
				warnings = append(warnings, elem)
			} else if previous == ERROR {
				errors = append(errors, elem)
			}
		}
	}

	return warnings, errors, results
}

func handleErrors(err error, methodName string) {
	if err != nil {
		fmt.Println("pkg dowonload", methodName, err)
	}
}

// avoiding reflection here
func patchDataField(mediaInfo e.MediaInformation) e.MediaInformation {

	const plainChannel string = "Channel: "
	const plainTitle string = "Title: "
	const plainDescription string = "Description: "
	const plainTags string = "Tags: "
	const plainCategories string = "Categories: "

	var queries []string
	//Only checking for title description, Channel Name errors here.
	switch {
	case mediaInfo.Channel == "":
		queries = append(queries, Plaintext_Channel)
	case mediaInfo.Title == "":
		queries = append(queries, Plaintext_Title)
	case mediaInfo.Description == "":
		queries = append(queries, Plaintext_Description)
	case len(mediaInfo.Tags) == 0:
		queries = append(queries, Plaintext_Tags)
	case len(mediaInfo.Categories) == 0:
		queries = append(queries, Plaintext_Categories)
	default:
		break
	}

	for _, elem := range queries {
		var args []string
		args = append(args, mediaInfo.OriginalURL)
		args = append(args, SkipDownload)
		args = append(args, elem)

		options := strings.Join(args, Space)
		cmd, stdout := buildProcess(options, GetCommandString())

		err := cmd.Start()
		handleErrors(err, "patchDataField - Cmd.Start")

		procResult := executeProcess(stdout, false)

		for i := range procResult {
			if idx := strings.Index(procResult[i], plainChannel); idx == 0 {
				mediaInfo.Channel = procResult[i][len(plainChannel):]
			} else if idx := strings.Index(procResult[i], plainTitle); idx == 0 {
				mediaInfo.Title = procResult[i][len(plainTitle):]
			} else if idx := strings.Index(procResult[i], plainDescription); idx == 0 {
				if len(procResult[i]) > 1 {
					//first
					procResult[0] = procResult[0][len(plainDescription):]
					//rest
					mediaInfo.Description = strings.Join(procResult, "<br />")
				} else {
					mediaInfo.Description = procResult[i][len(plainDescription):]
				}
			} else if idx := strings.Index(procResult[i], plainTags); idx == 0 {
				mediaInfo.Tags = []string{procResult[i][len(plainTags):]}
			} else if idx := strings.Index(procResult[i], plainCategories); idx == 0 {
				mediaInfo.Categories = []string{procResult[i][len(plainCategories):]}
			}
		}
	}

	return mediaInfo
}

func getVideoThumbnailFiles(contentId string, filepath string, videoId int, playlistId int, contentUrl string) e.Files {

	//read all files in the directory where the thumbnails are saved
	c, err := os.ReadDir(filepath)
	handleErrors(err, "network - ExtractThumbnail")

	var file e.Files
	for _, entry := range c {
		info, _ := entry.Info()
		splits := strings.SplitN(info.Name(), ".", -1)
		fs_filename := info.Name()

		if strings.Contains(fs_filename, contentId) {
			file = e.Files{
				VideoId:      videoId,
				PlaylistId:   playlistId,
				FileType:     "Thumbnail",
				SourceId:     e.Downloaded,
				FilePath:     filepath,
				FileName:     info.Name(),
				Extension:    splits[len(splits)-1],
				FileSize:     int(info.Size()),
				FileSizeUnit: "bytes",
				NetworkPath:  contentUrl,
				IsDeleted:    0,
				CreatedDate:  info.ModTime().Unix(),
			}
		}
	}

	return file
}

// #endregion

// Result Type for entity binding and result parsing
const (
	Indicator            = iota
	VideoMinimalMetadata = iota
	VideoMetadata        = iota
	PlaylistMetadata     = iota
	Download             = iota
)
