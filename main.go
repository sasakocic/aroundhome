package main

import (
	"database/sql"
	"errors"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"sort"
	"strconv"
	"strings"

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
	Rating             float32
	FlooringExperience string
}

type PartnerWithDistance struct {
	Partner  Partner
	Distance float32
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
	app.Get("/partners/:id", func(ctx *fiber.Ctx) error {
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
	id, err := strconv.ParseInt(c.Params("id"), 10, 16)
	if err != nil {
		return err
	}
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
		err := rows.Scan(&rec.Id, &rec.Name, &rec.Lat, &rec.Lng, &rec.Radius, &rec.Rating, &rec.FlooringExperience)
		if err != nil {
			log.Fatal(err)
		}
		recs = append(recs, rec)
	}
	if err := c.JSON(recs[0]); err != nil {
		return err
	}

	return nil
}

func queryHandler(c *fiber.Ctx, db *sql.DB) error {
	//qString := string(c.Request().URI().QueryString())
	phone := c.Query("phone", "")
	sqm := c.Query("sqm", "")
	address := strings.Split(c.Query("address", "0,0"), ",")
	lat, err := strconv.ParseFloat(address[0], 32)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(address[1], 32)
	if err != nil {
		return err
	}
	material := strings.Split(c.Query("material"), ",")
	sort.Strings(material)
	// prevent SQL injection by strictly checking types of material
	for _, v := range material {
		if v != "carpet" && v != "tiles" && v != "wood" {
			return errors.New("material " + v + " is not allowed in query")
		}
	}
	// manually injecting materials into sql because JSON array is changed during injection into db.Query
	materialString := "'" + strings.Join(material, `','`) + "'"
	rows, err := db.Query(querySql(materialString), lat, lng)
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
		err := rows.Scan(&rec.Partner.Id, &rec.Partner.Name, &rec.Partner.Lat, &rec.Partner.Lng, &rec.Partner.Radius, &rec.Partner.Rating, &rec.Partner.FlooringExperience, &rec.Distance)
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
	response := map[string]interface{}{
		"phone":    phone,
		"partners": recs,
		"sqm":      sqm,
	}
	if err := c.JSON(response); err != nil {
		return err
	}

	return nil
}

func querySql(materialString string) string {
	return fmt.Sprintf("select\n    Id, Name, Lat, Lng, Radius, Rating, flooring_experience AS FlooringExperience,\n    getDistance($1, $2, Lat, Lng) AS Distance\nfrom\n    partners\nwhere\n    getDistance($1, $2, Lat, Lng) < Radius AND flooring_experience @> ARRAY[%s]\norder by\n    Rating DESC,\n    Distance;", materialString)
}

func partnerSql() string {
	return "select\n    Id, Name, Lat, Lng, Radius, Rating, flooring_experience AS FlooringExperience\nfrom\n    partners\nwhere\n    id = $1;"
}
