package global

import (
	"strconv"
	"sync"

	cfg "github.com/rs-anantmishra/streamsphere/config"
)

var onceDownloadStatus sync.Once
var onceActiveItem sync.Once
var onceCurrentQueueIndex sync.Once
var onceQueueAlive sync.Once

// -- todo --------------------------------------------------
// handle case where Video is already added - Skip download.
// If Video is downloaded and is part of playlist being downloaded then,
// video will be added twice.
// -- todo --------------------------------------------------

// Singletons
type DownloadStatus struct {
	VideoId       int    // Local VideoId
	VideoURL      string // Network Video Id
	StatusMessage string // completion percentage is in here
	State         int    // 0 = Unassigned, 1 = Queued, 2 = Downloading, 3 = Completed, 4 = Failed
	Title         string // VideoTitle
	Channel       string // Channel
	Thumbnail     string // Thumbnail base64
	Duration      int    // content duration
	Description   string // Video descriptions
	Likes         int    //Likes Count
	// PlaylistOrVideo string // Possible Values: "Playlist" or "Video"
	// PlaylistId      int    // Playlist Id from db
}

var (
	dsQueue           []DownloadStatus
	activeItem        []DownloadStatus
	currentQueueIndex []int
	queueAlive        []int
)

func NewActiveItem() []DownloadStatus {
	onceActiveItem.Do(func() { // <-- atomic, does not allow repeating
		activeItem = make([]DownloadStatus, 1) // <-- thread safe
	})

	return activeItem
}

func NewDownloadStatus() []DownloadStatus {

	maxQueueLength, err := strconv.Atoi((cfg.Config("MAX_QUEUE", false)))
	if err != nil {
		maxQueueLength = 2000
	}

	onceDownloadStatus.Do(func() { // <-- atomic, does not allow repeating
		dsQueue = make([]DownloadStatus, maxQueueLength) // <-- thread safe
	})

	return dsQueue
}

func NewQueueAlive() []int {
	onceQueueAlive.Do(func() { // <-- atomic, does not allow repeating
		queueAlive = make([]int, 1) // <-- thread safe
	})

	return queueAlive
}

func NewCurrentQueueIndex() []int {
	onceCurrentQueueIndex.Do(func() { // <-- atomic, does not allow repeating
		currentQueueIndex = make([]int, 1) // <-- thread safe
		currentQueueIndex[0] = 0
	})

	return currentQueueIndex
}

// call when downloading is not in progress.
// set DequeueRequired and use within the service.
func DefragmentQueue() {

	var items []DownloadStatus
	for k := range dsQueue {
		if dsQueue[k].State == Completed {
			dsQueue[k].VideoId = 0
			dsQueue[k].VideoURL = ""
			dsQueue[k].StatusMessage = ""
			dsQueue[k].State = Unassigned
		} else if dsQueue[k].State == Queued {
			items = append(items, dsQueue[k])
		}
	}

	currentQueueIndex[0] = len(items)
	copy(dsQueue, items)
}

const (
	Unassigned  = iota
	Queued      = iota
	Downloading = iota
	Completed   = iota
	Failed      = iota
)
