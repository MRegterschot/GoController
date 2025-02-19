package config

import (
	"os"

	"go.uber.org/zap"
)

func Setup() {
	// load env
	err := LoadEnv()
	if err != nil {
		os.Exit(1)
	}

	// setup logger
	SetupLogger()

	defer zap.L().Sync()
}