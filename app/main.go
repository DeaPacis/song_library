package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"song_library/config"
	"song_library/db"
	"song_library/handlers"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg)
	defer db.Db.Close()

	router := gin.Default()

	router.GET("/songs", handlers.GetSongs)
	router.GET("/songs/:song_id/lyrics", handlers.GetSongLyrics)
	router.DELETE("/songs/:song_id", handlers.DeleteSong)
	router.PUT("/songs/:song_id", handlers.UpdateSong)
	router.POST("/songs", handlers.AddSong(cfg))

	log.Infof("Backend API running on port %s", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
