package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	sql "github.com/rs-anantmishra/streamsphere/database"
	"github.com/rs-anantmishra/streamsphere/pkg/videos"
)

// Get All Videos
func GetAllVideos(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	//Instantiate
	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.GetVideos()
	if err != nil {
		log.Info("error fetching all videos", err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// Get All Playlists
func GetAllPlaylists(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	status := `success`
	message := ``

	//Instantiate
	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.GetPlaylists()
	if err != nil {
		log.Info("error fetching all playlists", err)
		status = `failure`
	}

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}

// Get All Playlists
func GetPlaylistsVideo(c *fiber.Ctx) error {
	//log context
	id, _ := c.ParamsInt("id") // int 123 and no error
	log.Info("Request Params:", id)

	status := `success`
	message := ``

	//Instantiate
	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.GetPlaylistVideos(id)
	if err != nil {
		log.Info("error fetching all playlists", err)
		status = `failure`
	}

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}

// Get Video by Id
func GetContentById(c *fiber.Ctx) error {
	//log context
	id, _ := c.ParamsInt("id") // int 123 and no error
	log.Info("Request Params:", id)

	status := `success`
	message := ``

	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.GetContentById(id)
	if err != nil {
		log.Info("error fetching video", err)
		status = `failure`
	}

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}

// Get Video by Id
func DeleteContentById(c *fiber.Ctx) error {
	//log context
	id, _ := c.ParamsInt("id") // int 123 and no error
	log.Info("Request Params:", id)

	status := `success`
	message := `content deleted`

	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.DeleteContentById(id)
	if err != nil {
		log.Info("error deleting content", err)
		status = `failure`
		message = ``
	}

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}

// Get All Audios
func GetAllAudios(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Get Media By Tags
func GetMediaByTags(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Get Media By Categories
func GetVideosByCategories(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Get Media By Domain
func GetVideosByDomain(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Get Media By Channel
func GetVideosByChannel(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Search Media Files
func GetMediaBySearch(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Search all media by youtube Id
func GetMediaByYoutubeId(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

func GetMediaByPhysicalLocation(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

func GetVideoSearchData(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	//Instantiate
	svcRepo := videos.NewVideoRepo(sql.DB)
	svcVideos := videos.NewVideoService(svcRepo)

	result, err := svcVideos.GetVideoSearchData()
	message := strconv.Itoa(len(result)) + ` records found`
	status := `success`
	if err != nil {
		log.Info("error fetching all videos", err)
		status = `failure`
	}
	return c.Status(fiber.StatusOK).Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}
