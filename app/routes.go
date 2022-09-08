package app

import (
	"aroundHome/app/controllers"
	"database/sql"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, db *sql.DB) {
	// Routes
	app.Get("/", controllers.HealthCheck)
	//app.Get("/swagger/*", swagger.HandlerDefault)     // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "https://raw.githubusercontent.com/sasakocic/aroundhome/master/docs/swagger.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))
	app.Get("/partners/:id", func(ctx *fiber.Ctx) error {
		return controllers.PartnersHandler(ctx, db)
	})
	app.Get("/query/*", func(ctx *fiber.Ctx) error {
		return controllers.QueryHandler(ctx, db)
	})
}
