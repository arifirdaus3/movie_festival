package movieusecase

import (
	"context"
	"errors"
	"moviefestival/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type MovieUsecase struct {
	repository repository
}
type repository interface {
	InsertMovie(ctx context.Context, movie model.Movie) error
	GetMovies(ctx context.Context, pagination model.Pagination) ([]model.Movie, error)
	GetGenresByIDs(ctx context.Context, genreIDs []uint) ([]model.Genre, error)
	GetArtistsByIDs(ctx context.Context, artistIDs []uint) ([]model.Artist, error)
}

func NewMovieUsecase(repository repository) *MovieUsecase {
	return &MovieUsecase{
		repository: repository,
	}
}

func (a *MovieUsecase) InsertMovie(ctx context.Context, movies model.Movie) error {
	genreIDs := []uint{}
	for _, v := range movies.Genres {
		genreIDs = append(genreIDs, v.ID)
	}
	genres, err := a.repository.GetGenresByIDs(ctx, genreIDs)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch genre")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	if len(genres) != len(movies.Genres) {
		return echo.NewHTTPError(http.StatusBadRequest, "Some of genres are not found, please use available genre")
	}

	artistIDs := []uint{}
	for _, v := range movies.Artists {
		artistIDs = append(artistIDs, v.ID)
	}
	artists, err := a.repository.GetArtistsByIDs(ctx, artistIDs)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch artist")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	if len(artists) != len(movies.Artists) {
		return echo.NewHTTPError(http.StatusBadRequest, "Some of artists are not found, please use available artist")
	}

	err = a.repository.InsertMovie(ctx, movies)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusBadRequest, "Movie already exist!")
		}
		log.Error().Err(err).Msg("error when insert movie")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}

func (a *MovieUsecase) GetMovies(ctx context.Context, pagination model.Pagination) ([]model.Movie, error) {
	movies, err := a.repository.GetMovies(ctx, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Movie{}, echo.NewHTTPError(http.StatusNotFound, "Movie not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch movie")
		return []model.Movie{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return movies, nil
}
