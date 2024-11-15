package authusecase

import (
	"context"
	"errors"
	"fmt"
	"moviefestival/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	repository repository
	config     *model.Config
}
type repository interface {
	GetUser(ctx context.Context, user model.User) (model.User, error)
	InsertUser(ctx context.Context, user model.User) error
}

func NewAuthUsecase(config *model.Config, repository repository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
		config:     config,
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

func (a *AuthUsecase) Login(ctx context.Context, user model.User) (model.Token, error) {
	foundUser, err := a.repository.GetUser(ctx, user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Token{}, echo.NewHTTPError(http.StatusNotFound, "User not found!")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when get user")
		return model.Token{}, echo.NewHTTPError(http.StatusInternalServerError, "Error while fetching our database, please try again later!")
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return model.Token{}, echo.NewHTTPError(http.StatusUnauthorized, "Wrong password")
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(a.config.AccessTokenExpirationMinute))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:    foundUser.Name,
		Email:   foundUser.Email,
		IsAdmin: foundUser.IsAdmin,
	})
	signedAccessToken, err := accessToken.SignedString([]byte(a.config.SignTokenSecret))
	if err != nil {
		log.Error().Err(err).Msg("error when signing token")
		return model.Token{}, echo.NewHTTPError(http.StatusInternalServerError, "There is problem in our server, please try again later!")
	}
	refreshTokenExpiration := time.Now().Add(time.Minute * time.Duration(a.config.RefreshTokenExpirationMinute))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:    foundUser.Name,
		Email:   foundUser.Email,
		IsAdmin: foundUser.IsAdmin,
	})
	signeRefreshToken, err := refreshToken.SignedString([]byte(a.config.SignTokenSecret))
	if err != nil {
		log.Error().Err(err).Msg("error when signing token")
		return model.Token{}, echo.NewHTTPError(http.StatusInternalServerError, "There is problem in our server, please try again later!")
	}
	return model.Token{AcceessToken: signedAccessToken, RefreshToken: signeRefreshToken, RefreshTokenExpiration: refreshTokenExpiration}, nil
}

func (a *AuthUsecase) Refresh(ctx context.Context, token string) (model.Token, error) {
	verifiedToken, err := jwt.ParseWithClaims(token, &model.CustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.config.SignTokenSecret), nil
	})
	if err != nil {
		return model.Token{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}
	claim, ok := verifiedToken.Claims.(*model.CustomClaim)
	if !ok {
		return model.Token{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid claim")
	}
	if claim.ExpiresAt.Before(time.Now()) {
		return model.Token{}, echo.NewHTTPError(http.StatusUnauthorized, "Expired refresh token")
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(a.config.AccessTokenExpirationMinute))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:    claim.Name,
		Email:   claim.Email,
		IsAdmin: claim.IsAdmin,
	})
	signedAccessToken, err := accessToken.SignedString([]byte(a.config.SignTokenSecret))
	if err != nil {
		log.Error().Err(err).Msg("error when signing token")
		return model.Token{}, echo.NewHTTPError(http.StatusInternalServerError, "There is problem in our server, please try again later!")
	}
	return model.Token{
		AcceessToken: signedAccessToken,
	}, nil
}
