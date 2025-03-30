package windows

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/ui"
	"go.uber.org/zap"
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

	go c.Chat(fmt.Sprintf("#Primary#Forced %s to spectator", player.NickName), login)
	zap.L().Debug("Forced player to spectator", zap.String("login", player.Login), zap.String("nickname", player.NickName))

	for _, item := range plw.Items {
		if item[1] == player.Login {
			if colorMap, ok := item[3].(models.Toggle); ok {
				colorMap.Color = "Red"
				item[3] = colorMap
			}
		}
	}

	plw.Refresh()
}
