package videos

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// videoOnDisk is a downloaded video as the app stores it: rows in the database and the two files it points to.
type videoOnDisk struct {
	id           int
	contentPath  string
	thumbnailPth string
}

// addVideoOnDisk adds a video with real files under dir, tagged and categorised, optionally linked to a playlist.
func addVideoOnDisk(t *testing.T, db *sql.DB, dir string, videoID, playlistID, index int) videoOnDisk {
	t.Helper()
	v := videoOnDisk{
		id:           videoID,
		contentPath:  filepath.Join(dir, fmt.Sprintf("video%d.mp4", videoID)),
		thumbnailPth: filepath.Join(dir, fmt.Sprintf("thumb%d.jpg", videoID)),
	}
	for _, p := range []string{v.contentPath, v.thumbnailPth} {
		if err := os.WriteFile(p, []byte("file"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	stmts := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO tblVideos (Id, OriginalURL, WebpageURL, YoutubeVideoId, CreatedDate) VALUES (?, 'o', 'w', ?, 1)`, []any{videoID, videoID}},
		{`INSERT INTO tblFiles (VideoId, FileType, FilePath, FileName) VALUES (?, 'Thumbnail', ?, ?)`, []any{videoID, dir, filepath.Base(v.thumbnailPth)}},
		{`INSERT INTO tblFiles (VideoId, FileType, FilePath, FileName) VALUES (?, 'Video', ?, ?)`, []any{videoID, dir, filepath.Base(v.contentPath)}},
		// tag 1 and category 1 are shared by every video; the second pair belongs to this video alone
		{`INSERT OR IGNORE INTO tblTags (Id, Name, IsUsed, CreatedDate) VALUES (1, 'shared', 1, 1)`, nil},
		{`INSERT OR IGNORE INTO tblCategories (Id, Name, IsUsed, CreatedDate) VALUES (1, 'shared', 1, 1)`, nil},
		{`INSERT INTO tblTags (Id, Name, IsUsed, CreatedDate) VALUES (?, 'own', 1, 1)`, []any{100 + videoID}},
		{`INSERT INTO tblVideoFileTags (TagId, VideoId, CreatedDate) VALUES (1, ?, 1)`, []any{videoID}},
		{`INSERT INTO tblVideoFileTags (TagId, VideoId, CreatedDate) VALUES (?, ?, 1)`, []any{100 + videoID, videoID}},
		{`INSERT INTO tblVideoFileCategories (CategoryId, VideoId, CreatedDate) VALUES (1, ?, 1)`, []any{videoID}},
		{`INSERT INTO tblPlaylistVideoFiles (VideoId, PlaylistId, PlaylistVideoIndex, CreatedDate) VALUES (?, ?, ?, 1)`, []any{videoID, playlistID, index}},
	}
	for _, s := range stmts {
		if _, err := db.Exec(s.query, s.args...); err != nil {
			t.Fatal(err)
		}
	}
	return v
}

func rows(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// rowsOf counts everything the app keeps for one video.
func rowsOf(t *testing.T, db *sql.DB, videoID int) map[string]int {
	t.Helper()
	return map[string]int{
		"tblVideos":              rows(t, db, `SELECT COUNT(*) FROM tblVideos WHERE Id = ?`, videoID),
		"tblFiles":               rows(t, db, `SELECT COUNT(*) FROM tblFiles WHERE VideoId = ?`, videoID),
		"tblVideoFileTags":       rows(t, db, `SELECT COUNT(*) FROM tblVideoFileTags WHERE VideoId = ?`, videoID),
		"tblVideoFileCategories": rows(t, db, `SELECT COUNT(*) FROM tblVideoFileCategories WHERE VideoId = ?`, videoID),
		"tblPlaylistVideoFiles":  rows(t, db, `SELECT COUNT(*) FROM tblPlaylistVideoFiles WHERE VideoId = ?`, videoID),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestDeleteRemovesTheVideoItsRowsAndItsFiles(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist")
	v := addVideoOnDisk(t, db, dir, 1, 1, 1)
	want := map[string]int{"tblVideos": 1, "tblFiles": 2, "tblVideoFileTags": 2, "tblVideoFileCategories": 1, "tblPlaylistVideoFiles": 1}
	if got := rowsOf(t, db, 1); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("before the delete: %v, want %v", got, want)
	}

	ok, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(1)

	if !ok || err != nil {
		t.Fatalf("DeleteContentById() = %v, %v, want true, nil", ok, err)
	}
	for table, n := range rowsOf(t, db, 1) {
		if n != 0 {
			t.Errorf("%s still has %d row(s) for the video", table, n)
		}
	}
	if fileExists(v.contentPath) || fileExists(v.thumbnailPth) {
		t.Error("the video or its thumbnail file is still on disk")
	}
}

func TestDeleteLeavesOtherVideosTheirFilesAndSharedTagsAlone(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist")
	addVideoOnDisk(t, db, dir, 1, 1, 1)
	other := addVideoOnDisk(t, db, dir, 2, 1, 2)
	before := rowsOf(t, db, 2)

	if _, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(1); err != nil {
		t.Fatal(err)
	}

	if after := rowsOf(t, db, 2); fmt.Sprint(after) != fmt.Sprint(before) {
		t.Errorf("the other video's rows changed: %v, was %v", after, before)
	}
	if !fileExists(other.contentPath) || !fileExists(other.thumbnailPth) {
		t.Error("the other video's files were removed")
	}
	if rows(t, db, `SELECT COUNT(*) FROM tblTags WHERE Id = 1`) != 1 || rows(t, db, `SELECT COUNT(*) FROM tblCategories WHERE Id = 1`) != 1 {
		t.Error("the shared tag or category was removed")
	}
}

func TestDeleteInOnePlaylistLeavesOtherPlaylistsAlone(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "First")
	addVideoOnDisk(t, db, dir, 1, 1, 1)
	addVideoOnDisk(t, db, dir, 2, 1, 2)
	addPlaylist(t, db, 2, "Second")
	addVideoOnDisk(t, db, dir, 3, 2, 1)
	addVideoOnDisk(t, db, dir, 4, 2, 2)
	repo := NewVideoRepo(db)

	if _, err := NewVideoService(repo).DeleteContentById(1); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("playlists = %+v, want both still listed", got)
	}
	if n := rows(t, db, `SELECT COUNT(*) FROM tblPlaylistVideoFiles WHERE PlaylistId = 2`); n != 2 {
		t.Errorf("the second playlist has %d linked videos, want 2", n)
	}
}

func TestDeleteOfAnUnknownIdChangesNothing(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist")
	v := addVideoOnDisk(t, db, dir, 1, 1, 1)
	before := rowsOf(t, db, 1)

	ok, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(99)

	if ok || err == nil {
		t.Errorf("DeleteContentById(99) = %v, %v, want false and an error", ok, err)
	}
	if after := rowsOf(t, db, 1); fmt.Sprint(after) != fmt.Sprint(before) {
		t.Errorf("rows changed: %v, was %v", after, before)
	}
	if !fileExists(v.contentPath) || !fileExists(v.thumbnailPth) {
		t.Error("files were removed")
	}
}

func TestDeleteSucceedsWhenTheFilesAreAlreadyGone(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist")
	v := addVideoOnDisk(t, db, dir, 1, 1, 1)
	_ = os.Remove(v.contentPath)
	_ = os.Remove(v.thumbnailPth)

	ok, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(1)

	if !ok || err != nil {
		t.Fatalf("DeleteContentById() = %v, %v, want true, nil", ok, err)
	}
	if n := rows(t, db, `SELECT COUNT(*) FROM tblVideos WHERE Id = 1`); n != 0 {
		t.Error("the video row is still there")
	}
}

func TestDeleteWorksForAVideoThatIsInNoPlaylist(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	// the app files a standalone video under the "No Playlist" row, id -1
	if _, err := db.Exec(`INSERT INTO tblPlaylists (Id, Title, ItemCount, PlaylistUploader, YoutubePlaylistId, CreatedDate) VALUES (-1, 'No Playlist', 0, 'u', 'Not Available', 1)`); err != nil {
		t.Fatal(err)
	}
	addVideoOnDisk(t, db, dir, 1, -1, 0)

	ok, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(1)

	if !ok || err != nil {
		t.Fatalf("DeleteContentById() = %v, %v, want true, nil", ok, err)
	}
	for table, n := range rowsOf(t, db, 1) {
		if n != 0 {
			t.Errorf("%s still has %d row(s) for the video", table, n)
		}
	}
}

// The next two pin what a delete does NOT do to the playlist today, so a later fix shows up as a deliberate change.

func TestDeleteDoesNotRenumberTheOtherVideosOrUpdateTheItemCount(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist") // ItemCount 3
	addVideoOnDisk(t, db, dir, 1, 1, 1)
	addVideoOnDisk(t, db, dir, 2, 1, 2)
	addVideoOnDisk(t, db, dir, 3, 1, 3)

	if _, err := NewVideoService(NewVideoRepo(db)).DeleteContentById(2); err != nil {
		t.Fatal(err)
	}

	indices := fmt.Sprint(rows(t, db, `SELECT COUNT(*) FROM tblPlaylistVideoFiles WHERE PlaylistId = 1 AND PlaylistVideoIndex = 1`),
		rows(t, db, `SELECT COUNT(*) FROM tblPlaylistVideoFiles WHERE PlaylistId = 1 AND PlaylistVideoIndex = 2`),
		rows(t, db, `SELECT COUNT(*) FROM tblPlaylistVideoFiles WHERE PlaylistId = 1 AND PlaylistVideoIndex = 3`))
	if indices != "1 0 1" {
		t.Errorf("videos at indices 1, 2, 3: %s, want 1 0 1 (a gap at 2)", indices)
	}
	if n := rows(t, db, `SELECT ItemCount FROM tblPlaylists WHERE Id = 1`); n != 3 {
		t.Errorf("ItemCount = %d, want it left at 3", n)
	}
}

func TestDeletingEveryVideoLeavesTheEmptyPlaylistRowButHidesIt(t *testing.T) {
	db, dir := newPlaylistsDB(t), t.TempDir()
	addPlaylist(t, db, 1, "A playlist")
	addVideoOnDisk(t, db, dir, 1, 1, 1)
	addVideoOnDisk(t, db, dir, 2, 1, 2)
	repo := NewVideoRepo(db)
	svc := NewVideoService(repo)

	for _, id := range []int{1, 2} {
		if _, err := svc.DeleteContentById(id); err != nil {
			t.Fatal(err)
		}
	}

	if n := rows(t, db, `SELECT COUNT(*) FROM tblPlaylists WHERE Id = 1`); n != 1 {
		t.Errorf("playlist rows = %d, want the emptied playlist's row kept", n)
	}
	got, err := repo.GetAllPlaylists()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("playlists = %+v, want the emptied playlist hidden", got)
	}
}
