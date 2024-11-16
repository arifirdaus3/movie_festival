package voteusecase

import (
	"context"
	"errors"
	"moviefestival/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type VoteUsecase struct {
	repository repository
}
type repository interface {
	GetUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) (model.UserMovieVote, error)
	InsertUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error
	DeleteUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error
	UpdateUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error
	GetUserVote(ctx context.Context, userID uint) ([]model.UserMovieVote, error)
	GetMostVotedMovie(ctx context.Context) ([]model.VotedMovieCount, error)
}

func NewVoteUsecase(repository repository) *VoteUsecase {
	return &VoteUsecase{
		repository: repository,
	}
}

func (a *VoteUsecase) VoteMovie(ctx context.Context, vote model.UserMovieVote) error {
	userMovieVote, err := a.repository.GetUserMovieVoteIDs(ctx, vote)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("error when get vote")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	if userMovieVote.MovieID != 0 {
		if userMovieVote.Type == vote.Type {
			return nil
		}
		userMovieVote.Type = vote.Type
		err = a.repository.UpdateUserMovieVoteIDs(ctx, userMovieVote)
		if err != nil {
			log.Error().Err(err).Msg("error when update vote")
			return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
		}
		return nil
	}
	err = a.repository.InsertUserMovieVoteIDs(ctx, vote)
	if err != nil {
		log.Error().Err(err).Msg("error when insert vote")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}

func (a *VoteUsecase) UnvoteMovie(ctx context.Context, vote model.UserMovieVote) error {
	userMovieVote, err := a.repository.GetUserMovieVoteIDs(ctx, vote)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User vote not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when get vote")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}

	err = a.repository.DeleteUserMovieVoteIDs(ctx, userMovieVote)
	if err != nil {
		log.Error().Err(err).Msg("error when delete vote")
		return echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}
	return nil
}
func (a *VoteUsecase) GetUserVote(ctx context.Context, userID uint) ([]model.UserMovieVote, error) {
	userMovieVote, err := a.repository.GetUserVote(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.UserMovieVote{}, echo.NewHTTPError(http.StatusNotFound, "User vote not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when get vote")
		return []model.UserMovieVote{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}

	return userMovieVote, nil
}
func (a *VoteUsecase) GetMostVotedMovies(ctx context.Context) ([]model.VotedMovieCount, error) {
	movieVotedCount, err := a.repository.GetMostVotedMovie(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.VotedMovieCount{}, echo.NewHTTPError(http.StatusNotFound, "User vote not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("error when get vote")
		return []model.VotedMovieCount{}, echo.NewHTTPError(http.StatusInternalServerError, "There is error with our database, please try again later!")
	}

	return movieVotedCount, nil
}
