package model

import "gorm.io/gorm"

type UserMovieVote struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID;references:ID"`
	MovieID uint
	Movie   Movie `gorm:"foreignKey:MovieID;references:ID"`
	Type    bool  // true upvote, false downvote
}
