package model

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name   string
	Movies []Movie `gorm:"many2many:user_movies;"`
}
