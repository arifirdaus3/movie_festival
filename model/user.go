package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string
	Email          string
	Password       string
	Movies         []Movie `gorm:"many2many:user_movies;"`
	UserMovieVotes []UserMovieVote
}
