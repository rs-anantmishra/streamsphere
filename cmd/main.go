package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	m "github.com/rs-anantmishra/streamsphere/api/middleware"
	r "github.com/rs-anantmishra/streamsphere/api/routes"
	c "github.com/rs-anantmishra/streamsphere/config"
	"github.com/rs-anantmishra/streamsphere/database"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Streamsphere",
		UnescapePath:  true,
	})

	port := ":" + c.Config("APP_PORT", false)

	m.SetupMiddleware(app)
	r.SetupRoutes(app)

	database.ConnectDB()
	defer database.CloseDB()

	log.Fatal(app.Listen(port))
}
