package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found bhai, using environment variables")
	}
}

func Env(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func JWTSecret() string {
	secret := Env("JWT_SECRET", "my_secret_key")
	return secret
}
