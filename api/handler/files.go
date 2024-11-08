package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	sql "github.com/rs-anantmishra/streamsphere/database"
	"github.com/rs-anantmishra/streamsphere/pkg/files"
)

func StorageStatus(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	status := `success`
	message := ``

	//Instantiate
	filesRepo := files.NewFilesRepo(sql.DB)
	filesSvc := files.NewFilesService(filesRepo)

	result, err := filesSvc.StorageStatusInfo()

	if err != nil {
		status = `failure`
		log.Info("error: ", err)

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": status, "message": message, "data": result})
}

func UploadMedia(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

// Local or Network Media via NFS or SAMBA
func IngestMedia(c *fiber.Ctx) error {
	//log context
	log.Info("Request Params:", c)

	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
