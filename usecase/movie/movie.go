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
	GetMovies(ctx context.Context, filter model.FilterMovie) ([]model.Movie, error)
	GetGenresByIDs(ctx context.Context, genreIDs []uint) ([]model.Genre, error)
	GetArtistsByIDs(ctx context.Context, artistIDs []uint) ([]model.Artist, error)
	GetMoviesByIDs(ctx context.Context, movieIDs []uint) ([]model.Movie, error)
	UpdateMovie(ctx context.Context, movie model.Movie, extraArgs model.UpdateMovieArgs) error
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

func (a *MovieUsecase) GetMovies(ctx context.Context, filter model.FilterMovie) ([]model.Movie, error) {
	movies, err := a.repository.GetMovies(ctx, filter)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Movie{}, echo.NewHTTPError(http.StatusNotFound, "Movie not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch movie")
		return []model.Movie{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return movies, nil
}

func (a *MovieUsecase) UpdateMovie(ctx context.Context, updateMovie model.UpdateMovie) error {
	movies, err := a.repository.GetMoviesByIDs(ctx, []uint{updateMovie.ID})
	if len(movies) != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Movie not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch movie")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	extra := model.UpdateMovieArgs{}
	if updateMovie.Genres != nil {
		genres := []model.Genre{}
		genreIDs := []uint{}
		for _, v := range *updateMovie.Genres {
			genres = append(genres, model.Genre{Model: gorm.Model{ID: v}})
			genreIDs = append(genreIDs, v)
		}
		foundGenres, err := a.repository.GetGenresByIDs(ctx, genreIDs)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
		}
		if err != nil {
			log.Error().Err(err).Msg("error when fetch genre")
			return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
		}
		if len(foundGenres) != len(*updateMovie.Genres) {
			return echo.NewHTTPError(http.StatusBadRequest, "Some of genres are not found, please use available genre")
		}
		movies[0].Genres = genres
		extra.UpdateGenre = true
	}
	if updateMovie.Artists != nil {
		artistIDs := []uint{}
		artist := []model.Artist{}
		for _, v := range *updateMovie.Artists {
			artistIDs = append(artistIDs, v)
			artist = append(artist, model.Artist{Model: gorm.Model{ID: v}})
		}
		foundArtist, err := a.repository.GetArtistsByIDs(ctx, artistIDs)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Genre not found")
		}
		if err != nil {
			log.Error().Err(err).Msg("error when fetch artist")
			return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
		}
		if len(foundArtist) != len(*updateMovie.Artists) {
			return echo.NewHTTPError(http.StatusBadRequest, "Some of artists are not found, please use available artist")
		}
		movies[0].Artists = artist
		extra.UpdateArtist = true
	}
	if updateMovie.Description != nil {
		movies[0].Description = *updateMovie.Description
	}

	if updateMovie.Title != nil {
		movies[0].Title = *updateMovie.Title
	}
	if updateMovie.Duration != nil {
		movies[0].Duration = *updateMovie.Duration
	}
	if updateMovie.WatchURL != nil {
		movies[0].WatchURL = *updateMovie.WatchURL
	}
	err = a.repository.UpdateMovie(ctx, movies[0], extra)
	if err != nil {
		log.Error().Err(err).Msg("error when insert movie")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}
