package app

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Routes
	app.Get("/", HealthCheck)
	//app.Get("/swagger/*", swagger.HandlerDefault)     // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "https://raw.githubusercontent.com/sasakocic/aroundhome/master/docs/swagger.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))
	app.Get("/partners/:id", func(ctx *fiber.Ctx) error {
		return partnersHandler(ctx, db)
	})
	app.Get("/query/*", func(ctx *fiber.Ctx) error {
		return queryHandler(ctx, db)
	})
}
