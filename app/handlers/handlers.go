package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"song_library/config"
	"song_library/db"
	"song_library/models"
	"strconv"
	"strings"
)

func GetSongs(c *gin.Context) {
	filters := []string{}
	args := []interface{}{}
	i := 1

	if group := c.Query("group"); group != "" {
		filters = append(filters, fmt.Sprintf("group_name ILIKE $%d", i))
		args = append(args, "%"+group+"%")
		i++
	}
	if song := c.Query("song"); song != "" {
		filters = append(filters, fmt.Sprintf("song_name ILIKE $%d", i))
		args = append(args, "%"+song+"%")
		i++
	}
	if releaseDate := c.Query("releaseDate"); releaseDate != "" {
		filters = append(filters, fmt.Sprintf("release_date = $%d", i))
		args = append(args, releaseDate)
		i++
	}

	query := "SELECT song_id, group_name, song_name, release_date, lyrics, link FROM songs"
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY song_id LIMIT %d OFFSET %d", limit, offset)

	log.Debugf("Database request: %s with args %v", query, args)

	rows, err := db.Db.Query(query, args...)
	if err != nil {
		log.Error("Database request error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Text, &s.Link)
		if err != nil {
			log.Error(err)
			continue
		}
		songs = append(songs, s)
	}
	c.JSON(http.StatusOK, songs)
}

func GetSongLyrics(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error("Invalid song ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var lyrics string
	err = db.Db.QueryRow("SELECT lyrics FROM songs WHERE song_id = $1", songID).Scan(&lyrics)
	if err != nil {
		log.Error("Error getting song lyrics: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}

	verses := strings.Split(lyrics, "\n\n")
	if len(verses) == 1 {
		verses = strings.Split(lyrics, "\n")
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit < 1 {
		limit = 1
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		c.JSON(http.StatusOK, []string{})
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	c.JSON(http.StatusOK, verses[start:end])
}

func DeleteSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error("Invalid song ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	result, err := db.Db.Exec("DELETE FROM songs WHERE song_id = $1", songID)
	if err != nil {
		log.Error("Error deleting a song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		log.Errorf("Song with id = %d is not found", songID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song was deleted"})
}

func UpdateSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error("Invalid song ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var s models.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		log.Error("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}
	query := `
		UPDATE songs 
		SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5
		WHERE song_id = $6
	`
	result, err := db.Db.Exec(query, s.Group, s.Song, s.ReleaseDate, s.Text, s.Link, songID)
	affectedRows, _ := result.RowsAffected()
	if err != nil {
		log.Error("Error updating a song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else if affectedRows == 0 {
		log.Errorf("Song with id = %d is not found", songID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song was updated"})
}

func AddSong(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Group string `json:"group" binding:"required"`
			Song  string `json:"song" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Error("Error binding JSON: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
			return
		}

		apiURL := fmt.Sprintf("%s/info?group=%s&song=%s", cfg.ExternalAPIURL, input.Group, input.Song)
		log.Infof("External API call: %s", apiURL)
		resp, err := http.Get(apiURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Error("External API call error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get song info"})
			return
		}
		defer resp.Body.Close()

		var detail models.SongDetail
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			log.Error("Error decoding external API answer: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Song info processing error"})
			return
		}

		query := `
			INSERT INTO songs (group_name, song_name, release_date, lyrics, link)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING song_id
		`
		var songID int
		err = db.Db.QueryRow(query, input.Group, input.Song, detail.ReleaseDate, detail.Text, detail.Link).Scan(&songID)
		if err != nil {
			log.Error("Error adding a song: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		log.Infof("Song with ID = %d was added", songID)
		c.JSON(http.StatusOK, gin.H{"song_id": songID})
	}
}
