package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func DatabaseConnect() *sql.DB {
	host := os.Getenv("PG_HOSTNAME")
	if host == "" {
		host = "localhost"
	}
	pgPort := os.Getenv("PG_PORT")
	if pgPort == "" {
		pgPort = "5432"
	}
	port, err := strconv.ParseInt(pgPort, 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("PG_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("PG_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("PG_DATABASE")
	if dbname == "" {
		dbname = "aroundhome"
	}
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
