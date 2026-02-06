package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBDsn   string
}

func Load() (*Config, error) {
	
	if err := godotenv.Load("config.env"); err != nil {
		return nil, err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	return &Config{
		AppPort: os.Getenv("APP_PORT"),
		DBDsn:   dsn,
	}, nil
}

