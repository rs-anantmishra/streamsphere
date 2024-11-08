package handler

import "github.com/gofiber/fiber/v2"

// This removes the need for the Firefox Plugin I use to extract URL's for videos matching
// a certain pattern - but I have to do that for each page, and page no. pattern is available
// in URL so the entire thing can be automated, for generic website sources, where Channel
// download option isn't available.

//Video in Pages
func VideoURLPattern(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}

//Pages in URL's
func SourceURLPattern(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
