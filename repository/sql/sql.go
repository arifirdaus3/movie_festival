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
func (r *RepositorySQL) GetUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) (model.UserMovieVote, error) {
	var userMovieVote model.UserMovieVote
	err := r.db.Find(&userMovieVote, "user_id=? AND movie_id=?", vote.UserID, vote.MovieID).Error
	return userMovieVote, err
}
func (r *RepositorySQL) InsertUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error {
	return r.db.Create(&vote).Error
}
func (r *RepositorySQL) UpdateUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error {
	return r.db.Model(&model.UserMovieVote{}).Where("movie_id = ? AND user_id = ?", vote.MovieID, vote.UserID).Update("type", vote.Type).Error
}
func (r *RepositorySQL) DeleteUserMovieVoteIDs(ctx context.Context, vote model.UserMovieVote) error {
	return r.db.Model(&model.UserMovieVote{}).Where("movie_id = ? AND user_id = ?", vote.MovieID, vote.UserID).Delete(&vote).Error
}
func (r *RepositorySQL) GetUserVote(ctx context.Context, userID uint) ([]model.UserMovieVote, error) {
	var userMovieVote []model.UserMovieVote
	err := r.db.Find(&userMovieVote, "user_id=? ", userID).Error
	return userMovieVote, err
}
func (r *RepositorySQL) GetMostVotedMovie(ctx context.Context) ([]model.VotedMovieCount, error) {
	rows, err := r.db.Model(&model.UserMovieVote{}).Select("movie_id, count(*) as total").Group("movie_id").Order("total desc").Rows()
	if err != nil {
		return []model.VotedMovieCount{}, err
	}
	var votedMovie []model.VotedMovieCount
	defer rows.Close()
	for rows.Next() {
		var (
			id    uint
			count int64
		)
		err := rows.Scan(&id, &count)
		if err != nil {
			return votedMovie, err
		}
		votedMovie = append(votedMovie, model.VotedMovieCount{MovieID: id, Count: count})
	}
	return votedMovie, err
}
