package handler

import (
	"encoding/json"
	"errors"
	"strconv"
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
	if len(result) == 0 {
		err = errors.New("invalid url provided")
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

func DownloadStatus(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	//global MPI
	const _blank string = ""
	activeItem := g.NewActiveItem()
	mt = websocket.TextMessage
	terminate := false

	for {
		if len(activeItem) > 0 && activeItem[0].VideoURL != _blank {

			dsr := res.DownloadStatusResponse{Message: activeItem[0].StatusMessage, VideoURL: activeItem[0].VideoURL}
			jsonData, e := json.Marshal(dsr)
			if e != nil {
				log.Info(e)
			}
			msg = []byte(jsonData)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Info("write:", err)
				break
			}
		} else {
			// Send a WebSocket close message
			deadline := time.Now().Add(time.Second * 5)
			err := c.Conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				deadline,
			)
			if err != nil {
				log.Info("ws error:", err)
			}

			// Set deadline for reading the next message
			err = c.Conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				log.Info("ws error:", err)
			}
			// Read messages until the close message is confirmed
			for {
				_, _, err = c.Conn.NextReader()
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					terminate = true
					break
				}
				if err != nil {
					terminate = true
					break
				}
			}
			// Close the TCP connection
			err = c.Conn.Close()
			if err != nil {
				log.Info("ws error:", err)
			}

			if terminate {
				break
			}
		}

		//transmit data once per second
		duration := time.Second
		time.Sleep(duration)
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
