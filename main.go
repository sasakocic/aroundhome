package main

import (
	"aroundHome/app"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// @title Fiber Swagger API
// @version 2.0
// @description This is an API server
// @termsOfService http://swagger.io/terms/

// @contact.Name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.Name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {

	// Fiber instance
	webApp := fiber.New()

	// Middleware
	webApp.Use(recover.New())
	webApp.Use(cors.New())

	db := app.DatabaseConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	app.Routes(webApp, db)

	// Start Server
	webPort := os.Getenv("PORT")
	if webPort == "" {
		webPort = "3000"
	}
	if err := webApp.Listen(":" + webPort); err != nil {
		log.Fatal(err)
	}
}
