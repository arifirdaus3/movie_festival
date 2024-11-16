package genreusecase

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

type GenreUsecase struct {
	repository repository
}
type repository interface {
	InsertGenre(ctx context.Context, genre []model.Genre) error
	GetGenres(ctx context.Context, pagination model.Pagination) ([]model.Genre, error)
}

func NewGenreUsecase(repository repository) *GenreUsecase {
	return &GenreUsecase{
		repository: repository,
	}
}

func (a *GenreUsecase) InsertGenre(ctx context.Context, genres []model.Genre) error {
	err := a.repository.InsertGenre(ctx, genres)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusBadRequest, "Genre already exist!")
		}
		log.Error().Err(err).Msg("error when insert genre")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}

func (a *GenreUsecase) GetGenres(ctx context.Context, pagination model.Pagination) ([]model.Genre, error) {
	genres, err := a.repository.GetGenres(ctx, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Genre{}, echo.NewHTTPError(http.StatusNotFound, "Genre not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch genre")
		return []model.Genre{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return genres, nil
}
