package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"song_library/config"
	"song_library/db"
	"song_library/models"
	"strconv"
	"strings"
)

// GetSongs returns a list of songs with filtering and pagination
// @Summary Get a list of songs
// @Description Returns a list of songs with optional filters by group name, song name, and release date
// @Tags Songs
// @Produce json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param releaseDate query string false "Release date (format: DD.MM.YYYY)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of songs per page (default: 5)"
// @Success 200 {array} models.Song
// @Failure 500 {object} map[string]string "Database error"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	log.Debug().Msg("Processing GetSongs request")
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

	log.Debug().Msgf("Executing database query: %s with args %v", query, args)

	rows, err := db.Db.Query(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Database request error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Text, &s.Link)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning song row")
			continue
		}
		songs = append(songs, s)
	}
	log.Info().Msgf("Found %d songs", len(songs))
	c.JSON(http.StatusOK, songs)
}

// GetSongLyrics returns song lyrics with pagination by verses
// @Summary Get song lyrics
// @Description Returns the lyrics of a song, split into verses, with pagination support
// @Tags Lyrics
// @Produce json
// @Param song_id path int true "Song ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of verses per page (default: 1)"
// @Success 200 {array} string
// @Failure 400 {object} map[string]string "Invalid song ID"
// @Failure 404 {object} map[string]string "Song not found"
// @Router /songs/lyrics/{song_id} [get]
func GetSongLyrics(c *gin.Context) {
	log.Debug().Msg("Processing GetSongLyrics request")
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var lyrics string
	err = db.Db.QueryRow("SELECT lyrics FROM songs WHERE song_id = $1", songID).Scan(&lyrics)
	if err != nil {
		log.Error().Err(err).Msgf("Song with ID %d not found", songID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}

	verses := strings.Split(lyrics, "\n\n")
	if len(verses) == 1 {
		verses = strings.Split(lyrics, "\n")
	}

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr == "" && limitStr == "" {
		c.JSON(http.StatusOK, verses)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
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

	log.Info().Msgf("Returning lyrics for song ID %d", songID)
	c.JSON(http.StatusOK, verses[start:end])
}

// DeleteSong removes a song by ID
// @Summary Delete a song
// @Description Deletes a song by the specified ID
// @Tags Songs
// @Param song_id path int true "Song ID"
// @Success 200 {object} map[string]string "Song deleted successfully"
// @Failure 400 {object} map[string]string "Invalid song ID"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /songs/{song_id} [delete]
func DeleteSong(c *gin.Context) {
	log.Debug().Msg("Processing DeleteSong request")
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	result, err := db.Db.Exec("DELETE FROM songs WHERE song_id = $1", songID)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		log.Warn().Msgf("Song with ID %d not found", songID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}

	log.Info().Msgf("Song with ID %d deleted", songID)
	c.JSON(http.StatusOK, gin.H{"message": "Song was deleted"})
}

// UpdateSong updates song details
// @Summary Update a song
// @Description Updates song details by its ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param song_id path int true "Song ID"
// @Param song body models.Song true "Updated song details"
// @Success 200 {object} map[string]string "Song updated successfully"
// @Failure 400 {object} map[string]string "Invalid data format"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /songs/{song_id} [put]
func UpdateSong(c *gin.Context) {
	log.Debug().Msg("Processing UpdateSong request")
	songID, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var s models.Song
	if err := c.ShouldBindJSON(&s); err != nil {
		log.Error().Err(err).Msg("Error binding JSON")
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
		log.Error().Err(err).Msg("Error updating song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else if affectedRows == 0 {
		log.Warn().Msgf("Song with ID %d not found", songID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song is not found"})
		return
	}

	log.Info().Msgf("Song with ID %d updated", songID)
	c.JSON(http.StatusOK, gin.H{"message": "Song was updated"})
}

// AddSong adds a new song using external API data
// @Summary Add a song
// @Description Adds a new song by fetching its details from an external API
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body object true "Song data (group, song)"
// @Success 200 {object} map[string]int "Added song ID"
// @Failure 400 {object} map[string]string "Invalid request or missing song information"
// @Failure 500 {object} map[string]string "Database error"
// @Router /songs [post]
func AddSong(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Msg("Processing AddSong request")
		var input struct {
			Group string `json:"group" binding:"required"`
			Song  string `json:"song" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Error().Err(err).Msg("Error binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
			return
		}

		groupEncoded := strings.ReplaceAll(input.Group, " ", "+")
		songEncoded := strings.ReplaceAll(input.Song, " ", "+")

		apiURL := fmt.Sprintf("%s/info?group=%s&song=%s", cfg.ExternalAPIURL, groupEncoded, songEncoded)
		log.Info().Msgf("Calling external API: %s", apiURL)
		resp, err := http.Get(apiURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Error().Err(err).Msg("External API call failed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get song info"})
			return
		}
		defer resp.Body.Close()

		var detail models.SongDetail
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			log.Error().Err(err).Msg("Error decoding external API response")
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
			log.Error().Err(err).Msg("Error inserting song into database")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Info().Msgf("Song with ID %d added", songID)
		c.JSON(http.StatusOK, gin.H{"song_id": songID})
	}
}
