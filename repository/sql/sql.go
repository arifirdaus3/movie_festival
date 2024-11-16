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
func (r *RepositorySQL) GetArtists(ctx context.Context, pagination model.Pagination) ([]model.Artist, error) {
	var artists []model.Artist
	err := r.db.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&artists).Error
	return artists, err
}
func (r *RepositorySQL) InsertArtist(ctx context.Context, artist []model.Artist) error {
	return r.db.Create(&artist).Error
}
func (r *RepositorySQL) GetArtistsByIDs(ctx context.Context, artistIDs []uint) ([]model.Artist, error) {
	var artists []model.Artist
	err := r.db.Find(&artists, artistIDs).Error
	return artists, err
}
func (r *RepositorySQL) GetGenres(ctx context.Context, pagination model.Pagination) ([]model.Genre, error) {
	var genres []model.Genre
	err := r.db.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&genres).Error
	return genres, err
}
func (r *RepositorySQL) GetGenresByIDs(ctx context.Context, genreIDs []uint) ([]model.Genre, error) {
	var genres []model.Genre
	err := r.db.Find(&genres, genreIDs).Error
	return genres, err
}
func (r *RepositorySQL) InsertGenre(ctx context.Context, genre []model.Genre) error {
	return r.db.Create(&genre).Error
}
