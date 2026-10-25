package videos

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// newPlaylistsDB builds an empty database from the app's own schema script.
func newPlaylistsDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// every connection to ":memory:" is a separate database
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	schema, err := os.ReadFile("../../database/scripts/create.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	return db
}

func addPlaylist(t *testing.T, db *sql.DB, id int, title string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO tblPlaylists (Id, Title, ItemCount, PlaylistUploader, YoutubePlaylistId, CreatedDate) VALUES (?, ?, 3, 'uploader', 'PLtest', 1)`, id, title)
	if err != nil {
		t.Fatal(err)
	}
}

// addPlaylistVideo adds a downloaded video with its thumbnail and video file rows, linked to a playlist at an index.
func addPlaylistVideo(t *testing.T, db *sql.DB, playlistID, videoID, index int) {
	t.Helper()
	stmts := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO tblVideos (Id, OriginalURL, WebpageURL, YoutubeVideoId, CreatedDate) VALUES (?, 'o', 'w', ?, 1)`, []any{videoID, videoID}},
		{`INSERT INTO tblFiles (VideoId, FileType, FilePath, FileName) VALUES (?, 'Thumbnail', '/thumbs', ?)`, []any{videoID, "thumb" + string(rune('0'+videoID)) + ".jpg"}},
		{`INSERT INTO tblFiles (VideoId, FileType, FilePath, FileName) VALUES (?, 'Video', '/videos', ?)`, []any{videoID, "video" + string(rune('0'+videoID)) + ".mp4"}},
		{`INSERT INTO tblPlaylistVideoFiles (VideoId, PlaylistId, PlaylistVideoIndex, CreatedDate) VALUES (?, ?, ?, 1)`, []any{videoID, playlistID, index}},
	}
	for _, s := range stmts {
		if _, err := db.Exec(s.query, s.args...); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPlaylistsShowTheFirstVideosThumbnail(t *testing.T) {
	db := newPlaylistsDB(t)
	addPlaylist(t, db, 1, "A playlist")
	addPlaylistVideo(t, db, 1, 1, 1)
	addPlaylistVideo(t, db, 1, 2, 2)
	addPlaylistVideo(t, db, 1, 3, 3)

	got, err := NewVideoRepo(db).GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || !strings.HasSuffix(got[0].ThumbnailURL, "thumb1.jpg") {
		t.Errorf("playlists = %+v, want one playlist with the thumbnail of video 1", got)
	}
}

func TestPlaylistStaysListedWhenItsFirstVideoIsDeleted(t *testing.T) {
	db := newPlaylistsDB(t)
	addPlaylist(t, db, 1, "A playlist")
	addPlaylistVideo(t, db, 1, 1, 1)
	addPlaylistVideo(t, db, 1, 2, 2)
	addPlaylistVideo(t, db, 1, 3, 3)
	repo := NewVideoRepo(db)

	if _, _, _, err := repo.DeleteContentById(1); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}

	// nothing renumbers the indices on a delete, so the lowest left (2) stands in for the first
	if len(got) != 1 || !strings.HasSuffix(got[0].ThumbnailURL, "thumb2.jpg") {
		t.Errorf("playlists = %+v, want the playlist still listed with the thumbnail of video 2", got)
	}
}

func TestPlaylistStaysListedWhenOnlyALaterVideoIsLeft(t *testing.T) {
	db := newPlaylistsDB(t)
	addPlaylist(t, db, 1, "A playlist")
	addPlaylistVideo(t, db, 1, 1, 1)
	addPlaylistVideo(t, db, 1, 2, 2)
	addPlaylistVideo(t, db, 1, 3, 3)
	repo := NewVideoRepo(db)

	for _, id := range []int{1, 2} {
		if _, _, _, err := repo.DeleteContentById(id); err != nil {
			t.Fatal(err)
		}
	}
	got, err := repo.GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || !strings.HasSuffix(got[0].ThumbnailURL, "thumb3.jpg") {
		t.Errorf("playlists = %+v, want the playlist listed with the thumbnail of video 3", got)
	}
}

func TestPlaylistWithNoVideosIsNotListed(t *testing.T) {
	db := newPlaylistsDB(t)
	addPlaylist(t, db, 1, "Emptied")
	addPlaylistVideo(t, db, 1, 1, 1)
	addPlaylist(t, db, 2, "Has a video")
	addPlaylistVideo(t, db, 2, 2, 1)
	repo := NewVideoRepo(db)

	if _, _, _, err := repo.DeleteContentById(1); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].Id != 2 {
		t.Errorf("playlists = %+v, want only playlist 2", got)
	}
}

func TestEachPlaylistIsListedOnceInIdOrder(t *testing.T) {
	db := newPlaylistsDB(t)
	addPlaylist(t, db, 2, "Second")
	addPlaylistVideo(t, db, 2, 4, 1)
	addPlaylistVideo(t, db, 2, 5, 2)
	addPlaylist(t, db, 1, "First")
	addPlaylistVideo(t, db, 1, 1, 2)
	addPlaylistVideo(t, db, 1, 2, 3)

	got, err := NewVideoRepo(db).GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got[0].Id != 1 || got[1].Id != 2 {
		t.Fatalf("playlists = %+v, want playlists 1 and 2, once each, in id order", got)
	}
	if !strings.HasSuffix(got[0].ThumbnailURL, "thumb1.jpg") || !strings.HasSuffix(got[1].ThumbnailURL, "thumb4.jpg") {
		t.Errorf("thumbnails = %q, %q, want the lowest-index video of each playlist", got[0].ThumbnailURL, got[1].ThumbnailURL)
	}
}
