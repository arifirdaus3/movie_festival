package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title          string
	Description    string
	Duration       int64
	WatchURL       string
	Users          []User   `gorm:"many2many:user_movies;"`
	Artists        []Artist `gorm:"many2many:movie_artists;"`
	Genres         []Genre  `gorm:"many2many:movie_genres;"`
	UserMovieVotes []UserMovieVote
}

type CreateMovie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	WatchURL    string `json:"watchURL"`
	Genres      []uint `json:"genres"`
	Artists     []uint `json:"artists"`
}

func (c *CreateMovie) ToMovie() Movie {
	genre := []Genre{}
	for _, v := range c.Genres {
		genre = append(genre, Genre{Model: gorm.Model{ID: v}})
	}
	artist := []Artist{}
	for _, v := range c.Artists {
		artist = append(artist, Artist{Model: gorm.Model{ID: v}})
	}
	return Movie{
		Title:       c.Title,
		Description: c.Description,
		Duration:    c.Duration,
		WatchURL:    c.WatchURL,
		Genres:      genre,
		Artists:     artist,
	}
}
