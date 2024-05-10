package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	APP_PORT string
	APP_NAME string
	APP_ENV  string

	JWT_SECRET         string
	JWT_EXPIRY         int
	JWT_REFRESH_SECRET string
	JWT_REFRESH_EXPIRY int

	DB_HOST     string
	DB_PORT     string
	DB_DATABASE string
	DB_USERNAME string
	DB_PASSWORD string
}

var Get *config

func Load(path string) {
	godotenv.Load(path)
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		panic("INVALID JWT_EXPIRY, check .env")
	}
	refreshExp, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		panic("INVALID JWT_REFRESH_EXPIRY, check .env")
	}
	Get = &config{
		APP_PORT: os.Getenv("APP_PORT"),
		APP_NAME: os.Getenv("APP_NAME"),
		APP_ENV:  os.Getenv("APP_ENV"),

		JWT_SECRET:         os.Getenv("JWT_SECRET"),
		JWT_EXPIRY:         exp,
		JWT_REFRESH_SECRET: os.Getenv("JWT_REFRESH_SECRET"),
		JWT_REFRESH_EXPIRY: refreshExp,

		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_DATABASE: os.Getenv("DB_DATABASE"),
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	}
}
