package windows

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/ui"
	"go.uber.org/zap"
	"slices"
)

type PlayersListWindow struct {
	*ui.ListWindow
}

func CreatePlayersListWindow(login *string) *PlayersListWindow {
	plw := &PlayersListWindow{
		ListWindow: ui.NewListWindow(login),
	}

	return plw
}

func (plw *PlayersListWindow) OnSpectatorToggle(login string, data any, _ any) {
	player := data.(models.DetailedPlayer)
	if player.IsSpectator {
		return
	}

	c := app.GetGoController()

	if err := c.Server.Client.ForceSpectator(player.Login, 3); err != nil {
		go c.ChatError("Error forcing spectator", err, login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Forced #White#%s #Primary#to spectator", player.NickName), login)
	zap.L().Debug("Forced player to spectator", zap.String("admin", login), zap.String("player", player.Login))

	for _, item := range plw.Items {
		if item[1] == player.Login {
			for i, field := range item {
				if colorMap, ok := field.(models.Toggle); ok {
					colorMap.Color = "Red"
					item[i] = colorMap
				}
			}
		}
	}

	plw.Refresh()
}

func (plw *PlayersListWindow) OnBan(login string, data any, _ any) {
	player := data.(models.DetailedPlayer)
	c := app.GetGoController()

	c.CommandManager.ExecuteCommand(login, "//ban", []string{player.Login}, true)

	// Remove item from list
	for i, item := range plw.Items {
		if item[1] == player.Login {
			plw.Items = slices.Delete(plw.Items, i, i+1)
			break
		}
	}

	plw.Refresh()
}