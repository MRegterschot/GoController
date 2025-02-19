package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/MRegterschot/GoController/config"
	"go.uber.org/zap"
)

type SettingsManager struct {
	Settings     map[string]string
	Admins       []string
	MasterAdmins []string

	AdminsFile string
}

// Creates a new SettingsManager
func NewSettingsManager() *SettingsManager {
	adminPath := "./userdata/admins.json"

	masterAdmins := strings.Split(config.AppEnv.MasterAdmins, ",")
	for i, admin := range masterAdmins {
		masterAdmins[i] = strings.TrimSpace(admin)
	}

	sm := &SettingsManager{
		Settings:   make(map[string]string),
		MasterAdmins: masterAdmins,
		AdminsFile: adminPath,
	}

	err := sm.CreateFile(adminPath, []string{})
	if err != nil {
		zap.L().Fatal("Failed to create admins file", zap.Error(err))
	}

	return sm
}

func (sm *SettingsManager) Init() {
	zap.L().Info("Initializing SettingsManager")
	sm.LoadAdmins()
	zap.L().Info("SettingsManager initialized")
}

func (sm *SettingsManager) LoadAdmins() {
	f, err := os.Open(sm.AdminsFile)
	if err != nil {
		zap.L().Fatal("Failed to open admins file", zap.Error(err))
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&sm.Admins)
	if err != nil {
		zap.L().Fatal("Failed to decode admins file", zap.Error(err))
	}

	sm.Admins = append(sm.Admins, sm.MasterAdmins...)
}

// Create file with data if it doesn't exist
func (sm *SettingsManager) CreateFile(file string, data interface{}) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return err // Return error instead of panicking
		}
		defer f.Close()

		// Format data into JSON and write it to the file
		err = json.NewEncoder(f).Encode(data)
		if err != nil {
			return err // Return error instead of panicking
		}
	}
	return nil
}
