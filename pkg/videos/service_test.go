package videos

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs-anantmishra/streamsphere/pkg/entities"
)

// fakeRepo answers only DeleteContentById; the other methods are never called by these tests.
type fakeRepo struct {
	deleted                bool
	contentPath, thumbPath string
	err                    error
}

func (f fakeRepo) DeleteContentById(int) (bool, string, string, error) {
	return f.deleted, f.contentPath, f.thumbPath, f.err
}
func (fakeRepo) GetAllVideos() ([]entities.Videos, error)              { return nil, nil }
func (fakeRepo) GetContentById(int) ([]entities.Videos, error)         { return nil, nil }
func (fakeRepo) GetPlaylistVideos(int) ([]entities.Videos, error)      { return nil, nil }
func (fakeRepo) GetAllPlaylists() ([]entities.Playlist, error)         { return nil, nil }
func (fakeRepo) GetVideoSearchInfo() ([]entities.ContentSearch, error) { return nil, nil }

func tempFile(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestDeleteContentByIdRemovesBothFiles(t *testing.T) {
	dir := t.TempDir()
	content, thumb := tempFile(t, dir, "video.mp4"), tempFile(t, dir, "video.webp")

	ok, err := NewVideoService(fakeRepo{deleted: true, contentPath: content, thumbPath: thumb}).DeleteContentById(1)

	if !ok || err != nil {
		t.Fatalf("DeleteContentById() = %v, %v, want true, nil", ok, err)
	}
	if exists(content) || exists(thumb) {
		t.Error("expected both files to be removed")
	}
}

func TestDeleteContentByIdIgnoresFilesAlreadyGone(t *testing.T) {
	dir := t.TempDir()
	thumb := tempFile(t, dir, "video.webp")

	ok, err := NewVideoService(fakeRepo{deleted: true, contentPath: filepath.Join(dir, "missing.mp4"), thumbPath: thumb}).DeleteContentById(1)

	if !ok || err != nil {
		t.Fatalf("DeleteContentById() = %v, %v, want true, nil", ok, err)
	}
	if exists(thumb) {
		t.Error("expected the thumbnail to be removed")
	}
}

func TestDeleteContentByIdReportsAFileThatCannotBeRemoved(t *testing.T) {
	dir := t.TempDir()
	notEmpty := filepath.Join(dir, "folder")
	if err := os.Mkdir(notEmpty, 0o755); err != nil {
		t.Fatal(err)
	}
	tempFile(t, notEmpty, "inside.txt")

	ok, err := NewVideoService(fakeRepo{deleted: true, contentPath: notEmpty, thumbPath: filepath.Join(dir, "x.webp")}).DeleteContentById(1)

	if ok || err == nil || !strings.Contains(err.Error(), notEmpty) {
		t.Fatalf("DeleteContentById() = %v, %v, want false and an error naming %s", ok, err, notEmpty)
	}
}

func TestDeleteContentByIdTouchesNoFilesWhenTheDatabaseDeleteFails(t *testing.T) {
	dir := t.TempDir()
	content := tempFile(t, dir, "video.mp4")

	ok, err := NewVideoService(fakeRepo{deleted: false, contentPath: content, err: errors.New("db failure")}).DeleteContentById(1)

	if ok || err == nil {
		t.Fatalf("DeleteContentById() = %v, %v, want false and an error", ok, err)
	}
	if !exists(content) {
		t.Error("the file must stay when the database delete failed")
	}
}
