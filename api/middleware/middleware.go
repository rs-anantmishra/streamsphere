package middleware

import (
	"net/http"
	"strings"
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

	//the built UI
	serveStatic(app, http.Dir("../browser/"))
}

// serveStatic serves the built UI from root
func serveStatic(app *fiber.App, root http.FileSystem) {
	// index.html and main-*.js hold the addresses of the environment: the entrypoint rewrites main-*.js at every start
	// under the same file name, so a browser must ask again for them instead of keeping them for the hour
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if path := c.Path(); path == "/" || path == "/index.html" || strings.HasPrefix(path, "/main-") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return err
	})

	// Or extend your config for customization
	app.Use(filesystem.New(filesystem.Config{
		Root:   root,
		Browse: false,
		Index:  "index.html",
		MaxAge: 3600,
	}))
}
