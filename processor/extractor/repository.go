package extractor

import (
	"database/sql"
	"fmt"
	"time"

	p "github.com/rs-anantmishra/streamsphere/utils/processor/database/queries"
	e "github.com/rs-anantmishra/streamsphere/utils/processor/domain"
)

type IRepository interface {
	SaveMetadata(e.MediaInformation) e.SavedInfo
	SaveThumbnail(e.Files) []int
	SaveSubtitles(e.Files) int
	SaveMediaContent(e.Files) int
	GetVideoFileInfo(int) (e.SavedInfo, e.Filepath, error)
	// GetQueuedVideoDetails(videoId int) (e.MinimalCardsInfo, error)
	SavePlaylist(e.Playlist) int
	GetVideoIdByContentId(contentId string) (int, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(Database *sql.DB) IRepository {
	return &repository{
		db: Database,
	}
}

func (r *repository) SaveMetadata(metadata e.MediaInformation) e.SavedInfo {

	var savedInfo e.SavedInfo
	elem := metadata

	//Channel will be same for all items in playlist.
	channelId := genericCheck(*r, elem.ChannelId, "Channel", p.InsertChannelCheck)
	if channelId <= 0 {
		var args []any
		args = append(args, elem.Channel)
		args = append(args, elem.ChannelFollowerCount)
		args = append(args, elem.ChannelURL)
		args = append(args, elem.ChannelId)
		args = append(args, time.Now().Unix())

		channelId = genericSave(*r, args, p.InsertChannel)
		_ = channelId
	}

	//playlist info will be same for all in the playlist.
	//playlistChannelId := getPlaylistChannelId(elem.ChannelId, isSingleChannelPl, elem.PlaylistId) //check why this is needed and if this is numeric or yt id
	playlistId := genericCheck(*r, elem.YoutubePlaylistId, "Playlist", p.InsertPlaylistCheck)
	if playlistId <= 0 && (elem.PlaylistTitle != "" && elem.PlaylistCount > 0) {
		var args []any
		args = append(args, elem.PlaylistTitle)
		args = append(args, elem.PlaylistCount)
		args = append(args, elem.PlaylistChannel)
		args = append(args, elem.PlaylistChannelId)
		args = append(args, elem.PlaylistUploader)
		args = append(args, elem.PlaylistUploaderId)
		args = append(args, 0)
		args = append(args, elem.YoutubePlaylistId)
		args = append(args, time.Now().Unix())

		playlistId = genericSave(*r, args, p.InsertPlaylist)
		_ = playlistId
	}

	//Domain will be same for all items in playlist.
	domainId := genericCheck(*r, elem.Domain, "Domain", p.InsertDomainCheck)
	if domainId <= 0 {
		var args []any
		args = append(args, elem.Domain)
		args = append(args, time.Now().Unix())

		domainId = genericSave(*r, args, p.InsertDomain)
		_ = domainId
	}

	//Format will NOT be same for all items in playlist.
	formatId := genericCheck(*r, elem.Format, "Format", p.InsertFormatCheck)
	if formatId <= 0 {
		var args []any
		args = append(args, elem.Format)
		args = append(args, elem.FormatNote)
		args = append(args, elem.Resolution)
		args = append(args, "Video") //It should be Audio for audio only files
		args = append(args, time.Now().Unix())

		formatId = genericSave(*r, args, p.InsertFormat)
		_ = formatId
	}

	ytVideoId := genericCheck(*r, elem.YoutubeVideoId, "Metadata", p.InsertMetadataCheck)
	if ytVideoId < 0 {
		var args []any

		args = append(args, elem.Title)
		args = append(args, elem.Description)
		args = append(args, elem.Duration)
		args = append(args, elem.OriginalURL)
		args = append(args, elem.WebpageURL)
		args = append(args, elem.LiveStatus)
		args = append(args, elem.Availability)
		args = append(args, elem.YoutubeViewCount)
		args = append(args, elem.LikeCount)
		args = append(args, elem.DislikeCount)
		args = append(args, elem.License)
		args = append(args, elem.AgeLimit)
		args = append(args, elem.PlayableInEmbed)
		args = append(args, elem.UploadDate)
		args = append(args, elem.ReleaseTimestamp)
		args = append(args, elem.ModifiedTimestamp)
		args = append(args, 0) //IsFileDownloaded
		args = append(args, 0) //FileId
		args = append(args, channelId)
		args = append(args, domainId)
		args = append(args, formatId)
		args = append(args, elem.YoutubeVideoId)
		args = append(args, 0)                 //WatchCount
		args = append(args, 0)                 //IsDeleted
		args = append(args, time.Now().Unix()) //CreatedDate

		ytVideoId = genericSave(*r, args, p.InsertMetadata)
	}

	//check can be used for any one-to-many tables
	playlistVideoId := tagsOrCategoriesCheck(*r, playlistId, "PlaylistVideos", p.InsertPlaylistVideosCheck, ytVideoId)
	if playlistVideoId <= 0 && playlistId > 0 {
		var args []any
		args = append(args, ytVideoId)
		args = append(args, playlistId)
		args = append(args, elem.PlaylistVideoIndex)
		args = append(args, time.Now().Unix())

		playlistVideoId = genericSave(*r, args, p.InsertPlaylistVideos)
		_ = playlistVideoId
	}

	//if more items are added to playlist below logic will update the playlist
	if playlistId > 0 && ytVideoId > 0 {
		//Update Total Videos in Playlist
		//update bindings in tblVideos -- order of arguments is important.
		var argsUpdateTotalItems []any
		argsUpdateTotalItems = append(argsUpdateTotalItems, elem.PlaylistCount)
		argsUpdateTotalItems = append(argsUpdateTotalItems, playlistId)
		rowsAffectedTotalItems := genericUpdate(*r, argsUpdateTotalItems, p.UpdatePlaylistItemCount)
		_ = rowsAffectedTotalItems // can do something with it?

		//Update Video Index
		var argsUpdateItemIndex []any
		argsUpdateItemIndex = append(argsUpdateItemIndex, elem.PlaylistVideoIndex)
		argsUpdateItemIndex = append(argsUpdateItemIndex, playlistId)
		argsUpdateItemIndex = append(argsUpdateItemIndex, ytVideoId)
		rowsAffectedItemIndex := genericUpdate(*r, argsUpdateItemIndex, p.UpdatePlaylistVideoIndex)
		_ = rowsAffectedItemIndex // can do something with it?
	}

	//Tags will NOT be same for all items in playlist.
	var lstTagId []int
	for _, element := range elem.Tags {
		tagId := genericCheck(*r, element, "Tag", p.InsertTagsCheck)
		if tagId <= 0 {
			var args []any
			args = append(args, element)
			args = append(args, 1) // IsUsed
			args = append(args, time.Now().Unix())

			tagId = genericSave(*r, args, p.InsertTags)
			_ = tagId
			//here, we should have a map between tag Id,
			//VideoId which should be used to populate VideoFileTags
			//So this would be like a User Defined Type in MSSQL which can
			//be sent at once to SQLite
			lstTagId = append(lstTagId, tagId)
			_ = lstTagId
		}

		videoFileTagId := tagsOrCategoriesCheck(*r, tagId, "VideoFileTag", p.InsertVideoFileTagsCheck, ytVideoId)
		if videoFileTagId < 0 {
			var args []any
			args = append(args, tagId)
			args = append(args, ytVideoId)
			args = append(args, 0)
			args = append(args, time.Now().Unix())
			videoFileTagId = genericSave(*r, args, p.InsertVideoFileTags)
			_ = videoFileTagId
		}
	}

	//Categories will NOT be same for all items in playlist.
	var lstCategoryId []int
	for _, element := range elem.Categories {
		categoryId := genericCheck(*r, element, "Category", p.InsertCategoriesCheck)
		if categoryId <= 0 {
			var args []any
			args = append(args, element)
			args = append(args, 1) // IsUsed
			args = append(args, time.Now().Unix())

			categoryId = genericSave(*r, args, p.InsertCategories)
			_ = categoryId
			//here, we should have a map between categoryId Id,
			//VideoId which should be used to populate VideoFileCategories
			//So this would be like a User Defined Type in MSSQL which can
			//be sent at once to SQLite
			lstCategoryId = append(lstCategoryId, categoryId)
			_ = lstCategoryId
		}

		videoFileCategoryId := tagsOrCategoriesCheck(*r, categoryId, "VideoFileCategory", p.InsertVideoFileCategoriesCheck, ytVideoId)
		if videoFileCategoryId < 0 {
			var args []any
			args = append(args, categoryId)
			args = append(args, ytVideoId)
			args = append(args, 0)
			args = append(args, time.Now().Unix())
			videoFileCategoryId = genericSave(*r, args, p.InsertVideoFileCategories)
			_ = videoFileCategoryId
		}
	}

	//complete result
	savedInfo = e.SavedInfo{
		VideoId:        ytVideoId,
		YoutubeVideoId: elem.YoutubeVideoId,
		PlaylistId:     playlistId,
		ChannelId:      channelId,
		DomainId:       domainId,
		FormatId:       formatId,
		MediaInfo:      elem,
	}

	return savedInfo
}

func (r *repository) SaveThumbnail(file e.Files) []int {

	var lstFileIds []int

	thumbnailFileId := filesCheck(*r, file.FileType, file.VideoId, p.InsertThumbnailFileCheck)
	if thumbnailFileId <= 0 {
		var args []any
		args = append(args, file.VideoId)
		args = append(args, file.FileType)
		args = append(args, file.SourceId)
		args = append(args, file.FilePath)
		args = append(args, file.FileName)
		args = append(args, file.Extension)
		args = append(args, file.FileSize)
		args = append(args, file.FileSizeUnit)
		args = append(args, file.NetworkPath)
		args = append(args, file.IsDeleted)
		args = append(args, time.Now().Unix())

		thumbnailFileId = genericSave(*r, args, p.InsertFile)
		lstFileIds = append(lstFileIds, thumbnailFileId)
		_ = thumbnailFileId
	}

	return lstFileIds
}

func (r *repository) SaveSubtitles(file e.Files) int {

	subsFileId := subsFilesCheck(*r, file.FileType, file.VideoId, file.FileName, p.InsertSubsFileCheck)
	if subsFileId <= 0 {
		var args []any
		args = append(args, file.VideoId)
		args = append(args, file.FileType)
		args = append(args, file.SourceId)
		args = append(args, file.FilePath)
		args = append(args, file.FileName)
		args = append(args, file.Extension)
		args = append(args, file.FileSize)
		args = append(args, file.FileSizeUnit)
		args = append(args, file.NetworkPath)
		args = append(args, file.IsDeleted)
		args = append(args, time.Now().Unix())

		subsFileId = genericSave(*r, args, p.InsertFile)
		_ = subsFileId
	}

	return subsFileId
}

func (r *repository) SaveMediaContent(elem e.Files) int {

	var fileId int

	mediaFileId := filesCheck(*r, elem.FileType, elem.VideoId, p.InsertThumbnailFileCheck)
	if mediaFileId <= 0 {
		var args []any
		args = append(args, elem.VideoId)
		args = append(args, elem.FileType)
		args = append(args, elem.SourceId)
		args = append(args, elem.FilePath)
		args = append(args, elem.FileName)
		args = append(args, elem.Extension)
		args = append(args, elem.FileSize)
		args = append(args, elem.FileSizeUnit)
		args = append(args, elem.NetworkPath)
		args = append(args, elem.IsDeleted)
		args = append(args, time.Now().Unix())

		mediaFileId = genericSave(*r, args, p.InsertFile)
		fileId = mediaFileId
		_ = mediaFileId

		//update bindings in tblVideos -- order of arguments is important.
		var argsUpdate []any
		argsUpdate = append(argsUpdate, 1)
		argsUpdate = append(argsUpdate, mediaFileId)
		argsUpdate = append(argsUpdate, elem.VideoId)
		rowsAffected := genericUpdate(*r, argsUpdate, p.UpdateVideoFileFields)
		_ = rowsAffected // can do something with it?

		//update bindings in tblPlaylistVideoFiles
		//THERE SHOULD BE A PLAYLISTID ALSO IN WHERE CLAUSE
		var argsPVFUpdate []any
		argsPVFUpdate = append(argsPVFUpdate, mediaFileId)
		argsPVFUpdate = append(argsPVFUpdate, elem.VideoId)
		rowsAffectedPVF := genericUpdate(*r, argsPVFUpdate, p.UpdatePVFFileId)
		_ = rowsAffectedPVF // can do something with it?
	}

	return fileId
}

func (r *repository) GetVideoFileInfo(videoId int) (e.SavedInfo, e.Filepath, error) {

	var smi e.SavedInfo
	var fPath e.Filepath

	smi.VideoId = videoId
	row := r.db.QueryRow(p.GetVideoInformationById, videoId)
	if err := row.Scan(&smi.MediaInfo.Title, &smi.MediaInfo.PlaylistTitle, &fPath.Channel, &fPath.Domain, &fPath.PlaylistTitle, &smi.YoutubeVideoId, &smi.MediaInfo.WebpageURL); err != nil {
		if err == sql.ErrNoRows {
			return smi, fPath, fmt.Errorf("VideoId %d: no such video", videoId)
		}
		return smi, fPath, fmt.Errorf("VideoById %d: %v", videoId, err)
	}

	smi.MediaInfo.Domain = fPath.Domain
	smi.MediaInfo.Channel = fPath.Channel

	return smi, fPath, nil
}

// func (r *repository) GetQueuedVideoDetails(videoId int) (e.MinimalCardsInfo, error) {

// 	var minInfo e.MinimalCardsInfo

// 	minInfo.VideoId = videoId
// 	row := r.db.QueryRow(p.GetQueuedVideoDetailsById, videoId)
// 	if err := row.Scan(&minInfo.VideoId, &minInfo.Title, &minInfo.Channel, &minInfo.Description, &minInfo.Duration, &minInfo.WebpageURL, &minInfo.Thumbnail); err != nil {
// 		if err == sql.ErrNoRows {
// 			return minInfo, fmt.Errorf("VideoId %d: no such video", videoId)
// 		}
// 		return minInfo, fmt.Errorf("VideoById %d: %v", videoId, err)
// 	}
// 	return minInfo, nil
// }

// region Private Methods ////////////////////////////////////////////////////

func genericSave(r repository, args []any, genericQuery string) int {
	resultId := -1

	if resultId < 0 {

		//generic save
		result, err := r.db.Exec(genericQuery, args...)

		//check for errors
		if err != nil {
			fmt.Println("error:", err)
			return resultId
		}

		//get the inserted records Id
		var id int64
		if id, err = result.LastInsertId(); err != nil {
			return resultId
		}
		resultId = int(id)
	}

	return resultId
}

func genericUpdate(r repository, args []any, genericQuery string) int {
	rowsAffected := -1

	if rowsAffected < 0 {

		//generic save
		result, err := r.db.Exec(genericQuery, args...)

		//check for errors
		if err != nil {
			fmt.Println("error:", err)
			return rowsAffected
		}

		//get the inserted records Id
		var id int64
		if id, err = result.RowsAffected(); err != nil {
			return rowsAffected
		}
		rowsAffected = int(id)
	}

	return rowsAffected
}

func genericCheck(r repository, Id any, idContext string, genericQuery string) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(genericQuery, Id)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for", idContext, "Id", Id)
	}

	return resultId
}

func tagsOrCategoriesCheck(r repository, Id any, idContext string, genericQuery string, videoId int) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(genericQuery, Id, videoId)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for", idContext, "Id", Id)
	}

	return resultId
}

func filesCheck(r repository, fileType string, videoId int, filesCheckQuery string) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(filesCheckQuery, fileType, videoId)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for FileType:", fileType, "VideoId:", videoId)
	}

	return resultId
}

func subsFilesCheck(r repository, fileType string, videoId int, filename string, filesCheckQuery string) int {
	resultId := -1

	//Check if entry for Id exists?
	chk := r.db.QueryRow(filesCheckQuery, fileType, videoId, filename)
	if err := chk.Scan(&resultId); err == sql.ErrNoRows {
		fmt.Println("No Rows found for FileType:", fileType, "VideoId:", videoId, "FileName:", filename)
	}

	return resultId
}

// endregion Private Methods ////////////////////////////////////////////////////

func (r *repository) SavePlaylist(playlist e.Playlist) int {

	playlistId := genericCheck(*r, playlist.YoutubePlaylistId, "Playlist", p.InsertPlaylistCheck)
	if playlistId <= 0 && (playlist.Title != "" && playlist.ItemCount > 0) {
		var args []any
		args = append(args, playlist.Title)
		args = append(args, playlist.ItemCount)
		args = append(args, playlist.PlaylistChannel)
		args = append(args, playlist.PlaylistChannelId)
		args = append(args, playlist.PlaylistUploader)
		args = append(args, playlist.PlaylistUploaderId)
		args = append(args, 0)
		args = append(args, playlist.YoutubePlaylistId)
		args = append(args, time.Now().Unix())

		playlistId = genericSave(*r, args, p.InsertPlaylist)
		_ = playlistId
	}

	return playlistId
}

func (r *repository) GetVideoIdByContentId(contentId string) (int, int, error) {

	var id int
	var channelId int
	row := r.db.QueryRow(p.GetVideoIdByContentId, contentId)
	if err := row.Scan(&id, &channelId); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, fmt.Errorf("content-id %v: no such video", contentId)
		}
		return -1, -1, fmt.Errorf("content-id %v: %v", contentId, err)
	}
	return id, channelId, nil
}
