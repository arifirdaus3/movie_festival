package handlerhttp

import (
	"context"
	"io"
	"moviefestival/model"
	"net/http"
	"os"
	"path/filepath"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type HandlerHTTP struct {
	config        *model.Config
	authUsecase   authUsecase
	genreUsecase  genreUsecase
	artistUsecase artistUsecase
	movieUsecase  movieUsecase
}

type authUsecase interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, user model.User) (model.Token, error)
	Refresh(ctx context.Context, token string) (model.Token, error)
}
type genreUsecase interface {
	InsertGenre(ctx context.Context, genres []model.Genre) error
	GetGenres(ctx context.Context, pagination model.Pagination) ([]model.Genre, error)
}

type artistUsecase interface {
	InsertArtist(ctx context.Context, artists []model.Artist) error
	GetArtists(ctx context.Context, pagination model.Pagination) ([]model.Artist, error)
}
type movieUsecase interface {
	InsertMovie(ctx context.Context, movies model.Movie) error
	GetMovies(ctx context.Context, filter model.FilterMovie) ([]model.Movie, error)
	UpdateMovie(ctx context.Context, updateMovie model.UpdateMovie) error
	ViewedMovie(ctx context.Context, user model.User, movie model.Movie) error
}

func NewHandlerHTTP(config *model.Config, authUsecase authUsecase, artistUsecase artistUsecase, genreUsecase genreUsecase, movieUsecase movieUsecase) *HandlerHTTP {
	return &HandlerHTTP{
		config:        config,
		authUsecase:   authUsecase,
		genreUsecase:  genreUsecase,
		artistUsecase: artistUsecase,
		movieUsecase:  movieUsecase,
	}
}

func (h *HandlerHTTP) InitRoute(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	adminRoute := e.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.config.SignTokenSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &model.CustomClaim{}
		},
	}), func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			user, ok := token.Claims.(*model.CustomClaim)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			if !user.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to access this endpoint")
			}
			return next(c)
		}
	})
	userRoute := e.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.config.SignTokenSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &model.CustomClaim{}
		},
	}))

	e.Static("/public", "public")
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
	e.POST("/auth/refresh", h.Refresh)

	adminRoute.GET("/artist", h.GetArtists)
	adminRoute.POST("/artist", h.InsertArtist)

	adminRoute.GET("/genre", h.GetGenres)
	adminRoute.POST("/genre", h.InsertGenre)

	adminRoute.POST("/movie", h.InsertMovie)
	adminRoute.POST("/movie/upload", h.UploadMovie)
	adminRoute.PUT("/movie", h.UpdateMovie)

	e.GET("/movie", h.GetMovies)

	userRoute.POST("/movie/viewed", h.ViewedMovie)

	e.POST("/movie/vote", h.Refresh)
	e.GET("/movie/voted", h.Refresh)
	e.GET("/movie/most-voted", h.Refresh)
	e.GET("/movie/most-viewed", h.Refresh)
}

func (h *HandlerHTTP) Register(c echo.Context) error {
	var user model.User
	c.Bind(&user)
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 30)),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = h.authUsecase.Register(c.Request().Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}

func (h *HandlerHTTP) Login(c echo.Context) error {
	var user model.User
	email, password, ok := c.Request().BasicAuth()
	if !ok {
		c.Bind(&user)
	} else {
		user.Email = email
		user.Password = password
	}

	err := validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 30)),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := h.authUsecase.Login(c.Request().Context(), user)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = token.RefreshToken
	cookie.Expires = token.RefreshTokenExpiration
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token.AcceessToken,
	},
	)
}

func (h *HandlerHTTP) Refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}
	err = validation.Validate(&refreshToken.Value, validation.Required.Error("Invalid refresh token"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := h.authUsecase.Refresh(c.Request().Context(), refreshToken.Value)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token.AcceessToken,
	},
	)
}

func (h *HandlerHTTP) InsertGenre(c echo.Context) error {
	var genre []model.Genre
	c.Bind(&genre)
	for _, v := range genre {
		err := validation.ValidateStruct(&v,
			validation.Field(&v.Name, validation.Required, is.Alphanumeric),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	err := h.genreUsecase.InsertGenre(c.Request().Context(), genre)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
func (h *HandlerHTTP) GetGenres(c echo.Context) error {
	var pagination model.Pagination
	c.Bind(&pagination)
	pagination.Default()
	genres, err := h.genreUsecase.GetGenres(c.Request().Context(), pagination)
	if err != nil {
		return err
	}
	res := []model.GenreHTTPResponse{}
	for _, v := range genres {
		res = append(res, model.NewGenreHTTPResponse(v))
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data:   res,
	})
}

func (h *HandlerHTTP) InsertArtist(c echo.Context) error {
	var artist []model.Artist
	c.Bind(&artist)
	for _, v := range artist {
		err := validation.ValidateStruct(&v,
			validation.Field(&v.Name, validation.Required, is.Alphanumeric),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	err := h.artistUsecase.InsertArtist(c.Request().Context(), artist)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
func (h *HandlerHTTP) GetArtists(c echo.Context) error {
	var pagination model.Pagination
	c.Bind(&pagination)
	pagination.Default()
	artists, err := h.artistUsecase.GetArtists(c.Request().Context(), pagination)
	if err != nil {
		return err
	}
	res := []model.ArtistHTTPResponse{}
	for _, v := range artists {
		res = append(res, model.NewArtistHTTPResponse(v))
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data:   res,
	})
}

func (h *HandlerHTTP) InsertMovie(c echo.Context) error {
	var movie model.CreateMovie
	c.Bind(&movie)
	err := validation.ValidateStruct(&movie,
		validation.Field(&movie.Duration, validation.Required),
		validation.Field(&movie.Title, validation.Required),
		validation.Field(&movie.Genres, validation.Required, validation.Length(1, 99)),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.movieUsecase.InsertMovie(c.Request().Context(), movie.ToMovie())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
func (h *HandlerHTTP) UpdateMovie(c echo.Context) error {
	var movie model.UpdateMovie
	c.Bind(&movie)

	err := h.movieUsecase.UpdateMovie(c.Request().Context(), movie)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
func (h *HandlerHTTP) UploadMovie(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide video file")
	}
	ext := filepath.Ext(file.Filename)
	if !model.ValidVideoExt[ext] {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file extension")
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	newFileName := "/public/" + uuid.New().String() + ext
	// Destination
	dst, err := os.Create("." + newFileName)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data: map[string]interface{}{
			"url": newFileName,
		},
	})
}

func (h *HandlerHTTP) GetMovies(c echo.Context) error {
	var filter model.FilterMovie
	c.Bind(&filter)
	filter.Default()
	movies, err := h.movieUsecase.GetMovies(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	res := []model.MovieHTTPResponse{}
	for _, v := range movies {
		res = append(res, model.NewMovieHTTPResponse(v))
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data:   res,
	})
}

func (h *HandlerHTTP) ViewedMovie(c echo.Context) error {
	var movie model.MovieHTTPResponse
	c.Bind(&movie)
	err := validation.ValidateStruct(&movie,
		validation.Field(&movie.ID, validation.Required),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	user, ok := token.Claims.(*model.CustomClaim)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	err = h.movieUsecase.ViewedMovie(c.Request().Context(), model.User{
		Model: gorm.Model{ID: user.ID},
	}, model.Movie{
		Model: gorm.Model{ID: movie.ID},
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
