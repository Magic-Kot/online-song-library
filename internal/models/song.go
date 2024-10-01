package models

type User struct {
	ID       int    `json:"id"           validate:"required,min=1"`
	Age      int    `json:"age"          validate:"lte=120"` //gte=14,
	Username string `json:"login"        validate:"max=20"`  //min=4,
	Name     string `json:"name"         validate:"min=1,max=20"`
	Surname  string `json:"surname"      validate:"min=1,max=20"`
	Email    string `json:"email"        validate:"email"`
	Avatar   string `json:"avatar"`
}

type CreateSong struct {
	Group string `json:"group"       validate:"required,min=2,max=20"`
	Song  string `json:"song"        validate:"required,min=2"`
}

type RequestGetAll struct {
	Filter string `json:"filter"`
	Value  string `json:"value"`
}

type SongsResponse struct {
	Id        int    `json:"id" db:"id"`
	GroupSong string `json:"group_song" db:"group_song"`
	Song      string `json:"song" db:"song"`
}
