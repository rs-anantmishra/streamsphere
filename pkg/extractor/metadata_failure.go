package extractor

import "strings"

// explainMetadataFailure turns yt-dlp's output (stdout and stderr share one stream) into a plain reason
// for the metadata call returning nothing, so the page can tell the user what to do instead of "invalid url".
func explainMetadataFailure(lines []string) string {
	errorLine := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "ERROR:") {
			errorLine = line
		}
	}

	switch {
	case strings.Contains(errorLine, "Sign in to confirm"):
		return "YouTube is asking this server to sign in (bot check), so the video can't be read right now. Try again later, or give yt-dlp your cookies."
	case strings.Contains(errorLine, "Private video"):
		return "This video is private."
	case strings.Contains(errorLine, "unavailable"), strings.Contains(errorLine, "not available"):
		return "This video is unavailable on YouTube."
	case errorLine != "":
		return "yt-dlp could not read this URL: " + strings.TrimSpace(strings.TrimPrefix(errorLine, "ERROR:"))
	}
	// yt-dlp prints nothing for a video its download archive already holds, so that case lands here too
	return "No video information was returned for this URL. If it was downloaded before, yt-dlp's download archive may be skipping it."
}
