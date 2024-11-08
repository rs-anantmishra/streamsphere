package extractor

import (
	"os"
	"runtime"
	"strings"

	c "github.com/rs-anantmishra/metubeplus/config"
	e "github.com/rs-anantmishra/metubeplus/pkg/entities"
)

type CSwitch struct {
	Index     int
	Name      string
	Value     string
	DataField bool
	Group     FxGroups
}

type FxGroups struct {
	Playlist Functions
	Video    Functions
	Generic  Functions
}

const (
	Generic  = iota
	Video    = iota
	Playlist = iota
)

type Functions struct {
	Metadata  bool
	Download  bool
	Subtitle  bool
	Thumbnail bool
}

func BuilderOptions() []CSwitch {

	//these true false patterns are talking about default download options
	//this forms the basis of the execute-custom-commands that may be implemented later on
	//flexibility of cutom commands may still be a question mark
	//ideally this should be moved to a db or read from a config file.
	defaults := []CSwitch{
		{Index: 1, Name: `Filepath`, Value: Filepath, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 2, Name: `Channel`, Value: Channel, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 3, Name: `Title`, Value: Title, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 4, Name: `Description`, Value: Description, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 5, Name: `Extension`, Value: Extension, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 6, Name: `Duration`, Value: Duration, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 7, Name: `URLDomain`, Value: URLDomain, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 8, Name: `OriginalURL`, Value: OriginalURL, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 9, Name: `Tags`, Value: Tags, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 10, Name: `YTFormatString`, Value: YTFormatString, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 11, Name: `FileSizeApprox`, Value: FileSizeApprox, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 12, Name: `FormatNote`, Value: FormatNote, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 13, Name: `Resolution`, Value: Resolution, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 14, Name: `Categories`, Value: Categories, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 15, Name: `YoutubeVideoId`, Value: YoutubeVideoId, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 16, Name: `Availability`, Value: Availability, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 17, Name: `LiveStatus`, Value: LiveStatus, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 18, Name: `ChannelId`, Value: ChannelId, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 19, Name: `ChannelURL`, Value: ChannelURL, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 20, Name: `PlaylistTitle`, Value: PlaylistTitle, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 21, Name: `PlaylistIndex`, Value: PlaylistIndex, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 22, Name: `PlaylistCount`, Value: PlaylistCount, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 23, Name: `PlaylistId`, Value: PlaylistId, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 24, Name: `ShowProgress`, Value: ShowProgress, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true}},
		},
		{Index: 25, Name: `ProgressDelta`, Value: ProgressDelta, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true}},
		},
		{Index: 26, Name: `QuietDownload`, Value: QuietDownload, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: true, Subtitle: false, Thumbnail: false}},
		},
		{Index: 27, Name: `ProgressNewline`, Value: ProgressNewline, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true}},
		},
		{Index: 28, Name: `SkipDownload`, Value: SkipDownload, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: true, Thumbnail: true},
			Video:    Functions{Metadata: true, Download: false, Subtitle: true, Thumbnail: true}},
		},
		{Index: 29, Name: `WriteSubtitles`, Value: WriteSubtitles, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: true, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: true, Thumbnail: false}},
		},
		{Index: 30, Name: `WriteThumbnail`, Value: WriteThumbnail, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: true}},
		},
		{Index: 31, Name: `MediaDirectory`, Value: GetMediaDirectory(true), DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: true, Subtitle: true, Thumbnail: true}},
		},
		{Index: 32, Name: `OutputPlaylistVideoFile`, Value: OutputPlaylistVideoFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: true, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 33, Name: `OutputPlaylistSubtitleFile`, Value: OutputPlaylistSubtitleFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: true, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 34, Name: `OutputPlaylistThumbnailFile`, Value: OutputPlaylistThumbnailFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: true},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 35, Name: `OutputVideoFile`, Value: OutputVideoFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: true, Subtitle: false, Thumbnail: false}},
		},
		{Index: 36, Name: `OutputSubtitleFile`, Value: OutputSubtitleFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: true, Thumbnail: false}},
		},
		{Index: 37, Name: `OutputThumbnailFile`, Value: OutputThumbnailFile, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: false, Download: false, Subtitle: false, Thumbnail: true}},
		},
		{Index: 38, Name: `ThumbnailURL`, Value: ThumbnailURL, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 39, Name: `License`, Value: License, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 40, Name: `ChannelFollowerCount`, Value: ChannelFollowerCount, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 41, Name: `UploadDate`, Value: UploadDate, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 42, Name: `ReleaseTimestamp`, Value: ReleaseTimestamp, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 43, Name: `ModifiedTimestamp`, Value: ModifiedTimestamp, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 44, Name: `ViewCount`, Value: ViewCount, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 45, Name: `LikeCount`, Value: LikeCount, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 46, Name: `DislikeCount`, Value: DislikeCount, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 47, Name: `AgeLimit`, Value: AgeLimit, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 48, Name: `PlayableInEmbed`, Value: PlayableInEmbed, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 49, Name: `WebpageURL`, Value: WebpageURL, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 50, Name: `EncodingExperimental`, Value: EncodingExperimental, DataField: false, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 51, Name: `PlaylistChannel`, Value: PlaylistChannel, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 52, Name: `PlaylistChannelId`, Value: PlaylistChannelId, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 53, Name: `PlaylistUploader`, Value: PlaylistUploader, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		{Index: 54, Name: `PlaylistUploaderId`, Value: PlaylistUploaderId, DataField: true, Group: FxGroups{
			Playlist: Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false},
			Video:    Functions{Metadata: true, Download: false, Subtitle: false, Thumbnail: false}},
		},
		//Audio only file options to be added later
	}

	return defaults
}

func GetCommandString() string {
	cmdPath := c.Config("YTDLP_PATH", true)
	return cmdPath + string(os.PathSeparator) + GetCommandName()
}

func GetVideoFilepath(fp e.Filepath, fType int) string {
	var result string
	if fType == e.Thumbnail {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, "Videos", "Thumbnails"}, string(os.PathSeparator))
	} else if fType == e.Subtitles {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, "Videos", "Subtitles"}, string(os.PathSeparator))
	} else {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, "Videos"}, string(os.PathSeparator))
	}
	return result
}

func GetPlaylistFilepath(fp e.Filepath, fType int) string {

	var result string
	if fType == e.Thumbnail {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, fp.PlaylistTitle, "Thumbnails"}, string(os.PathSeparator))
	} else if fType == e.Subtitles {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, fp.PlaylistTitle, "Subtitles"}, string(os.PathSeparator))
	} else {
		result = strings.Join([]string{GetMediaDirectory(false), fp.Domain, fp.Channel, fp.PlaylistTitle}, string(os.PathSeparator))
	}
	return result
}

func cmdBuilderMetadata(url string) (string, string, int) {

	var args []string
	args = append(args, "\""+url+"\"")
	totalItems := 0

	bo := BuilderOptions()
	for _, elem := range bo {
		//Handle Video
		if elem.Group.Video.Metadata {
			if runtime.GOOS == "linux" {
				elem.Value = strings.ReplaceAll(elem.Value, "(", "'(")
				elem.Value = strings.ReplaceAll(elem.Value, ")", ")'")
			}

			if elem.DataField {
				totalItems++
			}
			args = append(args, elem.Value)
		}
	}
	arguments := strings.Join(args, Space)

	cmdPath := c.Config("YTDLP_PATH", true)
	cmd := cmdPath + string(os.PathSeparator) + GetCommandName()

	return arguments, cmd, totalItems
}

// Download Media Content
func cmdBuilderDownload(url string, savedInfo e.SavedInfo) (string, string) {

	var args []string
	args = append(args, "\""+url+"\"")

	//this is to get rid of the problem with special chars that windows does not support
	//while maintaining the directory structure and aethetics for fs access to your data
	contentFilepath, _ := buildDownloadPath(savedInfo, e.Video)

	bo := BuilderOptions()
	for _, elem := range bo {

		//Handle Video
		// if elem.Group.Video.Download && savedInfo.PlaylistId < 0 {
		if elem.Group.Video.Download {
			switch elem.Name {
			case "OutputVideoFile":
				args = append(args, contentFilepath)
			default:
				args = append(args, elem.Value)
			}
		}

		//Handle Playlist
		if elem.Group.Playlist.Download && savedInfo.PlaylistId > 0 {
			switch elem.Name {
			case "OutputPlaylistVideoFile":
				args = append(args, contentFilepath)
			default:
				args = append(args, elem.Value)
			}
		}
	}

	arguments := strings.Join(args, Space)
	cmdPath := c.Config("YTDLP_PATH", true)
	cmd := cmdPath + string(os.PathSeparator) + GetCommandName()

	return arguments, cmd
}

func cmdBuilderSubtitles(url string, savedInfo e.SavedInfo) (string, string) {

	var args []string
	args = append(args, "\""+url+"\"")

	//this is to get rid of the problem with special chars that windows does not support
	//while maintaining the directory structure and aethetics for fs access to your data
	subtitlesFilepath, _ := buildDownloadPath(savedInfo, e.Subtitles)

	bo := BuilderOptions()
	for _, elem := range bo {

		//Handle Video
		if elem.Group.Video.Subtitle && savedInfo.PlaylistId < 0 {
			switch elem.Name {
			case "OutputSubtitleFile":
				args = append(args, subtitlesFilepath)
			default:
				args = append(args, elem.Value)
			}
		}

		//Handle Playlist
		if elem.Group.Playlist.Subtitle && savedInfo.PlaylistId > 0 {
			switch elem.Name {
			case "OutputPlaylistSubtitleFile":
				args = append(args, subtitlesFilepath)
			default:
				args = append(args, elem.Value)
			}
		}
	}

	arguments := strings.Join(args, Space)
	cmdPath := c.Config("YTDLP_PATH", true)
	cmd := cmdPath + string(os.PathSeparator) + GetCommandName()

	return arguments, cmd
}

func cmdBuilderThumbnails(url string, savedInfo e.SavedInfo) (string, string) {

	var args []string
	args = append(args, "\""+url+"\"")

	//this is to get rid of the problem with special chars that windows does not support
	//while maintaining the directory structure and aethetics for fs access to your data
	thumbnailFilepath, _ := buildDownloadPath(savedInfo, e.Thumbnail)

	bo := BuilderOptions()
	for _, elem := range bo {

		//Handle Video
		if elem.Group.Video.Thumbnail && savedInfo.PlaylistId < 0 {
			switch elem.Name {
			case "OutputThumbnailFile":
				args = append(args, thumbnailFilepath)
			default:
				args = append(args, elem.Value)
			}
		}

		//Handle Playlist
		if elem.Group.Playlist.Thumbnail && savedInfo.PlaylistId > 0 {
			switch elem.Name {
			case "OutputPlaylistThumbnailFile":
				args = append(args, thumbnailFilepath)
			default:
				args = append(args, elem.Value)
			}
		}
	}

	arguments := strings.Join(args, Space)
	cmdPath := c.Config("YTDLP_PATH", true)
	cmd := cmdPath + string(os.PathSeparator) + GetCommandName()

	return arguments, cmd
}

func buildDownloadPath(savedInfo e.SavedInfo, pathType int) (string, string) {

	pathResult := ""
	dirResultPath := ""
	space := " "
	sep := string(os.PathSeparator)

	//Videos - In playlists or otherwise.
	{
		var directories []string
		directories = append(directories, savedInfo.MediaInfo.Domain)
		directories = append(directories, savedInfo.MediaInfo.Channel)
		directories = append(directories, "Videos")
		{
			if pathType == e.Thumbnail {
				directories = append(directories, "Thumbnails")
			} else if pathType == e.Subtitles {
				directories = append(directories, "Subtitles")
			}
		}
		dirResultPath = GetMediaDirectory(false) + string(os.PathSeparator) + strings.Join(directories, sep) //separate result for path only
		directories = append(directories, savedInfo.MediaInfo.Title+space+"[%(id)s].%(ext)s")
		pathResult = strings.Join(directories, sep)
	}

	switch pathType {
	case e.Thumbnail:
		pathResult = "\"thumbnail:" + pathResult + "\""
		pathResult = `-o ` + pathResult
	case e.Subtitles:
		pathResult = "\"subtitle:" + pathResult + "\""
		pathResult = `-o ` + pathResult
	case e.Video:
		pathResult = "\"" + pathResult + "\""
		pathResult = `-o ` + pathResult
	}

	return pathResult, dirResultPath
}
