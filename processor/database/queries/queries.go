package queries

// This file to be replaced with .sql files for each query
const InsertChannelCheck string = `Select Id From tblChannels Where YoutubeChannelId = ?`
const InsertChannel string = `INSERT INTO tblChannels Select NULL, ?, ?, ?, ?, ?;`

const InsertPlaylistCheck string = `Select Id From tblPlaylists WHERE YoutubePlaylistId = ?`
const InsertPlaylist string = `INSERT INTO tblPlaylists Select NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?;`
const UpdatePlaylistItemCount string = `UPDATE tblPlaylists SET ItemCount = ? WHERE Id = ?`
const UpdatePlaylistVideoIndex string = `UPDATE tblPlaylistVideoFiles SET PlaylistVideoIndex = ? WHERE PlaylistId = ? AND VideoId = ?`

const InsertPlaylistVideosCheck string = `Select Id From tblPlaylistVideoFiles WHERE PlaylistId = ? AND VideoId = ?`
const InsertPlaylistVideos string = `INSERT INTO tblPlaylistVideoFiles(Id, VideoId, PlaylistId, PlaylistVideoIndex, CreatedDate) Select NULL, ?, ?, ?, ?;`

const InsertDomainCheck string = `Select Id From tblDomains Where Domain = ?`
const InsertDomain string = `INSERT INTO tblDomains Select NULL, ?, ?;`

const InsertFormatCheck string = `Select Id From tblFormats Where Format = ?`
const InsertFormat string = `INSERT INTO tblFormats Select NULL, ?, ?, ?, ?, ?;`

const InsertMetadataCheck string = `Select Id From tblVideos Where YoutubeVideoId = ?`
const InsertMetadata string = `INSERT INTO tblVideos (
	Id
	,Title
	,Description
	,DurationSeconds
	,OriginalURL
	,WebpageURL
	,LiveStatus
	,Availability
	,YoutubeViewCount
	,LikeCount
	,DislikeCount
	,License
	,AgeLimit
	,PlayableInEmbed
	,UploadDate
	,ReleaseTimestamp
	,ModifiedTimestamp
	,IsFileDownloaded
	,FileId
	,ChannelId
	,DomainId
	,FormatId
	,YoutubeVideoId
	,WatchCount
	,IsDeleted
	,CreatedDate
)
VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const InsertTagsCheck string = `Select Id From tblTags Where Name = ?`
const InsertTags string = `INSERT INTO tblTags Select NULL, ?, ?, ?;`

const InsertCategoriesCheck string = `Select Id From tblCategories Where Name = ?`
const InsertCategories string = `INSERT INTO tblCategories Select NULL, ?, ?, ?;`

const InsertVideoFileTagsCheck string = `SELECT Id From tblVideoFileTags Where TagId = ? AND VideoId = ?`
const InsertVideoFileTags string = `INSERT INTO tblVideoFileTags SELECT NULL, ?, ?, ?, ?`

const InsertVideoFileCategoriesCheck string = `SELECT Id From tblVideoFileCategories Where CategoryId = ? AND VideoId = ?`
const InsertVideoFileCategories string = `INSERT INTO tblVideoFileCategories SELECT NULL, ?, ?, ?, ?`

const InsertThumbnailFileCheck string = `SELECT Id From tblFiles WHERE FileType = ? AND VideoId = ?`
const InsertSubsFileCheck string = `SELECT Id From tblFiles WHERE FileType = ? AND VideoId = ? AND FileName = ?`
const InsertFile string = `INSERT INTO tblFiles SELECT NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?`

//Files with duplicate names to be overwritten or renamed(Choice thru UI).
const InsertMediaFileCheck string = `SELECT Id From tblFiles WHERE FileType = ? AND VideoId = ? AND FileName = ?`

const GetNetworkVideoURLById string = `Select WebpageURL from tblVideos Where Id = ?`
const GetVideoInformationById string = `Select V.Title, P.Title, C.Name as 'Channel', D.Domain, P.Title as 'PlaylistTitle', YoutubeVideoId, V.WebpageURL
										FROM tblVideos V 
										INNER JOIN tblChannels C ON C.Id = V.ChannelId 
										INNER JOIN tblDomains D ON D.Id = V.DomainId
										INNER JOIN tblPlaylistVideoFiles PVF ON PVF.VideoId = V.Id
										INNER JOIN tblPlaylists P ON P.Id = PVF.PlaylistId
										WHERE V.Id = ?;`

const UpdateVideoFileFields string = `UPDATE tblVideos SET IsFileDownloaded = ?, FileId = ? WHERE Id = ?;`
const UpdatePVFFileId string = `UPDATE tblPlaylistVideoFiles SET FileId = ? WHERE VideoId = ?`

//below may not be used
const UpdatePVFThumbnailFileId string = `Select F.Id
										 From tblFiles F
										 INNER JOIN tblPlaylistVideoFiles PVF ON PVF.VideoId = F.VideoId
										 WHERE PVF.PlaylistId = ? AND F.FileType = 'Thumbnail'
										 AND F.VideoId = ?`

const GetVideoIdByContentId string = `SELECT Id, ChannelId FROM tblVideos WHERE YoutubeVideoId = ?;`
