package plugins

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/windows"
	"github.com/MRegterschot/GoController/ui"
	"github.com/MRegterschot/GoController/utils"
	"go.uber.org/zap"
)

type MapsPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateMapsPlugin() *MapsPlugin {
	return &MapsPlugin{
		Name:         "Maps",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *MapsPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "/maps",
		Callback: p.mapsCommand,
		Admin:    false,
		Help:     "Shows all available maps",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//localmaps",
		Callback: p.localMapsCommand,
		Admin:    true,
		Help:     "Manage local maps",
	})

	return nil
}

func (p *MapsPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("/maps")
	commandManager.RemoveCommand("//localmaps")

	return nil
}

func (p *MapsPlugin) mapsCommand(login string, args []string) {
	c := app.GetGoController()

	window := windows.CreateMapsGridWindow(&login)
	window.Title = "Maps"
	window.Items = make([]any, 0, len(c.MapManager.Maps))

	isAdmin := c.IsAdmin(login)
	window.IsAdmin = &isAdmin

	for _, m := range c.MapManager.Maps {
		if isAdmin {
			window.Actions["remove_"+m.UId] = app.GetUIManager().AddAction(window.HandleRemoveAnswer, m)
		}
		window.Actions["queue_"+m.UId] = app.GetUIManager().AddAction(window.HandleQueueAnswer, m)
		window.Items = append(window.Items, m)
	}

	go window.Display()
}

func (p *MapsPlugin) localMapsCommand(login string, args []string) {
	c := app.GetGoController()

	mapsPath := app.GetMapManager().MapsPath
	if mapsPath == "" {
		go c.ChatError("No maps directory found", nil, login)
		return
	}

	items := make([][]any, 0)

	err := filepath.WalkDir(mapsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			zap.L().Error("Failed to walk directory", zap.Error(err))
			return err
		}

		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".gbx") {
			return nil
		}

		parentDir := strings.TrimPrefix(filepath.Dir(path), mapsPath)
		name := utils.MapFileRegex.ReplaceAllString(filepath.Base(path), "")

		items = append(items, []any{
			parentDir,
			name,
			app.GetUIManager().AddAction(p.onAddMap, path),
		})
		
		return nil
	})
	if err != nil {
		go c.ChatError("Error walking directory", err, login)
		return
	}
	
	columns := []ui.Column{
		{Name: "Folder", Width: 40},
		{Name: "File Name", Width: 50},
		{Name: "Add", Width: 10, Type: "button"},
	}

	window := ui.NewListWindow(&login)
	window.Title = "Local Maps"
	window.Columns = columns
	window.Items = items

	go window.Display()
}

func (p *MapsPlugin) onAddMap(login string, data any, _ any) {
	file := data.(string)

	c := app.GetGoController()
	if err := c.Server.Client.AddMap(file); err != nil {
		zap.L().Error("Failed to add map", zap.String("file", file), zap.Error(err))
		go c.ChatError("Error adding map", err, login)
		return
	}

	c.MapManager.SyncMaps()
	go c.Chat("#Primary#Map added successfully", login)
	zap.L().Info("Map added successfully", zap.String("file", file))
}

func init() {
	mapsPlugin := CreateMapsPlugin()
	app.GetPluginManager().PreLoadPlugin(mapsPlugin)
}
