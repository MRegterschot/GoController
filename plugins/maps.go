package plugins

import "fmt"

type MapsPlugin struct {
	Name string
	Dependencies []string
	Loaded bool
}

func CreateMapsPlugin() *MapsPlugin {
	return &MapsPlugin{
		Name: "Maps",
		Dependencies: []string{},
		Loaded: false,
	}
}

func (m *MapsPlugin) Load() error {
	fmt.Println("Loading Maps plugin")
	return nil
}

func (m *MapsPlugin) Unload() error {
	fmt.Println("Unloading Maps plugin")
	return nil
}

func init() {
	mapsPlugin := CreateMapsPlugin()
	GetPluginManager().PreLoadPlugin(mapsPlugin)
}