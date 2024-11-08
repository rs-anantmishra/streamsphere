package router

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs-anantmishra/metubeplus/api/handler"
)

func SetupRoutes(app *fiber.App) {

	//Middlewares
	api := app.Group("/api", logger.New())
	api.Get("/hello", handler.Hello)

	//Network Downloads: Playlist, Videos or Audios
	download := api.Group("/download", logger.New())

	download.Get("/queued-items", handler.NetworkIngestQueuedItems) //Get Queued-Items
	download.Post("/metadata", handler.NetworkIngestMetadata)       //Download Metadata + Thumbnail and save to db for [Playlists, Videos]
	download.Post("/media", handler.NetworkIngestMedia)             //Download Media File(s) and update db for [Playlists, Videos]

	//Specific to network Video and Audio only
	download.Post("/autosubs", handler.NetworkIngestAutoSubs)   //Download auto-subs for a video that exists in library [Videos]
	download.Post("/thumbnail", handler.NetworkIngestThumbnail) //Download thumbnail for a video that exists in library [Videos]

	//Web-Sockets
	app.Get("/ws/status", websocket.New(handler.DownloadStatus))

	//Todo: Homepage
	homepage := api.Group("/homepage", logger.New())
	homepage.Get("/videos", handler.GetAllVideos)
	homepage.Get("/video/:id", handler.GetContentById)
	homepage.Get("/audios", handler.GetAllAudios)

	//playlist
	homepage.Get("/playlists", handler.GetAllPlaylists)
	homepage.Get("/playlists/:id", handler.GetPlaylistsVideo)

	//autocomplete search data
	search := api.Group("/search", logger.New())
	search.Get("/info", handler.GetVideoSearchData)

	//Todo: Tags & Categories

	//Todo: Files
	filesystem := api.Group("/storage", logger.New())
	filesystem.Get("/status", handler.StorageStatus)

	//Todo: Patterns
}
