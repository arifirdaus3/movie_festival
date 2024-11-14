package model

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresTimeZone string `mapstructure:"POSTGRES_TIMEZONE"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`

	HTTPPort string `mapstructure:"HTTP_PORT"`
}
