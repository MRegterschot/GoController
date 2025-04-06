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

	// Server
	ServerLogin string
	ServerPass  string

	// Contact info
	Contact string

	// Delimiter for theme processing
	Delimiter string

	LogLevel string
}

var AppEnv *Env

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	port, err := strconv.Atoi(os.Getenv("XMLRPC_PORT"))
	if err != nil {
		port = 5000
	}

	AppEnv = &Env{
		Host: os.Getenv("XMLRPC_HOST"),
		Port: port,
		User: os.Getenv("XMLRPC_USER"),
		Pass: os.Getenv("XMLRPC_PASS"),

		MasterAdmins: os.Getenv("MASTER_ADMINS"),

		MongoUri: os.Getenv("MONGO_URI"),
		MongoDb:  os.Getenv("MONGO_DB"),

		ServerLogin: os.Getenv("SERVER_LOGIN"),
		ServerPass:  os.Getenv("SERVER_PASS"),
		
		Contact: os.Getenv("CONTACT_INFO"),

		Delimiter: os.Getenv("DELIMITER"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}

	return nil
}
