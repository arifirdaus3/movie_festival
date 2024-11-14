package repositorysql

import (
	"context"
	"moviefestival/model"

	"gorm.io/gorm"
)

type RepositorySQL struct {
	db *gorm.DB
}

func NewRepositorySQL(db *gorm.DB) *RepositorySQL {
	return &RepositorySQL{
		db: db,
	}
}

func (r *RepositorySQL) GetUser(ctx context.Context, user model.User) (model.User, error) {
	var foundUser model.User
	res := r.db.First(&foundUser, "email = ?", user.Email)
	return foundUser, res.Error
}
func (r *RepositorySQL) InsertUser(ctx context.Context, user model.User) error {
	return r.db.Create(&user).Error
}
