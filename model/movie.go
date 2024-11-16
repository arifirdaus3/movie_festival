package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title          string
	Description    string
	Duration       int64
	WatchURL       string
	Users          []User   `gorm:"many2many:user_movies;"`
	Artists        []Artist `gorm:"many2many:movie_artists;"`
	Genres         []Genre  `gorm:"many2many:movie_genres;"`
	UserMovieVotes []UserMovieVote
}
type MovieHTTPResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Duration    int64                `json:"duration"`
	WatchURL    string               `json:"watchURL"`
	Genres      []GenreHTTPResponse  `json:"genres"`
	Artists     []ArtistHTTPResponse `json:"artists"`
}

type CreateMovie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	WatchURL    string `json:"watchURL"`
	Genres      []uint `json:"genres"`
	Artists     []uint `json:"artists"`
}
type UpdateMovie struct {
	ID          uint    `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Duration    *int64  `json:"duration,omitempty"`
	WatchURL    *string `json:"watchURL,omitempty"`
	Genres      *[]uint `json:"genres,omitempty"`
	Artists     *[]uint `json:"artists,omitempty"`
}
type UpdateMovieArgs struct {
	UpdateGenre  bool
	UpdateArtist bool
}

type FilterMovie struct {
	SearchBy string `json:"searchBy" query:"searchBy"`
	Search   string `json:"search" query:"search"`
	Pagination
}

func (c *CreateMovie) ToMovie() Movie {
	genre := []Genre{}
	for _, v := range c.Genres {
		genre = append(genre, Genre{Model: gorm.Model{ID: v}})
	}
	artist := []Artist{}
	for _, v := range c.Artists {
		artist = append(artist, Artist{Model: gorm.Model{ID: v}})
	}
	return Movie{
		Title:       c.Title,
		Description: c.Description,
		Duration:    c.Duration,
		WatchURL:    c.WatchURL,
		Genres:      genre,
		Artists:     artist,
	}
}

func NewMovieHTTPResponse(m Movie) MovieHTTPResponse {
	genres := []GenreHTTPResponse{}
	for _, v := range m.Genres {
		genres = append(genres, NewGenreHTTPResponse(v))
	}
	artists := []ArtistHTTPResponse{}
	for _, v := range m.Artists {
		artists = append(artists, NewArtistHTTPResponse(v))
	}
	return MovieHTTPResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		WatchURL:    m.WatchURL,
		Genres:      genres,
		Artists:     artists,
	}
}
