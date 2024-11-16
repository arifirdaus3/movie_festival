package artistusecase

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

type ArtistUsecase struct {
	repository repository
}
type repository interface {
	InsertArtist(ctx context.Context, artist []model.Artist) error
	GetArtists(ctx context.Context, pagination model.Pagination) ([]model.Artist, error)
}

func NewArtistUsecase(repository repository) *ArtistUsecase {
	return &ArtistUsecase{
		repository: repository,
	}
}

func (a *ArtistUsecase) InsertArtist(ctx context.Context, artists []model.Artist) error {
	err := a.repository.InsertArtist(ctx, artists)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusBadRequest, "Artist already exist!")
		}
		log.Error().Err(err).Msg("error when insert artist")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}
func (a *ArtistUsecase) GetArtists(ctx context.Context, pagination model.Pagination) ([]model.Artist, error) {
	artists, err := a.repository.GetArtists(ctx, pagination)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Artist{}, echo.NewHTTPError(http.StatusNotFound, "Artist not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetch artist")
		return []model.Artist{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return artists, nil
}
