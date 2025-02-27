package db

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"song_library/app/config"
)

var Db *sql.DB

func InitDB(cfg *config.Config) {
	var err error

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	Db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to database!")
}
