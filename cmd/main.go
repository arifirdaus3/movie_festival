package main

import (
	"fmt"
	"log"
	handlerhttp "moviefestival/handler/http"
	"moviefestival/model"
	repositorysql "moviefestival/repository/sql"
	artistusecase "moviefestival/usecase/artist"
	authusecase "moviefestival/usecase/auth"
	genreusecase "moviefestival/usecase/genre"
	movieusecase "moviefestival/usecase/movie"
	voteusecase "moviefestival/usecase/vote"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := initConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDB, config.PostgresPort, config.PostgresTimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open database, ", err)
	}
	err = db.AutoMigrate(&model.Artist{}, &model.User{}, &model.Genre{}, &model.Movie{}, &model.UserMovieVote{})
	if err != nil {
		log.Fatal("failed to migrate database model, ", err)
	}
	e := echo.New()

	repositorySQL := repositorysql.NewRepositorySQL(db)

	authUsecase := authusecase.NewAuthUsecase(&config, repositorySQL)
	artistUsecase := artistusecase.NewArtistUsecase(repositorySQL)
	genreUsecase := genreusecase.NewGenreUsecase(repositorySQL)
	movieUsecase := movieusecase.NewMovieUsecase(repositorySQL)
	voteUsecase := voteusecase.NewVoteUsecase(repositorySQL)

	handlerHTTP := handlerhttp.NewHandlerHTTP(&config, authUsecase, artistUsecase, genreUsecase, movieUsecase, voteUsecase)

	handlerHTTP.InitRoute(e)

	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}

func initConfig() model.Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	errReadConfig := viper.ReadInConfig()
	var config model.Config
	err := viper.Unmarshal(&config)
	if err != nil || config.PostgresHost == "" {
		log.Fatal("failed to unmarshal config, ", err, errReadConfig)
	}
	return config
}
