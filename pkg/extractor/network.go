package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	c "github.com/rs-anantmishra/streamsphere/config"
	e "github.com/rs-anantmishra/streamsphere/pkg/entities"
	g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

type IDownload interface {
	ExtractMetadata() ([]e.MediaInformation, e.Filepath)
	ExtractMediaContent(smi e.SavedInfo) int
	ExtractSubtitles(fp e.Filepath, lstSavedInfo []e.SavedInfo) []e.Files
	ExtractThumbnail(fp e.Filepath, lstSavedInfo []e.SavedInfo) []e.Files
	GetDownloadedMediaFileInfo(smi e.SavedInfo, fp e.Filepath) []e.Files
	Cleanup()
}

type download struct {
	p            e.IncomingRequest
	lstDownloads []g.DownloadStatus
}

func NewDownload(params e.IncomingRequest) IDownload {
	return &download{
		p: params,
	}
}

func (d *download) Cleanup() {
	var updated []g.DownloadStatus
	for i := 0; i < len(d.lstDownloads); i++ {
		ds := d.lstDownloads
		if ds[i].State != g.Completed {
			updated = append(updated, ds[i])
		}
	}
	d.lstDownloads = updated
}

func (d *download) ExtractMetadata() ([]e.MediaInformation, e.Filepath) {

	args, command, totalItems := cmdBuilderMetadata(d.p.Indicator)
	logCommand := command + Space + args
	_ = logCommand //log executed command - in activity log later

	strCommand := GetCommandString()
	cmd, stdout := buildProcess(args, strCommand)

	err := cmd.Start()
	handleErrors(err, "Metadata - Cmd.Start")

	var pResult []string
	var mediaInfo []e.MediaInformation

	pResult = executeProcess(stdout, false)
	//there are 22 properties that are to be claimed
	//this condition handles various errors from yt-dlp
	if len(pResult) >= (totalItems - 1) {
		mediaInfo = parseResults(pResult, VideoMetadata)

		// helper method fix
		mediaInfo = removeForbiddenChars(mediaInfo)
		mediaInfo = cleanDirectoryStructureFields(mediaInfo)

		////////////////////////////////////
		//handle shortened URL /////////////
		////////////////////////////////////

		fp := e.Filepath{Domain: mediaInfo[0].Domain, Channel: mediaInfo[0].Channel, PlaylistTitle: mediaInfo[0].PlaylistTitle}
		return mediaInfo, fp
	}
	return nil, e.Filepath{}
}

func (d *download) ExtractMediaContent(smi e.SavedInfo) int {

	activeItem := g.NewActiveItem()

	var (
		args    string
		command string
	)

	//get download command
	args, command = cmdBuilderDownload(activeItem[0].VideoURL, smi)
	logCommand := command + Space + args

	//log executed command - in activity log later
	_ = logCommand
	cmd, stdout := buildProcess(args, GetCommandString())

	err := cmd.Start()
	handleErrors(err, "Download - Cmd.Start")

	pResult := executeDownloadProcess(stdout, activeItem)
	_, errors, results := stripResultSections(pResult)

	//results are not really needed - except maybe to check for errors.
	_ = errors
	_ = results

	return activeItem[0].State
}

func (d *download) ExtractThumbnail(fPath e.Filepath, lstSavedInfo []e.SavedInfo) []e.Files {

	//#region [download thumbnails]

	//for multi-channel playlists
	var thumbnailFilePaths []string

	for i := 0; i < len(lstSavedInfo); i++ {

		args, command := cmdBuilderThumbnails(lstSavedInfo[i].MediaInfo.WebpageURL, lstSavedInfo[i])
		logCommand := command + Space + args

		//log executed command - in activity log later
		_ = logCommand
		cmd, stdout := buildProcess(args, GetCommandString())

		err := cmd.Start()
		handleErrors(err, "Thumbnail - Cmd.Start")

		pResult := executeProcess(stdout, true)
		_, errors, results := stripResultSections(pResult)

		//if no erros and it is a multi channel playlist
		if errors == nil {
			_, dirPath := buildDownloadPath(lstSavedInfo[i], e.Thumbnail)
			thumbnailFilePaths = append(thumbnailFilePaths, dirPath)
		}

		//results are not really needed - except maybe to check for errors.
		_ = errors
		_ = results

		if len(errors) > 0 {
			//Show error on UI
			log.Error(errors)
			return []e.Files{}
		}
	}
	//#endregion

	//Sort & Unique thumbnailFilePaths
	slices.Sort(thumbnailFilePaths)
	uniqThumbnailFilePaths := slices.Compact(thumbnailFilePaths)

	result := getVideoThumbnailFiles(lstSavedInfo, uniqThumbnailFilePaths)
	return result
}

func (d *download) ExtractSubtitles(fPath e.Filepath, lstSavedInfo []e.SavedInfo) []e.Files {

	//#region [download thumbnails]

	for i := 0; i < len(lstSavedInfo); i++ {
		args, command := cmdBuilderSubtitles(lstSavedInfo[i].MediaInfo.WebpageURL, lstSavedInfo[i])
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
			log.Error(errors)
			return []e.Files{}
		}
	}
	//#endregion

	//Get Subtitles FilePaths
	fp := getFilepaths(lstSavedInfo[0].PlaylistId, fPath, e.Subtitles)

	c, err := os.ReadDir(fp)
	handleErrors(err, "network - ExtractSubtitles")

	smiIndex := 0
	var files []e.Files
	//Todo: re-write this to compare filenames after removing all special characters.
	//if there is a match then do the assignment.
	for _, entry := range c {
		info, _ := entry.Info()
		splits := strings.SplitN(info.Name(), ".", -1)
		fs_filename := info.Name()

		for _, saved := range lstSavedInfo {
			if strings.Contains(fs_filename, saved.YoutubeVideoId) {
				f := e.Files{
					VideoId:      lstSavedInfo[smiIndex].VideoId,
					PlaylistId:   lstSavedInfo[smiIndex].PlaylistId,
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
				files = append(files, f)
				smiIndex++
			}
		}
	}
	return files
}

func (d *download) GetDownloadedMediaFileInfo(smi e.SavedInfo, fPath e.Filepath) []e.Files {

	//Get Video FilePaths
	_, dirPath := buildDownloadPath(smi, e.Video)
	c, err := os.ReadDir(dirPath)
	handleErrors(err, "network - ExtractMediaContent")

	smiIndex := 0
	var files []e.Files
	for _, entry := range c {
		info, _ := entry.Info()
		splits := strings.SplitN(info.Name(), ".", -1)
		fs_filename := info.Name()

		if strings.Contains(fs_filename, smi.YoutubeVideoId) {
			f := e.Files{
				VideoId:      smi.VideoId,
				PlaylistId:   smi.PlaylistId,
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
			files = append(files, f)
			smiIndex++
		}
	}

	return files
}

// // Build a Process to execute & attach pipe to it here
// func buildProcess(args string, command string) (*exec.Cmd, io.ReadCloser) {
// 	args = command + " " + args
// 	// arguments := strings.Split(args, " ")

// 	cmd := exec.Command("bash", "-c", args)
// 	//cmd.SysProcAttr = &syscall.SysProcAttr{}
// 	//cmd.SysProcAttr.CmdLine = command + Space + args

// 	stdout, err := cmd.StdoutPipe()
// 	cmd.Stderr = cmd.Stdout
// 	handleErrors(err, "Metadata - StdoutPipe")

// 	return cmd, stdout
// }

func executeProcess(stdout io.ReadCloser, isFileDownload bool) []string {
	// var result string
	var b bytes.Buffer
	for {
		//Read data from pipe into temp
		temp := make([]byte, 4096)
		n, e := stdout.Read(temp)
		b.WriteString(string(temp[:n]))

		//terminate loop at eof
		if e != nil {
			log.Info("Error Reading:", e)
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

func executeDownloadProcess(stdout io.ReadCloser, activeItem []g.DownloadStatus) []string {

	//Update State
	activeItem[0].State = g.Downloading
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
			activeItem[0].StatusMessage = results[len(results)-arrayOffset]
			log.Info("MESSAGE VALUE: ", activeItem[0].StatusMessage)
		}

		//terminate loop at eof
		if err != nil {
			log.Info("Error Reading:", err)
			if err == io.EOF {
				activeItem[0].StatusMessage = "Download completed successfully."
			} else {
				activeItem[0].StatusMessage = "Error:" + err.Error()
			}
			break
		}
	}
	//change state to stop reprocessing.
	activeItem[0].State = g.Completed

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

	//replace double-quotes with escaped-double-quotes
	//if condition because in some case double quotes is text qualifier
	//for the value field while the key is still using single quotes
	if !strings.Contains(data, ": \"") {
		data = strings.ReplaceAll(data, "\"", "\\\"")
	}

	//replace boolean values not enclosed in any quotes to be of lower case atleast.
	if strings.Contains(data, ": True}") || strings.Contains(data, ": False}") {
		data = strings.ToLower(data)
	}

	dQ := []byte("\"")[0]
	b := []byte(data)

	seqArraryCheck1 := strings.Index(data, ": ['")
	seqArraryCheck2 := strings.LastIndex(data, "']")
	if seqArraryCheck1 >= 0 && seqArraryCheck2 >= 0 {
		data = strings.ReplaceAll(data, "'", "\"")
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

	// replace escaped-single-quotes that may appear inside data with single-quotes since it is no longer the qualifier
	// placed at the end since, a string like -> {'channel': '/ Mad Moose Media \\\\'} <- here a backslash is escaped
	// right before the	single quote but the single quote is not inside data and that does not need to be escaped.
	data = strings.ReplaceAll(data, "\\'", "'")

	return data
}

func parseResults(pResult []string, metadataType int) []e.MediaInformation {

	_, _, results := stripResultSections(pResult)
	var itemCount int         //no of Videos. If its
	titleKey := "{\"title\":" //title of the video starts with in results json

	//count title keys
	for _, elem := range results {
		//if the string starts with a title key, it is a new video
		if strings.Index(elem, titleKey) == 0 {
			itemCount = itemCount + 1
		}
	}

	metaItemsCount := 0
	for _, elem := range BuilderOptions() {
		if metadataType == VideoMetadata && elem.Group.Video.Metadata && elem.DataField {
			metaItemsCount++
		} else if metadataType == PlaylistMetadata && elem.Group.Playlist.Metadata && elem.DataField {
			metaItemsCount++
		}
	}

	var lstMediaInfo []e.MediaInformation
	for k := 0; k < itemCount; k++ {
		mediaInfo := e.MediaInformation{}
		for i := (0 + k*metaItemsCount); i < (k+1)*metaItemsCount; i++ {

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

		lstMediaInfo = append(lstMediaInfo, mediaInfo)
	}

	// log Properties that were bound \
	// lstMediaInfo

	return lstMediaInfo
}

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

func getVideoThumbnailFiles(lstSavedInfo []e.SavedInfo, thumbnailFilePaths []string) []e.Files {

	var files []e.Files
	for _, path := range thumbnailFilePaths {

		//read all files in the directory where the thumbnails are saved
		c, err := os.ReadDir(path)
		handleErrors(err, "network - ExtractThumbnail")

		for _, entry := range c {
			info, _ := entry.Info()
			splits := strings.SplitN(info.Name(), ".", -1)
			fs_filename := info.Name()

			for idx, saved := range lstSavedInfo {
				if strings.Contains(fs_filename, saved.YoutubeVideoId) {
					f := e.Files{
						VideoId:      lstSavedInfo[idx].VideoId,
						PlaylistId:   saved.PlaylistId,
						FileType:     "Thumbnail",
						SourceId:     e.Downloaded,
						FilePath:     path,
						FileName:     info.Name(),
						Extension:    splits[len(splits)-1],
						FileSize:     int(info.Size()),
						FileSizeUnit: "bytes",
						NetworkPath:  saved.MediaInfo.ThumbnailURL,
						IsDeleted:    0,
						CreatedDate:  info.ModTime().Unix(),
					}
					files = append(files, f)
				}
			}
		}

	}
	return files
}

// Result Type for entity binding and result parsing
const (
	Indicator            = iota
	VideoMinimalMetadata = iota
	VideoMetadata        = iota
	PlaylistMetadata     = iota
	Download             = iota
)
