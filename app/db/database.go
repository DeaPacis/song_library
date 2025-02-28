package db

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"song_library/config"
)

var Db *sql.DB

func InitDB(cfg *config.Config) {
	var err error

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	Db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection failed")
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Database ping failed")
	}
	log.Info().Msg("Connected to database!")
}
