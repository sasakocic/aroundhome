package main

import (
	"database/sql"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "aroundhome"
)

type Partner struct {
	Id                 int16
	Name               string
	Lat                float32
	Lng                float32
	Radius             float32
	Sqm                float32
	Rating             float32
	FlooringExperience string
	Distance           float32
}

// @title Fiber Swagger Example API
// @version 2.0
// @description This is a sample server server.
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
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	// Routes
	app.Get("/", HealthCheck)
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/partners/*", Partners)

	// Start Server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func databaseConnect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

func Partners(c *fiber.Ctx) error {
	db := databaseConnect()
	defer db.Close()
	sqlQuery := buildSql(40.076762, 113.300129)
	//result, err := db.Exec(sqlQuery)
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	recs := make([]*Partner, 0)
	for rows.Next() {
		rec := new(Partner)
		err := rows.Scan(&rec.Id, &rec.Name, &rec.Lat, &rec.Lng, &rec.Radius, &rec.Sqm, &rec.Rating, &rec.FlooringExperience, &rec.Distance)
		//e, err := json.Marshal(rec)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//log.Fatal(e)
		if err != nil {
			log.Fatal(err)
		}
		recs = append(recs, rec)
	}
	//e, err := json.Marshal(recs)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(e)

	//log.Fatal(recs)

	//res := map[string]interface{}{
	//	"test": "ok",
	//}

	if err := c.JSON(recs); err != nil {
		return err
	}

	return nil
}

func buildSql(lat, lng float32) string {
	return fmt.Sprintf(
		"select\n    Id, Name, Lat, Lng, Radius, Sqm, Rating, flooring_experience AS FlooringExperience,\n    getDistance(%f, %f, Lat, Lng) AS Distance\nfrom\n    partners\nwhere\n    getDistance(40.076762, 113.300129, Lat, Lng) < Radius\norder by\n    Rating DESC,\n    Distance;", lat, lng)
}
