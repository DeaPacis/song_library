package models

type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Lyrics      string `json:"lyrics,omitempty"`
	Link        string `json:"link,omitempty"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}
