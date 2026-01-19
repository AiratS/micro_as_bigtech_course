package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return nil
}
