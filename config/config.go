package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type dbConfig struct {
	DB_URI        string
	DB_NAME       string
}

func GetDBConfig() dbConfig {
	return dbConfig{
		DB_URI:        os.Getenv("DB_URI_CLOUD"),
		DB_NAME:       os.Getenv("DB_NAME"),
	}
}
