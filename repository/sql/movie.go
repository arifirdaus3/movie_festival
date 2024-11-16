package repositorysql

import (
	"context"
	"moviefestival/model"
)

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
func (r *RepositorySQL) GetMovies(ctx context.Context, filter model.FilterMovie) ([]model.Movie, error) {
	var movies []model.Movie
	db := r.db.Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Preload("Artists").Preload("Genres")
	if filter.Search != "" {
		filter.Search = "%" + filter.Search + "%"
		switch filter.SearchBy {
		case "description":
			db = db.Where("description ILIKE ? ", filter.Search)
		case "artist":
			db = db.Where("artists.name ILIKE ? ", filter.Search)
		case "genre":
			db = db.Where("genres.name ILIKE ? ", filter.Search)
		default:
			db = db.Where("title ILIKE ? ", filter.Search)
		}
	}
	err := db.
		Joins("left join movie_artists ma on ma.movie_id = movies.id").
		Joins("left join movie_genres mg on mg.movie_id = movies.id").
		Joins("left join artists on artists.id = ma.artist_id").
		Joins("left join genres on genres.id = mg.genre_id").Find(&movies).Error
	return movies, err
}
func (r *RepositorySQL) GetMoviesByIDs(ctx context.Context, movieIDs []uint) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Find(&movies, movieIDs).Preload("Genres").Error
	return movies, err
}

func (r *RepositorySQL) InsertMovieUsers(ctx context.Context, user model.User, movie model.Movie) error {
	err := r.db.Model(&movie).Association("Users").Append(&user)
	if err != nil {
		return err
	}
	count := r.db.Model(&movie).Association("Users").Count()
	err = r.db.Model(&movie).Where("id = ?", movie.ID).Update("views", count).Error
	if err != nil {
		return err
	}

	return nil
}
