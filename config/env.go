package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	// XMLRPC server configuration
	Host string
	Port int
	User string
	Pass string

	// Master admins
	MasterAdmins string

	// Database configuration
	MongoUri string
	MongoDb  string

	// Delimiter for theme processing
	Delimiter string

	LogLevel string
}

// Global variable to store the loaded environment configuration
var AppEnv *Env

// LoadEnv loads the .env file into the environment and initializes AppEnv
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
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

		MasterAdmins: os.Getenv("MASTER_ADMINS"),

		MongoUri: os.Getenv("MONGO_URI"),
		MongoDb:  os.Getenv("MONGO_DB"),

		Delimiter: os.Getenv("DELIMITER"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}

	return nil
}
