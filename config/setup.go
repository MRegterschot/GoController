package config

import (
	"os"
)

func Setup() {
	// load env
	err := LoadEnv()
	if err != nil {
		os.Exit(1)
	}

	// setup logger
	SetupLogger()
	
	// load theme
	LoadTheme()
}