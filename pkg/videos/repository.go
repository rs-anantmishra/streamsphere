package videos

import (
	"database/sql"
	"fmt"

	q "github.com/rs-anantmishra/metubeplus/database/queries"
	"github.com/rs-anantmishra/metubeplus/pkg/entities"
)

type IRepository interface {
	GetAllVideos() ([]entities.Videos, error)
	GetContentById(int) ([]entities.Videos, error)
	GetPlaylistVideos(int) ([]entities.Videos, error)
	GetAllPlaylists() ([]entities.Playlist, error)
	GetVideoSearchInfo() ([]entities.ContentSearch, error)
}

type repository struct {
	db *sql.DB
}

func NewVideoRepo(Database *sql.DB) IRepository {
	return &repository{
		db: Database,
	}
}

// GetAllVideos implements IRepository.
func (r *repository) GetAllVideos() ([]entities.Videos, error) {

	var lstVideos []entities.Videos

	rows, err := r.db.Query(q.GetVideoMetadata_AllVideos)
	if err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var v entities.Videos
		if err := rows.Scan(&v.Id, &v.Title, &v.Description, &v.DurationSeconds, &v.OriginalURL, &v.WebpageURL, &v.IsFileDownloaded,
			&v.IsDeleted, &v.Channel.Name, &v.LiveStatus, &v.Domain.Domain, &v.LikesCount, &v.ViewsCount, &v.WatchCount, &v.UploadDate,
			&v.Availability, &v.Format.Format, &v.YoutubeVideoId, &v.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching videos: %v", err)
		}
		lstVideos = append(lstVideos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}

	//tags
	tagsRows, err := r.db.Query(q.GetVideoTags_AllVideos)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for Video: %v", err)
	}
	defer tagsRows.Close()

	var lstTags []entities.Tags
	for tagsRows.Next() {
		var tags entities.Tags
		if err := tagsRows.Scan(&tags.Id, &tags.Name, &tags.IsUsed, &tags.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching tags: %v", err)
		}
		lstTags = append(lstTags, tags)
	}

	//assign Tags to video
	for i := range lstVideos {
		for k, elem := range lstTags {
			if lstVideos[i].Id == lstTags[k].Id {
				lstVideos[i].Tags = append(lstVideos[i].Tags, elem)
			}
		}
	}

	//categories
	categoryRows, err := r.db.Query(q.GetVideoCategories_AllVideos)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for Video: %v", err)
	}
	defer categoryRows.Close()

	var lstCategories []entities.Categories
	for categoryRows.Next() {
		var categories entities.Categories
		if err := categoryRows.Scan(&categories.Id, &categories.Name, &categories.IsUsed, &categories.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching categories: %v", err)
		}
		lstCategories = append(lstCategories, categories)
	}

	//assign categories to video
	for i := range lstVideos {
		for k, elem := range lstCategories {
			if lstVideos[i].Id == lstCategories[k].Id {
				lstVideos[i].Categories = append(lstVideos[i].Categories, elem)
			}
		}
	}

	//files
	filesRows, err := r.db.Query(q.GetVideoFiles_AllVideos)
	if err != nil {
		return nil, fmt.Errorf("error fetching files for Video: %v", err)
	}
	defer filesRows.Close()

	var lstFiles []entities.Files
	for filesRows.Next() {
		var files entities.Files
		if err := filesRows.Scan(&files.VideoId, &files.FileType, &files.FileSize, &files.Extension, &files.FilePath, &files.FileName); err != nil {
			return nil, fmt.Errorf("error fetching files: %v", err)
		}
		lstFiles = append(lstFiles, files)
	}

	//assign Files to video
	for i := range lstVideos {
		for k, elem := range lstFiles {
			if lstVideos[i].Id == lstFiles[k].VideoId {
				lstVideos[i].Files = append(lstVideos[i].Files, elem)
			}
		}
	}

	return lstVideos, nil
}

func (r *repository) GetPlaylistVideos(playlistId int) ([]entities.Videos, error) {
	var lstPlaylistVideos []entities.Videos

	rows, err := r.db.Query(q.GetVideoMetadata_Playlists, playlistId)
	if err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var v entities.Videos
		if err := rows.Scan(&v.Id, &v.Title, &v.Description, &v.DurationSeconds, &v.OriginalURL, &v.WebpageURL, &v.IsFileDownloaded,
			&v.IsDeleted, &v.Channel.Name, &v.LiveStatus, &v.Domain.Domain, &v.LikesCount, &v.ViewsCount, &v.WatchCount, &v.UploadDate,
			&v.Availability, &v.Format.Format, &v.YoutubeVideoId, &v.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching videos: %v", err)
		}
		lstPlaylistVideos = append(lstPlaylistVideos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}

	//tags
	tagsRows, err := r.db.Query(q.GetVideoTags_Playlists, playlistId)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for Video: %v", err)
	}
	defer tagsRows.Close()

	var lstTags []entities.Tags
	for tagsRows.Next() {
		var tags entities.Tags
		if err := tagsRows.Scan(&tags.Id, &tags.Name, &tags.IsUsed, &tags.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching tags: %v", err)
		}
		lstTags = append(lstTags, tags)
	}

	//assign Tags to video
	for i := range lstPlaylistVideos {
		for k, elem := range lstTags {
			if lstPlaylistVideos[i].Id == lstTags[k].Id {
				lstPlaylistVideos[i].Tags = append(lstPlaylistVideos[i].Tags, elem)
			}
		}
	}

	//categories
	categoryRows, err := r.db.Query(q.GetVideoCategories_Playlists, playlistId)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for Video: %v", err)
	}
	defer categoryRows.Close()

	var lstCategories []entities.Categories
	for categoryRows.Next() {
		var categories entities.Categories
		if err := categoryRows.Scan(&categories.Id, &categories.Name, &categories.IsUsed, &categories.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching categories: %v", err)
		}
		lstCategories = append(lstCategories, categories)
	}

	//assign categories to video
	for i := range lstPlaylistVideos {
		for k, elem := range lstCategories {
			if lstPlaylistVideos[i].Id == lstCategories[k].Id {
				lstPlaylistVideos[i].Categories = append(lstPlaylistVideos[i].Categories, elem)
			}
		}
	}

	//files
	filesRows, err := r.db.Query(q.GetVideoFiles_Playlists, playlistId)
	if err != nil {
		return nil, fmt.Errorf("error fetching files for Video: %v", err)
	}
	defer filesRows.Close()

	var lstFiles []entities.Files
	for filesRows.Next() {
		var files entities.Files
		if err := filesRows.Scan(&files.VideoId, &files.FileType, &files.FileSize, &files.Extension, &files.FilePath, &files.FileName); err != nil {
			return nil, fmt.Errorf("error fetching files: %v", err)
		}
		lstFiles = append(lstFiles, files)
	}

	//assign Files to video
	for i := range lstPlaylistVideos {
		for k, elem := range lstFiles {
			if lstPlaylistVideos[i].Id == lstFiles[k].VideoId {
				lstPlaylistVideos[i].Files = append(lstPlaylistVideos[i].Files, elem)
			}
		}
	}

	return lstPlaylistVideos, nil
}

func (r *repository) GetAllPlaylists() ([]entities.Playlist, error) {

	//Id, Title, PlaylistUploader, ItemCount
	var lstPlaylists []entities.Playlist

	rows, err := r.db.Query(q.GetAllPlaylists)
	if err != nil {
		return nil, fmt.Errorf("error fetching playlists: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var p entities.Playlist
		if err := rows.Scan(&p.Id, &p.Title, &p.PlaylistUploader, &p.ItemCount, &p.YoutubePlaylistId, &p.ThumbnailURL); err != nil {
			return nil, fmt.Errorf("error fetching playlists: %v", err)
		}
		lstPlaylists = append(lstPlaylists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching playlists: %v", err)
	}

	return lstPlaylists, nil
}

func (r *repository) GetVideoSearchInfo() ([]entities.ContentSearch, error) {

	var lstSearchInfo []entities.ContentSearch

	rows, err := r.db.Query(q.GetVideoSearchChannelTitle)
	if err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var cs entities.ContentSearch
		if err := rows.Scan(&cs.VideoId, &cs.Title, &cs.Channel); err != nil {
			return nil, fmt.Errorf("error fetching search content: %v", err)
		}
		lstSearchInfo = append(lstSearchInfo, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}

	return lstSearchInfo, nil
}

func (r *repository) GetContentById(contentId int) ([]entities.Videos, error) {

	var lstVideos []entities.Videos

	rows, err := r.db.Query(q.GetVideoMetadata, contentId)
	if err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var v entities.Videos
		if err := rows.Scan(&v.Id, &v.Title, &v.Description, &v.DurationSeconds, &v.OriginalURL, &v.WebpageURL, &v.IsFileDownloaded,
			&v.IsDeleted, &v.Channel.Name, &v.LiveStatus, &v.Domain.Domain, &v.LikesCount, &v.ViewsCount, &v.WatchCount, &v.UploadDate,
			&v.Availability, &v.Format.Format, &v.YoutubeVideoId, &v.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching videos: %v", err)
		}
		lstVideos = append(lstVideos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching Videos: %v", err)
	}

	//tags
	tagsRows, err := r.db.Query(q.GetVideoTags, contentId)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for Video: %v", err)
	}
	defer tagsRows.Close()

	var lstTags []entities.Tags
	for tagsRows.Next() {
		var tags entities.Tags
		if err := tagsRows.Scan(&tags.Id, &tags.Name, &tags.IsUsed, &tags.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching tags: %v", err)
		}
		lstTags = append(lstTags, tags)
	}

	//assign Tags to video
	for i := range lstVideos {
		for k, elem := range lstTags {
			if lstVideos[i].Id == lstTags[k].Id {
				lstVideos[i].Tags = append(lstVideos[i].Tags, elem)
			}
		}
	}

	//categories
	categoryRows, err := r.db.Query(q.GetVideoCategories, contentId)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for Video: %v", err)
	}
	defer categoryRows.Close()

	var lstCategories []entities.Categories
	for categoryRows.Next() {
		var categories entities.Categories
		if err := categoryRows.Scan(&categories.Id, &categories.Name, &categories.IsUsed, &categories.CreatedDate); err != nil {
			return nil, fmt.Errorf("error fetching categories: %v", err)
		}
		lstCategories = append(lstCategories, categories)
	}

	//assign categories to video
	for i := range lstVideos {
		for k, elem := range lstCategories {
			if lstVideos[i].Id == lstCategories[k].Id {
				lstVideos[i].Categories = append(lstVideos[i].Categories, elem)
			}
		}
	}

	//files
	filesRows, err := r.db.Query(q.GetVideoFiles, contentId)
	if err != nil {
		return nil, fmt.Errorf("error fetching files for Video: %v", err)
	}
	defer filesRows.Close()

	var lstFiles []entities.Files
	for filesRows.Next() {
		var files entities.Files
		if err := filesRows.Scan(&files.VideoId, &files.FileType, &files.FileSize, &files.Extension, &files.FilePath, &files.FileName); err != nil {
			return nil, fmt.Errorf("error fetching files: %v", err)
		}
		lstFiles = append(lstFiles, files)
	}

	//assign Files to video
	for i := range lstVideos {
		for k, elem := range lstFiles {
			if lstVideos[i].Id == lstFiles[k].VideoId {
				lstVideos[i].Files = append(lstVideos[i].Files, elem)
			}
		}
	}

	return lstVideos, nil
}
