package extractor

import (
	"testing"
)

// the content host is added by the UI, so the backend only returns the path below the media directory
func TestUrlTransformsReturnsThePathWithoutAHost(t *testing.T) {
	t.Setenv("MEDIA_PATH", "..#content")
	t.Setenv("FILE_HOSTING", "http://old.host:3500")

	cases := []struct{ name, in, want string }{
		{"file under the media directory", "../content/youtube.com/jawed/Videos/x.webm", "/youtube.com/jawed/Videos/x.webm"},
		{"windows separators after the prefix", "../content/youtube.com\\jawed\\Videos\\x.webm", "/youtube.com/jawed/Videos/x.webm"},
		{"nothing to link to stays empty", "", ""},
		{"a path outside the media directory is left alone", "/elsewhere/x.webm", "/elsewhere/x.webm"},
	}
	for _, c := range cases {
		if got := urlTransforms(c.in); got != c.want {
			t.Errorf("%s: urlTransforms(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
	}
}
