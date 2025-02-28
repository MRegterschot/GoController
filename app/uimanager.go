package app

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
	"go.uber.org/zap"
)

type UIManager struct {
	Views *jet.Set
}

var (
	uiInstance *UIManager
	uiOnce     sync.Once
)

func GetUIManager() *UIManager {
	uiOnce.Do(func() {
		uiInstance = &UIManager{}
	})
	return uiInstance
}

func (uim *UIManager) Init() {
	zap.L().Info("Initializing UIManager")

	uim.Views = jet.NewSet(jet.NewOSFileSystemLoader("./ui"))

	zap.L().Info("UIManager initialized")
}