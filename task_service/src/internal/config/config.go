package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_Host     string
	DB_Port     string
	DB_User     string
	DB_Password string
	DB_Name     string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load("../.secrets/.env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Return Config struct populated with env vars
	return &Config{
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_User:     os.Getenv("DB_USER"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME"),
	}
}
