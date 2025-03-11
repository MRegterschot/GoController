package app

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type PluginManager struct {
	PreLoadedPlugins []any
	Plugins          []any
}

type BasePlugin struct {
	CommandManager  *CommandManager
	SettingsManager *SettingsManager
	GoController    *GoController
}

func GetBasePlugin() BasePlugin {
	return BasePlugin{
		CommandManager:  GetCommandManager(),
		SettingsManager: GetSettingsManager(),
		GoController:    GetGoController(),
	}
}

var (
	pmInstance *PluginManager
	pmOnce     sync.Once
)

func GetPluginManager() *PluginManager {
	pmOnce.Do(func() {
		pmInstance = &PluginManager{
			Plugins: []any{},
		}
	})
	return pmInstance
}

func (pm *PluginManager) Init() {
	zap.L().Info("Initializing PluginManager")
	pm.RegisterPlugins()
	pm.LoadPlugins()
	zap.L().Info("PluginManager initialized")
}

func (pm *PluginManager) PreLoadPlugin(plugin any) {
	pm.PreLoadedPlugins = append(pm.PreLoadedPlugins, plugin)
}

func (pm *PluginManager) RegisterPlugins() {
	for _, plugin := range pm.PreLoadedPlugins {
		name, ok := isPlugin(plugin)

		if !ok {
			zap.L().Error("Failed to register plugin", zap.String("plugin", reflect.TypeOf(plugin).Name()), zap.String("reason", "plugin is not a valid plugin"))
		}
		zap.L().Info("Registering plugin", zap.String("plugin", name))
		pm.Plugins = append(pm.Plugins, plugin)
	}
}

func (pm *PluginManager) LoadPlugins() {
	for _, plugin := range pm.Plugins {
		p := reflect.ValueOf(plugin).Elem()
		loaded := p.FieldByName("Loaded").Bool()

		if loaded {
			continue
		}

		pl, _ := plugin.(interface{ Load() error })
		err := pl.Load()
		if err != nil {
			zap.L().Error("Failed to load plugin", zap.String("plugin", p.FieldByName("Name").String()), zap.Error(err))
		}
		p.FieldByName("Loaded").SetBool(true)
		zap.L().Info("Loaded plugin", zap.String("plugin", p.FieldByName("Name").String()))
	}
}

func (pm *PluginManager) LoadPlugin(name string) error {
	for _, plugin := range pm.Plugins {
		p := reflect.ValueOf(plugin).Elem()
		if strings.ToLower(p.FieldByName("Name").String()) == name {
			pl, _ := plugin.(interface{ Load() error })
			err := pl.Load()
			if err != nil {
				zap.L().Error("Failed to load plugin", zap.String("plugin", p.FieldByName("Name").String()), zap.Error(err))
			}
			p.FieldByName("Loaded").SetBool(true)
			zap.L().Info("Loaded plugin", zap.String("plugin", p.FieldByName("Name").String()))
			return nil
		}
	}

	return errors.New("plugin not found")
}

func (pm *PluginManager) UnloadPlugin(name string) error {
	for _, plugin := range pm.Plugins {
		p := reflect.ValueOf(plugin).Elem()
		if strings.ToLower(p.FieldByName("Name").String()) == name {
			pl, _ := plugin.(interface{ Unload() error })
			err := pl.Unload()
			if err != nil {
				zap.L().Error("Failed to unload plugin", zap.String("plugin", p.FieldByName("Name").String()), zap.Error(err))
			}
			p.FieldByName("Loaded").SetBool(false)
			zap.L().Info("Unloaded plugin", zap.String("plugin", p.FieldByName("Name").String()))
			return nil
		}
	}

	return errors.New("plugin not found")
}

func isPlugin(plugin any) (string, bool) {
	t := reflect.TypeOf(plugin)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", false
	}

	_, hasName := t.FieldByName("Name")
	_, hasDependencies := t.FieldByName("Dependencies")
	_, hasLoaded := t.FieldByName("Loaded")
	_, hasLoad := plugin.(interface{ Load() error })
	_, hasUnload := plugin.(interface{ Unload() error })

	if !hasName {
		return "", false
	}

	name := reflect.ValueOf(plugin).Elem().FieldByName("Name").String()

	return name, hasName && hasDependencies && hasLoaded && hasLoad && hasUnload
}
