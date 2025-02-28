package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"song_library/config"
	"song_library/db"
	_ "song_library/docs"
	"song_library/handlers"
)

// @title Song Library API
// @version 1.0
// @description This is a simple API to manage a song library.
// @host localhost:8080
// @BasePath /

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg := config.LoadConfig()

	db.InitDB(cfg)
	defer db.Db.Close()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.GET("/songs", handlers.GetSongs)
	router.GET("/songs/lyrics/:song_id", handlers.GetSongLyrics)
	router.DELETE("/songs/:song_id", handlers.DeleteSong)
	router.PUT("/songs/:song_id", handlers.UpdateSong)
	router.POST("/songs", handlers.AddSong(cfg))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info().Msgf("Backend API running on port %s", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
