package extractor

import "testing"

func TestExplainMetadataFailure(t *testing.T) {
	cases := []struct {
		name  string
		lines []string
		want  string
	}{
		{"bot check", []string{"WARNING: [youtube] No supported JavaScript runtime could be found.", "ERROR: [youtube] abc: Sign in to confirm you’re not a bot. Use --cookies-from-browser"},
			"YouTube is asking this server to sign in (bot check), so the video can't be read right now. Try again later, or give yt-dlp your cookies."},
		{"unavailable", []string{"ERROR: [youtube] abc: This video is unavailable"}, "This video is unavailable on YouTube."},
		{"private", []string{"ERROR: [youtube] abc: Private video. Sign in if you've been granted access"}, "This video is private."},
		{"other error", []string{"ERROR: [youtube] abc: Something new went wrong"}, "yt-dlp could not read this URL: [youtube] abc: Something new went wrong"},
		{"last error wins", []string{"ERROR: first problem", "ERROR: [youtube] abc: This video is unavailable"}, "This video is unavailable on YouTube."},
		{"warnings only", []string{"WARNING: something"}, "No video information was returned for this URL. If it was downloaded before, yt-dlp's download archive may be skipping it."},
		{"no output", nil, "No video information was returned for this URL. If it was downloaded before, yt-dlp's download archive may be skipping it."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := explainMetadataFailure(tc.lines); got != tc.want {
				t.Errorf("explainMetadataFailure() = %q, want %q", got, tc.want)
			}
		})
	}
}
