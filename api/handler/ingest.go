package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/rs-anantmishra/streamsphere/api/presenter"
	res "github.com/rs-anantmishra/streamsphere/api/presenter"
	cfg "github.com/rs-anantmishra/streamsphere/config"
	sql "github.com/rs-anantmishra/streamsphere/database"
	en "github.com/rs-anantmishra/streamsphere/pkg/entities"
	ex "github.com/rs-anantmishra/streamsphere/pkg/extractor"
	g "github.com/rs-anantmishra/streamsphere/pkg/global"
)

func NetworkIngestMetadata(c *fiber.Ctx) error {

	status := `failure`
	message := ``

	//bind incoming data
	params := new(en.IncomingRequest)
	if err := c.BodyParser(params); err != nil {
		return nil
	}

	//log incoming data
	log.Info("Request Params:", params)

	//Instantiate
	svcDownloads := ex.NewDownload(*params)
	svcRepo := ex.NewDownloadRepo(sql.DB)
	svcVideos := ex.NewDownloadService(svcRepo, svcDownloads)

	//Process Request
	// No validations for URL/Playlist are needed.
	// If Metadata is not fetched, and there is an error message from yt-dlp
	// just show that error on the UI
	result, err := svcVideos.ExtractIngestMetadata(*params)
	if len(result) == 0 && err == nil {
		err = errors.New("no video information was returned for this url")
	}
	if err != nil {
		status = `failure`
		message = err.Error()
	}

	if result != nil {
		//queue downloads below
		//global MPI
		maxQueueLength, _ := strconv.Atoi((cfg.Config("MAX_QUEUE", false)))
		lstDownloads := g.NewDownloadStatus()
		qAlive := g.NewQueueAlive()
		currentQueueIndex := g.NewCurrentQueueIndex()

		if maxQueueLength-currentQueueIndex[0]-len(result) >= 0 {
			for idx := range result {
				lstDownloads[currentQueueIndex[0]] = g.DownloadStatus{VideoId: result[idx].VideoId,
					VideoURL:      result[idx].WebpageURL,
					StatusMessage: "",
					State:         g.Queued,
					Title:         result[idx].Title,
					Channel:       result[idx].Channel,
					Duration:      result[idx].Duration,
					Thumbnail:     result[idx].Thumbnail,
					Likes:         result[idx].LikesCount,
				}
				currentQueueIndex[0]++
			}
		} else {
			//send error response that queue is full. Please wait for existing downloads to complete.
		}

		if qAlive[0] != 1 {
			qAlive[0] = 1
			go svcVideos.ExtractIngestMedia()
		}
		status = `success`
		message = strconv.Itoa(len(result)) + ` records added too queue.`
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
	// return c.Status(fiber.StatusOK).JSON(result)
}

func NetworkIngestMedia(c *fiber.Ctx) error {

	maxQueueLength, _ := strconv.Atoi((cfg.Config("MAX_QUEUE", false)))
	//bind incoming data
	params := new(en.QueueDownloads)

	if err := c.BodyParser(params); err != nil {
		return err
	}

	//log incoming data
	log.Info("Request Params:", params)

	//global MPI
	lstDownloads := g.NewDownloadStatus()
	qAlive := g.NewQueueAlive()
	currentQueueIndex := g.NewCurrentQueueIndex()

	if maxQueueLength-currentQueueIndex[0]-len(params.DownloadVideos) >= 0 {
		for idx := range params.DownloadVideos {
			lstDownloads[currentQueueIndex[0]] = g.DownloadStatus{VideoId: params.DownloadVideos[idx].VideoId, VideoURL: params.DownloadVideos[idx].VideoURL, StatusMessage: "", State: g.Queued}
			currentQueueIndex[0]++
		}
	} else {
		//send error response that queue is full. Please wait for existing downloads to complete.
	}

	//Instantiate
	svcDownloads := ex.NewDownload(en.IncomingRequest{})
	svcRepo := ex.NewDownloadRepo(sql.DB)
	svcVideos := ex.NewDownloadService(svcRepo, svcDownloads)

	if qAlive[0] != 1 {
		qAlive[0] = 1
		go svcVideos.ExtractIngestMedia()
	}

	result := res.QueueResponse{Result: "Item added to download queue successfully."}
	return c.Status(fiber.StatusOK).JSON(result)
}

// downloadStatusResponse builds the status message; an idle item keeps only what identifies the last download (id, title, channel, thumbnail),
// so the page can show it again after a reload or a visit to another page.
func downloadStatusResponse(item g.DownloadStatus) res.DownloadStatusResponse {
	state := downloadState(item)
	if state == "idle" {
		return res.DownloadStatusResponse{State: state, VideoId: item.VideoId, Title: item.Title, Channel: item.Channel, Thumbnail: item.Thumbnail}
	}
	return res.DownloadStatusResponse{
		Message:   item.StatusMessage,
		VideoURL:  item.VideoURL,
		State:     state,
		VideoId:   item.VideoId,
		Title:     item.Title,
		Channel:   item.Channel,
		Thumbnail: item.Thumbnail,
	}
}

// downloadState maps the active item's status message to a wire-format state.
func downloadState(item g.DownloadStatus) string {
	switch {
	case item.VideoURL == "":
		return "idle"
	case strings.HasPrefix(item.StatusMessage, "Error:"):
		return "failed"
	case item.StatusMessage == "Download completed successfully.":
		return "completed"
	case item.StatusMessage == "[download] Saving file to library...":
		return "saving"
	default:
		return "downloading"
	}
}

// DownloadStatus keeps the socket open and streams the active item's state
// once per second; idle heartbeats are sent at most every 5 seconds.
func DownloadStatus(c *websocket.Conn) {
	const idleInterval = 5 * time.Second
	const writeTimeout = 5 * time.Second
	activeItem := g.NewActiveItem()
	var lastIdle time.Time

	// Read so the library processes the client's close frame; any error ends the handler.
	// Capture the raw conn: c is pooled and its embedded Conn is nilled when the handler returns.
	ws := c.Conn
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()
	// Stop the reader before returning: fasthttp recycles the connection and its buffers once the handler ends.
	defer func() {
		_ = ws.SetReadDeadline(time.Now())
		<-done
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		dsr := res.DownloadStatusResponse{State: "idle"}
		if len(activeItem) > 0 {
			dsr = downloadStatusResponse(activeItem[0])
		}

		if dsr.State != "idle" || time.Since(lastIdle) >= idleInterval {
			jsonData, e := json.Marshal(dsr)
			if e != nil {
				log.Info(e)
			}
			// A client that vanished without closing must fail the write instead of lingering in the OS buffer.
			_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Info("write:", err)
				break
			}
			if dsr.State == "idle" {
				lastIdle = time.Now()
			}
		}

		//transmit data once per second
		select {
		case <-done:
			return
		case <-ticker.C:
		}
	}
}

func NetworkIngestAutoSubs(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": "nil"})
}

func NetworkIngestThumbnail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": "nil"})
}

func NetworkIngestQueuedItems(c *fiber.Ctx) error {

	qry := c.Queries()
	state := qry["state"]

	allQueueItems := g.NewDownloadStatus()
	var queuedItems []presenter.CardsInfoResponse
	for _, elem := range allQueueItems {
		if elem.State == g.Queued && state == "queued" {
			queuedItems = append(queuedItems, presenter.CardsInfoResponse{
				VideoId:     elem.VideoId,
				Title:       elem.Title,
				Description: elem.Description,
				Duration:    elem.Duration,
				WebpageURL:  elem.VideoURL,
				Thumbnail:   elem.Thumbnail,
				Channel:     elem.Channel,
			})
		} else if elem.State == g.Downloading && state == "downloading" {
			queuedItems = append(queuedItems, presenter.CardsInfoResponse{
				VideoId:     elem.VideoId,
				Title:       elem.Title,
				Description: elem.Description,
				Duration:    elem.Duration,
				WebpageURL:  elem.VideoURL,
				Thumbnail:   elem.Thumbnail,
				Channel:     elem.Channel,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(queuedItems)
}
