package dto

type Genre struct {
	GenreName string `json:"genreName" db:"genre_name"`
}

type MovieGenres struct {
	IDMovie int `json:"idmovie" db:"id_movie" binding:"required"`
	IDGenre int `json:"idgenre" db:"id_genre" binding:"required"`
}