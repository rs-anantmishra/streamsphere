package queries

// All Videos Page Queries /////////////////////////////////////////////////////////
// const GetVideoMetadata_AllVideos string = `Select DISTINCT V.Id, V.Title, V.Description, V.DurationSeconds, V.OriginalURL, V.WebpageURL, V.IsFileDownloaded, V.IsDeleted, C.Name, V.LiveStatus, 
// D.Domain, V.LikeCount, V.YoutubeViewCount as 'ViewsCount', V.WatchCount, V.UploadDate, V.Availability, F.Format, V.YoutubeVideoId, V.CreatedDate
// FROM tblVideos V
// INNER JOIN tblChannels C ON V.ChannelId = C.Id
// INNER JOIN tblPlaylistVideoFiles PVF ON V.Id = PVF.VideoId
// INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
// INNER JOIN tblDomains D ON V.DomainId = D.Id
// INNER JOIN tblFormats F ON V.FormatId = F.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') ORDER BY V.Id ASC;`

// const GetVideoTags_AllVideos string = `Select V.Id, T.Name, T.IsUsed, T.CreatedDate
// FROM tblVideos V
// INNER JOIN tblVideoFileTags VFT ON V.Id = VFT.VideoId
// INNER JOIN tblTags T ON T.Id = VFT.TagId;`

// const GetVideoCategories_AllVideos string = `Select V.Id, C.Name, C.IsUsed, C.CreatedDate
// FROM tblVideos V
// INNER JOIN tblVideoFileCategories VFC ON V.Id = VFC.VideoId
// INNER JOIN tblCategories C ON C.Id = VFC.CategoryId;`

// const GetVideoFiles_AllVideos string = `Select V.Id, F.FileType, F.FileSize, F.Extension, F.FilePath, F.FileName
// FROM tblVideos V
// INNER JOIN tblFiles F ON V.Id = F.VideoId;`

// const GetVideoMetadata string = `Select DISTINCT V.Id, V.Title, V.Description, V.DurationSeconds, V.OriginalURL, V.WebpageURL, V.IsFileDownloaded, V.IsDeleted, C.Name, V.LiveStatus, 
// D.Domain, V.LikeCount, V.YoutubeViewCount as 'ViewsCount', V.WatchCount, V.UploadDate, V.Availability, F.Format, V.YoutubeVideoId, V.CreatedDate
// FROM tblVideos V
// INNER JOIN tblChannels C ON V.ChannelId = C.Id
// INNER JOIN tblPlaylistVideoFiles PVF ON V.Id = PVF.VideoId
// INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
// INNER JOIN tblDomains D ON V.DomainId = D.Id
// INNER JOIN tblFormats F ON V.FormatId = F.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video')
// WHERE V.Id = ? ORDER BY V.Id ASC;`

// const GetVideoTags string = `Select V.Id, T.Name, T.IsUsed, T.CreatedDate
// FROM tblVideos V
// INNER JOIN tblVideoFileTags VFT ON V.Id = VFT.VideoId
// INNER JOIN tblTags T ON T.Id = VFT.TagId
// WHERE V.Id = ?;`

// const GetVideoCategories string = `Select V.Id, C.Name, C.IsUsed, C.CreatedDate
// FROM tblVideos V
// INNER JOIN tblVideoFileCategories VFC ON V.Id = VFC.VideoId
// INNER JOIN tblCategories C ON C.Id = VFC.CategoryId
// WHERE V.Id = ?;`

// const GetVideoFiles string = `Select V.Id, F.FileType, F.FileSize, F.Extension, F.FilePath, F.FileName
// FROM tblVideos V
// INNER JOIN tblFiles F ON V.Id = F.VideoId
// WHERE V.Id = ?;`

// const GetQueuedVideoDetailsById string = `Select V.Id, V.Title, C.Name as 'Channel', V.Description, V.DurationSeconds as 'Duration', V.WebpageURL, FIThumbnail.FilePath || '\' || FIThumbnail.FileName as 'Thumbnail'
// FROM tblVideos V
// INNER JOIN tblChannels C ON V.ChannelId = C.Id
// INNER JOIN tblPlaylists P ON V.PlaylistId = P.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// WHERE V.Id = ?`

// const GetVideoSearchChannelTitle string = `Select DISTINCT V.Id 'VideoId', V.Title, C.Name as 'Channel'
// FROM tblVideos V
// INNER JOIN tblChannels C ON V.ChannelId = C.Id
// INNER JOIN tblPlaylistVideoFiles PVF ON V.Id = PVF.VideoId
// INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
// INNER JOIN tblDomains D ON V.DomainId = D.Id
// INNER JOIN tblFormats F ON V.FormatId = F.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') ORDER BY V.Id ASC;`

// All Videos Page Queries /////////////////////////////////////////////////////////

// Playlists Page Queries /////////////////////////////////////////////////////////
// const GetAllPlaylists string = `SELECT DISTINCT P.Id, P.Title, P.PlaylistUploader, P.ItemCount, P.YoutubePlaylistId, F.FilePath || '\' || F.FileName as 'ThumbnailURL'
// FROM tblPlaylists P
// INNER JOIN tblPlaylistVideoFiles PVF ON PVF.PlaylistId = P.Id
// INNER JOIN tblVideos V ON V.Id = PVF.VideoId
// INNER JOIN tblFiles F ON F.VideoId = PVF.VideoId
// WHERE P.Id > 0 
// 	AND PVF.PlaylistVideoIndex = 1 
// 	AND F.FileType = 'Thumbnail'
// ORDER BY P.Id ASC`

// const GetVideoMetadata_Playlists string = `Select DISTINCT V.Id, V.Title, V.Description, V.DurationSeconds, V.OriginalURL, V.WebpageURL, V.IsFileDownloaded, V.IsDeleted, C.Name, V.LiveStatus, 
// D.Domain, V.LikeCount, V.YoutubeViewCount as 'ViewsCount', V.WatchCount, V.UploadDate, V.Availability, F.Format, V.YoutubeVideoId, V.CreatedDate, PVF.PlaylistVideoIndex
// FROM tblVideos V
// INNER JOIN tblPlaylistVideoFiles PVF ON (V.Id = PVF.VideoId AND PVF.PlaylistId = ?)
// INNER JOIN tblChannels C ON V.ChannelId = C.Id
// INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
// INNER JOIN tblDomains D ON V.DomainId = D.Id
// INNER JOIN tblFormats F ON V.FormatId = F.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') 
// ORDER BY V.Id ASC;`

// const GetVideoTags_Playlists string = `Select V.Id, T.Name, T.IsUsed, T.CreatedDate
// FROM tblVideos V
// INNER JOIN tblPlaylistVideoFiles PVF ON (PVF.VideoId = V.Id AND PVF.PlaylistId = ?)
// INNER JOIN tblVideoFileTags VFT ON V.Id = VFT.VideoId
// INNER JOIN tblTags T ON T.Id = VFT.TagId;`

// const GetVideoCategories_Playlists string = `Select V.Id, C.Name, C.IsUsed, C.CreatedDate
// FROM tblVideos V
// INNER JOIN tblPlaylistVideoFiles PVF ON (PVF.VideoId = V.Id AND PVF.PlaylistId = ?)
// INNER JOIN tblVideoFileCategories VFC ON V.Id = VFC.VideoId
// INNER JOIN tblCategories C ON C.Id = VFC.CategoryId;`

// const GetVideoFiles_Playlists string = `Select V.Id, F.FileType, F.FileSize, F.Extension, F.FilePath, F.FileName
// FROM tblVideos V
// INNER JOIN tblPlaylistVideoFiles PVF ON (PVF.VideoId = V.Id AND PVF.PlaylistId = ?)
// INNER JOIN tblFiles F ON V.Id = F.VideoId;`

// Playlists Page Queries /////////////////////////////////////////////////////////

// Channels Page Queries /////////////////////////////////////////////////////////
// const GetVideoMetadata_Channels string = `Select DISTINCT V.Id, V.Title, V.Description, V.DurationSeconds, V.OriginalURL, V.WebpageURL, V.IsFileDownloaded, V.IsDeleted, C.Name, V.LiveStatus, 
// D.Domain, V.LikeCount, V.YoutubeViewCount as 'ViewsCount', V.WatchCount, V.UploadDate, V.Availability, F.Format, V.YoutubeVideoId, V.CreatedDate
// FROM tblVideos V
// INNER JOIN tblChannels C ON (V.ChannelId = C.Id AND C.Id = ?)
// INNER JOIN tblPlaylistVideoFiles PVF ON V.Id = PVF.VideoId
// INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
// INNER JOIN tblDomains D ON V.DomainId = D.Id
// INNER JOIN tblFormats F ON V.FormatId = F.Id
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') 
// ORDER BY V.Id ASC;`

// const GetVideoTags_Channels string = `Select V.Id, T.Name, T.IsUsed, T.CreatedDate
// FROM tblVideos V
// INNER JOIN tblChannels C ON (V.ChannelId = C.Id AND C.Id = ?)
// INNER JOIN tblVideoFileTags VFT ON V.Id = VFT.VideoId
// INNER JOIN tblTags T ON T.Id = VFT.TagId;`

// const GetVideoCategories_Channels string = `Select V.Id, C.Name, C.IsUsed, C.CreatedDate
// FROM tblVideos V
// INNER JOIN tblChannels Channel ON (V.ChannelId = Channel.Id AND Channel.Id = ?)
// INNER JOIN tblVideoFileCategories VFC ON V.Id = VFC.VideoId
// INNER JOIN tblCategories C ON C.Id = VFC.CategoryId;`

// const GetVideoFiles_Channels string = `Select V.Id, F.FileType, F.FileSize, F.Extension, F.FilePath, F.FileName
// FROM tblVideos V
// INNER JOIN tblChannels C ON (V.ChannelId = C.Id AND C.Id = ?)
// INNER JOIN tblFiles F ON V.Id = F.VideoId;`

//Channels Page Queries /////////////////////////////////////////////////////////

const GetVideoById string = `Select DISTINCT V.Id, V.Title, V.Description, V.DurationSeconds, V.OriginalURL, V.WebpageURL, V.IsFileDownloaded, V.IsDeleted, C.Name, V.LiveStatus, 
D.Domain, V.LikeCount, V.YoutubeViewCount as 'ViewsCount', V.WatchCount, V.UploadDate, V.Availability, F.Format, V.YoutubeVideoId, V.CreatedDate
FROM tblVideos V
INNER JOIN tblPlaylistVideoFiles PVF ON (V.Id = PVF.VideoId AND V.Id = ?)
INNER JOIN tblChannels C ON V.ChannelId = C.Id
INNER JOIN tblPlaylists P ON PVF.PlaylistId = P.Id
INNER JOIN tblDomains D ON V.DomainId = D.Id
INNER JOIN tblFormats F ON V.FormatId = F.Id
INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') ORDER BY V.Id ASC;`

/*Hard Delete Queries*/

// const GetVideoDeleteDetails string = `Select YoutubeVideoId, FIThumbnail.FilePath || '\' || FIThumbnail.FileName as 'ThumbnailFile', FIVideo.FilePath || '\' || FIVideo.FileName as 'ContentFile'
// FROM tblVideos V
// INNER JOIN tblFiles FIThumbnail ON (V.Id = FIThumbnail.VideoId AND FIThumbnail.FileType = 'Thumbnail')
// INNER JOIN tblFiles FIVideo ON (V.Id = FIVideo.VideoId AND FIVideo.FileType = 'Video') 
// WHERE V.Id = ?`

// const DeleteVideoFileTags string = `DELETE FROM tblVideoFileTags WHERE VideoId = ?`
// const DeleteVideoFileCategories string = `DELETE FROM tblVideoFileCategories WHERE VideoId = ?`
// const DeletePlaylistVideoFiles string = `DELETE FROM tblPlaylistVideoFiles WHERE VideoId = ?`
// const DeleteFiles string = `DELETE FROM tblFiles WHERE VideoId = ?`
// const DeleteVideos string = `DELETE FROM tblVideos WHERE Id = ?`

/*Hard Delete Queries*/
