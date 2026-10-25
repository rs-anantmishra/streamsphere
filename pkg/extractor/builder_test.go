package extractor

import (
	"strings"
	"testing"
)

func TestUrlArgsTurnsTheDownloadArchiveOff(t *testing.T) {
	got := urlArgs("https://www.youtube.com/watch?v=abc")
	want := []string{`"https://www.youtube.com/watch?v=abc"`, "--no-download-archive"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("urlArgs = %q, want %q", got, want)
	}
}

func TestMetadataCommandTurnsTheDownloadArchiveOff(t *testing.T) {
	args, _, _ := cmdBuilderMetadata("https://www.youtube.com/playlist?list=PLabc")
	if !strings.HasPrefix(args, `"https://www.youtube.com/playlist?list=PLabc" --no-download-archive `) {
		t.Errorf("metadata arguments start with %q", args[:min(len(args), 90)])
	}
}
