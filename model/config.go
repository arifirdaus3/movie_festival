package model

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresTimeZone string `mapstructure:"POSTGRES_TIMEZONE"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`

	HTTPPort string `mapstructure:"HTTP_PORT"`

	AccessTokenExpirationMinute  int64  `mapstructure:"ACCESS_TOKEN_EXPIRATION_MINUTE"`
	RefreshTokenExpirationMinute int64  `mapstructure:"REFRESH_TOKEN_EXPIRATION_MINUTE"`
	SignTokenSecret              string `mapstructure:"SIGN_TOKEN_SECRET"`
}
