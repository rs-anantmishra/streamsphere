package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupMiddleware(app *fiber.App) {

	//pprof
	//app.Use(pprof.New())

	//logger
	app.Use(logger.New())

	//setup CORS
	app.Use(cors.New())

	//Rate-Limiter
	app.Use(limiter.New(limiter.Config{
		Max:               500,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	//Web-Sockets Request
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Or extend your config for customization
	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("../browser/"),
		Browse: false,
		Index:  "index.html",
		MaxAge: 3600,
	}))
}
