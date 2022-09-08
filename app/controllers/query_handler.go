package controllers

import (
	"aroundHome/app/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
)

// QueryHandler godoc
// @Summary Get list of partners that satisfy given query.
// @Description Returns list of partners that satisfy given query.
// @Tags query
// @Accept */*
// @Produce json
// @Param phone  query string false "Phone number for contact" example(01604323444)
// @Param sqm  query decimal false "Square meters" example(65.22)
// @Param address  query string true "Address in format: Latitude,Longitude" example(40.076763,113.30013)
// @Param material query []string true "Material collection: carpet,tiles,wood" collectionFormat(csv) example(carpet,tiles,wood)
// @Success 200 {object} map[string]interface{}
// @Router /query/{id} [get]
func QueryHandler(c *fiber.Ctx, db *sql.DB) error {
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

	recs := make([]*models.PartnerWithDistance, 0)
	for rows.Next() {
		rec := new(models.PartnerWithDistance)
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
