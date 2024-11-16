package model

import "gorm.io/gorm"

type Artist struct {
	gorm.Model `json:"-"`
	Name       string  `gorm:"uniqueIndex"`
	Movies     []Movie `gorm:"many2many:user_movies;" json:"-"`
}

type ArtistHTTPResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewArtistHTTPResponse(a Artist) ArtistHTTPResponse {
	return ArtistHTTPResponse{
		Name: a.Name,
		ID:   a.ID,
	}
}
