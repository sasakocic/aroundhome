package controllers

import (
	"aroundHome/app/models"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// PartnersHandler godoc
// @Summary Get partners data for a given id.
// @Description Returns partners data for an id as integer.
// @Tags partners
// @Accept */*
// @Produce json
// @Param id  path int true "Partner ID"
// @Success 200 {object} map[string]interface{}
// @Router /partners/{id} [get]
func PartnersHandler(c *fiber.Ctx, db *sql.DB) error {
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
	recs := make([]*models.Partner, 0)
	for rows.Next() {
		rec := new(models.Partner)
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

func partnerSql() string {
	return "select\n    Id, Name, Lat, Lng, Radius, Rating, flooring_experience AS FlooringExperience\nfrom\n    partners\nwhere\n    id = $1;"
}
