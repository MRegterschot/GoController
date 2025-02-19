package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Host string
	Port int
	User string
	Pass string
}

// Global variable to store the loaded environment configuration
var AppEnv *Env

// LoadEnv loads the .env file into the environment and initializes AppEnv
func LoadEnv() error {
	goEnv := os.Getenv("GO_ENV")

	// Load .env file only in development mode
	if goEnv == "" || goEnv == "development" {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	// Convert PORT from string to int
	port, err := strconv.Atoi(os.Getenv("XMLRPC_PORT"))
	if err != nil {
		port = 5000 // Default port if conversion fails
	}

	// Initialize global AppEnv variable
	AppEnv = &Env{
		Host: os.Getenv("XMLRPC_HOST"),
		Port: port,
		User: os.Getenv("XMLRPC_USER"),
		Pass: os.Getenv("XMLRPC_PASS"),
	}

	return nil
}
