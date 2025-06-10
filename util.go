package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Loads the environment variables from the `.env` file.
func loadEnvVars() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("loading .env file failed: %v", err)
	}

	postgresUri = os.Getenv("POSTGRES_URI")
	if postgresUri == "" {
		return errors.New("POSTGRES_URI not set, please set it in your .env file")
	}

	staticDir = os.Getenv("STATIC_DIR")
	if staticDir == "" {
		return errors.New("STATIC_DIR not set, please set it in your .env file")
	}

	templateDir = os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		return errors.New("TEMPLATE_DIR not set, please set it in your .env file")
	}

	return nil
}
