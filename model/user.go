package model

import "gorm.io/gorm"

type User struct {
	gorm.Model     `json:"-"`
	Name           string          `json:"name"`
	Email          string          `json:"email" form:"username"`
	Password       string          `json:"password" form:"password"`
	IsAdmin        bool            `json:"isAdmin"`
	Movies         []Movie         `gorm:"many2many:user_movies;" json:"-"`
	UserMovieVotes []UserMovieVote `json:"-"`
}
