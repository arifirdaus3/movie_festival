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
func (r *RepositorySQL) InsertMovie(ctx context.Context, movie model.Movie) error {
	return r.db.Create(&movie).Error
}
func (r *RepositorySQL) UpdateMovie(ctx context.Context, movie model.Movie, extraArgs model.UpdateMovieArgs) error {
	err := r.db.Save(&movie).Error
	if err != nil {
		return err
	}
	if extraArgs.UpdateGenre {
		err = r.db.Model(&movie).Association("Genres").Replace(movie.Genres)
	}
	if extraArgs.UpdateArtist {
		err = r.db.Model(&movie).Association("Artists").Replace(movie.Artists)
	}
	return err
}
func (r *RepositorySQL) GetMovies(ctx context.Context, pagination model.Pagination) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Offset((pagination.Page - 1) * pagination.Limit).Limit(pagination.Limit).Find(&movies).Error
	return movies, err
}
func (r *RepositorySQL) GetMoviesByIDs(ctx context.Context, movieIDs []uint) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Find(&movies, movieIDs).Error
	return movies, err
}
