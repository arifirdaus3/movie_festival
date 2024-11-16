package model

import "gorm.io/gorm"

type Genre struct {
	gorm.Model `json:"-"`
	Name       string  `gorm:"uniqueIndex" json:"name"`
	Views      int64   `gorm:"default:0"`
	Movies     []Movie `gorm:"many2many:movie_genres;" json:"-"`
}
type GenreHTTPResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewGenreHTTPResponse(a Genre) GenreHTTPResponse {
	return GenreHTTPResponse{
		Name: a.Name,
		ID:   a.ID,
	}
}
