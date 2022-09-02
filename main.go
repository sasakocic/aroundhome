package main

import (
	"database/sql"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"strconv"

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
}

type PartnerWithDistance struct {
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
	db := databaseConnect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	// Routes
	app.Get("/", HealthCheck)
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/partners/*", func(ctx *fiber.Ctx) error {
		return partnersHandler(ctx, db)
	})
	app.Get("/query/*", func(ctx *fiber.Ctx) error {
		return queryHandler(ctx, db)
	})

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

func partnersHandler(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.ParseInt(c.Query("id"), 10, 16)
	if err != nil {
		return err
	}
	//result, err := db.Exec(sqlQuery)
	rows, err := db.Query(partnerSql(), id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	recs := make([]*Partner, 0)
	for rows.Next() {
		rec := new(Partner)
		err := rows.Scan(&rec.Id, &rec.Name, &rec.Lat, &rec.Lng, &rec.Radius, &rec.Sqm, &rec.Rating, &rec.FlooringExperience)
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

	if err := c.JSON(recs[0]); err != nil {
		return err
	}

	return nil
}

func queryHandler(c *fiber.Ctx, db *sql.DB) error {
	//qString := string(c.Request().URI().QueryString())
	lat, err := strconv.ParseFloat(c.Query("lat"), 32)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 32)
	if err != nil {
		return err
	}
	//fmt.Println(qString)
	//log.Fatal(qString)
	sqlQuery := querySql(lat, lng)
	//if err := c.JSON(sqlQuery); err != nil {
	//	return err
	//}

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	recs := make([]*PartnerWithDistance, 0)
	for rows.Next() {
		rec := new(PartnerWithDistance)
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

func querySql(lat, lng float64) string {
	return fmt.Sprintf(
		"select\n    Id, Name, Lat, Lng, Radius, Sqm, Rating, flooring_experience AS FlooringExperience,\n    getDistance(%f, %f, Lat, Lng) AS Distance\nfrom\n    partners\nwhere\n    getDistance(%f, %f, Lat, Lng) < Radius\norder by\n    Rating DESC,\n    Distance;", lat, lng, lat, lng)
}

func partnerSql() string {
	return "select\n    Id, Name, Lat, Lng, Radius, Sqm, Rating, flooring_experience AS FlooringExperience\nfrom\n    partners\nwhere\n    id = $1;"
}
