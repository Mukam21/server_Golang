package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	APIAgifyURL       string
	APIGenderizeURL   string
	APINationalizeURL string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort:        os.Getenv("SERVER_PORT"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		APIAgifyURL:       os.Getenv("API_AGIFY_URL"),
		APIGenderizeURL:   os.Getenv("API_GENDERIZE_URL"),
		APINationalizeURL: os.Getenv("API_NATIONALIZE_URL"),
	}, nil
}
