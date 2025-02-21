package plugins

import (
	"reflect"
	"sync"

	"go.uber.org/zap"
)

type PluginManager struct {
	PreLoadedPlugins []interface{}
	Plugins []interface{}
}

var (
	instance *PluginManager
	once     sync.Once
)

func GetPluginManager() *PluginManager {
	once.Do(func() {
		instance = &PluginManager{
			Plugins: []interface{}{},
		}
	})
	return instance
}

func (pm *PluginManager) Init() {
	zap.L().Info("Initializing PluginManager")
	for _, plugin := range pm.PreLoadedPlugins {
		name, ok := isPlugin(plugin)

		if !ok {
			zap.L().Error("Failed to register plugin", zap.String("plugin", reflect.TypeOf(plugin).Name()), zap.String("reason", "plugin is not a valid plugin"))
		}
		zap.L().Info("Registering plugin", zap.String("plugin", name))
		pm.Plugins = append(pm.Plugins, plugin)
	}
	zap.L().Info("PluginManager initialized")
}

func (pm *PluginManager) PreLoadPlugin(plugin interface{}) {
	pm.PreLoadedPlugins = append(pm.PreLoadedPlugins, plugin)
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
			zap.L().Fatal("Failed to load plugin", zap.String("plugin", p.FieldByName("Name").String()), zap.Error(err))
		}
		p.FieldByName("Loaded").SetBool(true)
	}
}

func isPlugin(plugin interface{}) (string, bool) {
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