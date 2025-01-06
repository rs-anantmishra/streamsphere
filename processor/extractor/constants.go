package extractor

import (
	"runtime"
	"strings"

	c "github.com/rs-anantmishra/streamsphere/utils/processor/config"
)

// --Options---------------------------------------------------------------------------------------//
const Separator string = `*MeTube+*`
const ShowProgress string = `--progress`
const ProgressDelta string = `--progress-delta 1` //seconds
const YoutubeVideoId string = `--print %(.{id})s`
const Availability string = `--print %(.{availability})s`
const LiveStatus string = `--print %(.{live_status})s`
const Filepath string = `--print %(.{filepath})s`
const Channel string = `--print %(.{channel})s`
const Title string = `--print %(.{title})s`
const Description string = `--print %(.{description})s`
const Extension string = `--print %(.{ext})s`
const Duration string = `--print %(.{duration})s`
const URLDomain string = `--print %(.{webpage_url_domain})s`
const OriginalURL string = `--print %(.{original_url})s`
const PlaylistTitle string = `--print %(.{playlist_title})s`
const PlaylistCount string = `--print %(.{playlist_count})s`
const PlaylistIndex string = `--print %(.{playlist_index})s`
const Tags string = `--print %(.{tags})s`
const YTFormatString string = `--print %(.{format})s`
const FileSizeApprox string = `--print %(.{filesize_approx})s`
const FormatNote string = `--print %(.{format_note})s`
const Resolution string = `--print %(.{resolution})s`
const Categories string = `--print %(.{categories})s`
const ChannelId string = `--print %(.{channel_id})s`
const ChannelURL string = `--print %(.{channel_url})s`
const PlaylistId string = `--print %(.{playlist_id})s`

// Added Later
const WebpageURL string = `--print %(.{webpage_url})s`
const ThumbnailURL string = `--print %(.{thumbnail})s`
const License string = `--print %(.{license})s`
const ChannelFollowerCount string = `--print %(.{channel_follower_count})s`
const UploadDate string = `--print %(.{upload_date})s`
const ReleaseTimestamp string = `--print %(.{release_timestamp})s`
const ModifiedTimestamp string = `--print %(.{modified_timestamp})s`
const ViewCount string = `--print %(.{view_count})s`
const LikeCount string = `--print %(.{like_count})s`
const DislikeCount string = `--print %(.{dislike_count})s`
const AgeLimit string = `--print %(.{age_limit})s`
const PlayableInEmbed string = `--print %(.{playable_in_embed})s`

const PlaylistChannel string = `--print %(.{playlist_channel})s`
const PlaylistChannelId string = `--print %(.{playlist_channel_id})s`
const PlaylistUploader string = `--print %(.{playlist_uploader})s`
const PlaylistUploaderId string = `--print %(.{playlist_uploader_id})s`

//list all video ids in channel
//yt-dlp_x86.exe --flat-playlist --print id https://www.youtube.com/@spinninrecords

//list all playlists from a channel
//yt-dlp_x86.exe --print "%(id)s|%(title)s|%(playlist_index)s" --flat-playlist "https://www.youtube.com/@spinninrecords/playlists"  --lazy-playlist

//all videos in a channel
//yt-dlp_x86.exe --print "%(id)s|%(title)s|%(playlist_index)s" --flat-playlist "https://www.youtube.com/@spinninrecords" --lazy-playlist

//videos in a playlist
//yt-dlp_x86.exe --print "%(id)s|%(title)s|%(playlist_index)s" --flat-playlist "https://www.youtube.com/playlist?list=PLA6DC52121CDE580A" --lazy-playlist

// --Options Plaintext-----------------------------------------------------------------------------//
const Plaintext_YoutubeVideoId string = `--print after_move:id`
const Plaintext_Availability string = `--print after_move:availability`
const Plaintext_LiveStatus string = `--print after_move:live_status`
const Plaintext_Filepath string = `--print after_move:filepath`
const Plaintext_Channel string = `--print before_dl:"Channel: %(channel)s"`             //Changed for patching data fields. `--print before_dl:channel`
const Plaintext_Title string = `--print before_dl:"Title: %(title)s"`                   //Changed for patching data fields. `--print before_dl:title`
const Plaintext_Description string = `--print before_dl:"Description: %(description)s"` //Changed for patching data fields. `--print before_dl:description`
const Plaintext_Tags string = `--print before_dl:"Tags: %(tags)s"`                      //Changed for patching data fields. `--print before_dl:tags`
const Plaintext_Categories string = `--print before_dl:"Categories: %(categories)s"`    //Changed for patching data fields. `--print before_dl:categories`
const Plaintext_Extension string = `--print before_dl:ext`
const Plaintext_Duration string = `--print before_dl:duration`
const Plaintext_URLDomain string = `--print before_dl:webpage_url_domain`
const Plaintext_OriginalURL string = `--print before_dl:original_url`
const Plaintext_PlaylistTitle string = `--print before_dl:playlist_title`
const Plaintext_PlaylistCount string = `--print before_dl:playlist_count`
const Plaintext_PlaylistIndex string = `--print before_dl:playlist_index`
const Plaintext_YTFormatString string = `--print before_dl:format`
const Plaintext_FileSizeApprox string = `--print before_dl:filesize_approx`
const Plaintext_FormatNote string = `--print before_dl:format_note`
const Plaintext_Resolution string = `--print before_dl:resolution`
const Plaintext_ChannelId string = `--print before_dl:channel_id`
const Plaintext_ChannelURL string = `--print before_dl:channel_url`
const Plaintext_PlaylistId string = `--print before_dl:playlist_id`
const Plaintext_ThumbnailURL string = `--print before_dl:thumbnail`

// Added later
const Plaintext_License string = `--print before_dl:license`
const Plaintext_ChannelFollowerCount string = `--print before_dl:channel_follower_count`
const Plaintext_UploadDate string = `--print before_dl:upload_date`
const Plaintext_ReleaseTimestamp string = `--print before_dl:release_timestamp`
const Plaintext_ModifiedTimestamp string = `--print before_dl:modified_timestamp`
const Plaintext_ViewCount string = `--print before_dl:view_count`
const Plaintext_LikeCount string = `--print before_dl:like_count`
const Plaintext_DislikeCount string = `--print before_dl:dislike_count`
const Plaintext_AgeLimit string = `--print before_dl:age_limit`
const Plaintext_PlayableInEmbed string = `--print before_dl:playable_in_embed`
const Plaintext_PlaylistChannel string = `--print before_dl:playlist_channel`
const Plaintext_PlaylistChannelId string = `--print before_dl:playlist_channel_id`
const Plaintext_PlaylistUploader string = `--print before_dl:playlist_uploader`
const Plaintext_PlaylistUploaderId string = `--print before_dl:playlist_uploader_id`

// --Extras----------------------------------------------------------------------------------------//
const WriteSubtitles string = `--write-auto-subs`
const WriteThumbnail string = `--write-thumbnail`
const SkipDownload string = `--skip-download`
const InfoJSON string = `--write-info-json`
const QuietDownload string = `--quiet`
const ProgressNewline string = `--newline`
const FlatPlaylist string = `--flat-playlist`
const LazyPlaylist string = `--lazy-playlist`

const EncodingExperimental string = `--encoding utf8`

// --Commands--------------------------------------------------------------------------------------//
// Playlist: Videos, Subtitles, Thumbnails
const OutputPlaylistVideoFile string = `-o "%(webpage_url_domain)s/%(channel)s/%(playlist)s/%(playlist_index)s - %(title)s [%(id)s].%(ext)s"`
const OutputPlaylistSubtitleFile string = `-o "subtitle:%(webpage_url_domain)s/%(channel)s/%(playlist)s/Subtitles/%(playlist_index)s - %(title)s [%(id)s].%(ext)s"`
const OutputPlaylistThumbnailFile string = `-o "thumbnail:%(webpage_url_domain)s/%(channel)s/%(playlist)s/Thumbnails/%(playlist_index)s - %(title)s [%(id)s].%(ext)s"`

// Videos: Videos, Subtitles, Thumbnails
const OutputVideoFile string = `-o "%(webpage_url_domain)s/%(channel)s/Videos/%(title)s [%(id)s].%(ext)s"`
const OutputSubtitleFile string = `-o "subtitle:%(webpage_url_domain)s/%(channel)s/Videos/Subtitles/%(title)s [%(id)s].%(ext)s"`
const OutputThumbnailFile string = `-o "thumbnail:%(webpage_url_domain)s/%(channel)s/Videos/Thumbnails/%(title)s [%(id)s].%(ext)s"`

// Output Parsing: Warnings & Errors
const WARNING string = `WARNING:`
const ERROR string = `ERROR:`
const ANSWER_START string = `{`

// --Testing----------------------------------------------------------------------------------------//
const TestURL1 string = `https://www.youtube.com/watch?v=GW2g-5WALrc`
const TestURL2 string = `https://www.youtube.com/watch?v=GW2g-5WALrc&list=PLFKeDWeuu3BZEBcRmolX6BDiFhK-GhCsd`
const TestURL3 string = `https://www.youtube.com/watch?v=-VC4FuG8P6Q`

// --Test Command Playlist--------------------------------------------------------------------------//
const TestCmdPlaylist string = `yt-dlp_x86.exe "https://www.youtube.com/watch?v=5WfiTHiU4x8&list=PLIhvC56v63IKrRHh3gvZZBAGvsvOhwrRF" -P "./" -o "%(webpage_url_domain)s/%(channel)s/%(playlist)s/%(playlist_index)s - %(title)s [%(id)s].%(ext)s" -o "subtitle:%(webpage_url_domain)s/%(channel)s/%(playlist)s/subs/%(playlist_index)s - %(title)s [%(id)s].%(ext)s" -o "thumbnail:%(webpage_url_domain)s/%(channel)s/%(playlist)s/thumbnails/%(playlist_index)s - %(title)s [%(id)s].%(ext)s" -S "res:240" --write-thumbnail --write-auto-subs`

// --Test Command Videos----------------------------------------------------------------------------//
const TestCmdVideo string = `yt-dlp_x86.exe "https://www.youtube.com/watch?v=AaseHnf0k2o" -P "./" -o "%(webpage_url_domain)s/%(channel)s/Videos/%(title)s [%(id)s].%(ext)s" -o "subtitle:%(webpage_url_domain)s/%(channel)s/Videos/Subtitles/%(title)s [%(id)s].%(ext)s" -o "thumbnail:%(webpage_url_domain)s/%(channel)s/Videos/Thumbnails/%(title)s [%(id)s].%(ext)s"  -S "res:240" --write-thumbnail --write-auto-subs`

// command path should be picked from .env
const CommandName string = `yt-dlp_x86.exe`
const Space string = " "
const FieldSeparator string = "#|#"
const Print string = `--print `
const ytUrl string = "https://youtube.com/"

func OptionsModifier(input string, remove_print bool, remove_postprocessor bool, answer_prefix string) string {
	const postprocessor_split string = `:`
	const answer_prefix_modifier string = `#answer_prefix#%(#data_field#)s`

	//--print before_dl:"Channel: %(channel)s"
	if answer_prefix != "" && input != "" {
		//get data-field
		result := strings.Split(input, postprocessor_split)
		data_field := result[1]
		postprocessor_value := strings.ReplaceAll(result[0], Print, Space)

		//apply answer-prefix
		input = strings.ReplaceAll(answer_prefix_modifier, "#data_field#", data_field)
		input = strings.ReplaceAll(input, "#answer_prefix#", answer_prefix)

		if !remove_print {
			input = Print + input
		}

		if !remove_postprocessor {
			input = strings.ReplaceAll(input, Print, Print+postprocessor_value+postprocessor_split)
		}
	} else if remove_print {
		input = strings.ReplaceAll(input, Print, Space)
	} else if remove_print && remove_postprocessor {
		result := strings.Split(input, postprocessor_split)
		input = result[1]
	} else if remove_postprocessor {
		result := strings.Split(input, postprocessor_split)
		input = Print + result[1]
	}

	return input
}

// Parent Directory
const mediaDirectory string = `-P {{MediaDir}}`

func GetMediaDirectory(keepParentDirectoryFlag bool) string {

	mediaDir := strings.ReplaceAll(mediaDirectory, "{{MediaDir}}", c.Config("MEDIA_PATH", true))
	if !keepParentDirectoryFlag {
		mediaDir = strings.ReplaceAll(mediaDir, "-P ", "")
	}

	return mediaDir
}

func GetCommandName() string {
	result := ""
	if runtime.GOOS == "windows" {
		result = c.Config("YTDLP_NAME_WINDOWS", false)
	}

	if runtime.GOOS == "linux" {
		result = c.Config("YTDLP_NAME_LINUX", false)
	}
	return result
}
