package models

type CreateSong struct {
	Group string `json:"group"       validate:"required,min=2,max=20"`
	Song  string `json:"song"        validate:"required,min=2"`
}

type RequestGetAll struct {
	Id     string `json:"id"`
	Limit  string `json:"limit"`
	Filter string `json:"filter"`
	Value  string `json:"value"`
}

type SongsResponse struct {
	Id          int    `json:"id" db:"id"`
	GroupSong   string `json:"group_song" db:"group_song"`
	Song        string `json:"song" db:"song"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}
