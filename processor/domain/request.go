package domain

// Populated by Streamsphere API ////////////////////////////////
// ContentFormat will come from ui, possible values:
// Audio-Best, Video-Best, FormatId
type Request struct {
	Id                int
	RequestUrl        string
	RequestType       string
	Metadata          int
	Thumbnail         int
	Content           int
	ContentFormat     string // will be used for video format
	Subtitles         int
	SubtitlesLanguage string //default
	IsProxied         int
	Proxy             string
	Scheduled         int //Schedule this request?
	CreatedDate       int
	ModifiedDate      int
}

//Populated by Processor CLI
type RequestStatus struct {
	Id            int
	RequestId     int
	RequestStatus string
	CreatedDate   int
	ModifiedDate  int
}

const RequestStatus_Recieved string = `Recieved`               //Set by API, request is saved and this entry is made in tblRequestStatus
const RequestStatus_Analyzing string = `Analyzing`             //Request expansion has started
const RequestStatus_Queued string = `Queued`                   //Request expansion is complete and expanded requests are queued
const RequestStatus_ProcessQueue string = `ProcessQueue`       //Queued requests are being processed
const RequestStatus_PartialComplete string = `PartialComplete` //completed status - Some Videos in playlist/channel failed to download (except private videos)
const RequestStatus_Complete string = `Complete`               //completed status - Successful
const RequestStatus_Failed string = `Failed`                   //completed status - Failed at 'Analyzing' or 'Post Queued'

//Populated by Processor CLI
type RequestQueue struct {
	Id            int
	RequestId     int
	ContentId     string
	ProcessStatus string
	RetryCount    int
	Message       string
	Cancelled     Bool
	RequestType   string
	CreatedDate   int64
	ModifiedDate  int64
}

const ProcessStatus_Queued string = `Queued`                          // Item is queued for processing
const ProcessStatus_ProcessMetadata string = `Downloading Metadata`   // Downloading Metadata
const ProcessStatus_ProcessThumbnail string = `Downloading Thumbnail` // Downloading Thumbnail
const ProcessStatus_ProcessContent string = `Downloading Content`     // Downloading Content
const ProcessStatus_ProcessSubs string = `Downloading Subtitles`      // Downloading Subtitles
const ProcessStatus_Complete string = `Complete`                      // Completed
const ProcessStatus_Failed string = `Failed`                          // Completed Status - Failed due to some reason

//-- Message --//
const Message_Success string = `Completed Successfully` // Success Message
const Message_Failed string = ``                        // Get Failure reason
//-- Message --//

//-- Populates tblPlaylists --//
type ChannelPlaylists struct {
	ChannelId string
	Playlist  []Playlist
}

//-- Populated tblPVF --//
type PlaylistContent struct {
	PlaylistId string
	Content    []PlaylistContentMeta
}

type PlaylistContentMeta struct {
	ContentId            string `json:"id"`
	PlaylistContentIndex int    `json:"playlist_index"`
}

type PlaylistUploader struct {
	PlaylistUploaderId string `json:"playlist_uploader_id"`
	PlaylistId         string `json:"playlist_id"`
}

type FilenameInfo struct {
	ContentId    string
	Id           int
	PlaylistId   int    //this may need to contain channelId going onwards -- to be verified.
	Domain       string `json:"webpage_url_domain"`
	Channel      string `json:"channel"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail"`
}

type Process struct {
	Id          int
	ProcessId   int //(pid for extractor, populated by API)
	RequestId   int //(Current RequestId that is being processed)
	Status      int //(boolean)[0 - Stopped, 1 - Processing]
	StartTime   int
	EndTime     int
	ProcessType string //(api-request-process/scheduled-request-process)
	Message     string //(If Panic is there, recover and log, if same repeats, log and exit)
	CreatedDate int
}

//1. Get all playlists for channel, then -
//2. Populate tblPlaylists with below (except itemcount)
//yt-dlp_x86.exe --print %(.{id})s --print %(.{title})s --print %(.{playlist_channel})s --flat-playlist --lazy-playlist "https://www.youtube.com/@sigfaults/playlists"  --print %(.{playlist_channel_id})s --print %(.{playlist_uploader})s --print %(.{playlist_uploader_id})s

//3. Get all videos for playlist
//yt-dlp_x86.exe --print %(.{id})s --print %(.{playlist_index})s --flat-playlist --lazy-playlist "PLBan2kCeFnBosc-AMMGXLXPKYDTA-7t_a"

//4. Dump contentId's into tblRequestQueue - NEXT > fetch metadata
