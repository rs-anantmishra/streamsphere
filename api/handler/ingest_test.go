package handler

import (
	"testing"

	g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

func TestDownloadState(t *testing.T) {
	cases := []struct {
		name string
		item g.DownloadStatus
		want string
	}{
		{"no active item", g.DownloadStatus{VideoURL: ""}, "idle"},
		{"error", g.DownloadStatus{VideoURL: "u", StatusMessage: "Error:exit status 1"}, "failed"},
		{"completed", g.DownloadStatus{VideoURL: "u", StatusMessage: "Download completed successfully."}, "completed"},
		{"saving", g.DownloadStatus{VideoURL: "u", StatusMessage: "[download] Saving file to library..."}, "saving"},
		{"progress", g.DownloadStatus{VideoURL: "u", StatusMessage: "[download]  45.3% of 1.00MiB at 1.0MiB/s ETA 00:01"}, "downloading"},
		{"empty message", g.DownloadStatus{VideoURL: "u", StatusMessage: ""}, "downloading"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := downloadState(tc.item); got != tc.want {
				t.Errorf("downloadState() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDownloadStatusResponseIdleCarriesNoDetails(t *testing.T) {
	last := g.DownloadStatus{VideoId: 42, Title: "T", Channel: "C", Thumbnail: "th.jpg", StatusMessage: ""}
	got := downloadStatusResponse(last)
	if got.State != "idle" || got.VideoId != 0 || got.Title != "" || got.Channel != "" || got.Thumbnail != "" {
		t.Errorf("idle response = %+v, want state idle and no details", got)
	}

	active := g.DownloadStatus{VideoURL: "u", VideoId: 42, Title: "T", Channel: "C", Thumbnail: "th.jpg", StatusMessage: "[download]  1.0% of 1MiB"}
	got = downloadStatusResponse(active)
	if got.State != "downloading" || got.VideoId != 42 || got.Title != "T" || got.Thumbnail != "th.jpg" {
		t.Errorf("active response = %+v, want state downloading with details", got)
	}
}
