package authusecase

import (
	"context"
	"errors"
	"fmt"
	"moviefestival/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	repository repository
}
type repository interface {
	GetUser(ctx context.Context, user model.User) (model.User, error)
	InsertUser(ctx context.Context, user model.User) error
}

func NewAuthUsecase(repository repository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
	}
}

func (a *AuthUsecase) Register(ctx context.Context, user model.User) error {
	foundUser, err := a.repository.GetUser(ctx, user)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("error when get user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while fetching our database, please try again later!")
	}
	if foundUser.Email != "" {
		log.Error().Err(fmt.Errorf("user with email %s already exist", foundUser.Email)).Msg("user already exist")
		return echo.NewHTTPError(http.StatusBadRequest, "User already exist")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Error().Err(err).Msg("error when hashing password")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while hashing password, please try again later!")
	}
	user.Password = string(bytes)
	err = a.repository.InsertUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("error when insert user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while fetching our database, please try again later!")
	}
	return nil
}
