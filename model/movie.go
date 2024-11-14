package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title          string
	Description    string
	Duration       string
	Views          int64
	WatchURL       string
	Users          []User   `gorm:"many2many:user_movies;"`
	Artists        []Artist `gorm:"many2many:movie_artists;"`
	Genre          []Genre  `gorm:"many2many:movie_genres;"`
	UserMovieVotes []UserMovieVote
}
