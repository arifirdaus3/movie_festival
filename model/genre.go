package model

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	Name   string
	Views  int64
	Movies []Movie `gorm:"many2many:movie_genres;"`
}
