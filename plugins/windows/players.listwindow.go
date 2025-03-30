package windows

import (
	"fmt"

	"slices"

	"github.com/MRegterschot/GbxRemoteGo/structs"
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

	go c.Chat(fmt.Sprintf("#Primary#Forced #White#%s #Primary#to spectator", player.NickName), login)
	zap.L().Debug("Forced player to spectator", zap.String("admin", login), zap.String("player", player.Login))

	for _, item := range plw.Items {
		if item[1] == player.Login {
			if spec, ok := item[2].(models.Toggle); ok {
				spec.Color = "Red"
				item[2] = spec
			}
		}
	}

	plw.Refresh()
}

func (plw *PlayersListWindow) OnKick(login string, data any, _ any) {
	player := data.(models.DetailedPlayer)
	c := app.GetGoController()

	c.CommandManager.ExecuteCommand(login, "//kick", []string{player.Login}, true)

	// Remove item from list
	for i, item := range plw.Items {
		if item[1] == player.Login {
			plw.Items = slices.Delete(plw.Items, i, i+1)
			break
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

func (plw *PlayersListWindow) OnGuestToggle(login string, data any, _ any) {
	player := data.(models.DetailedPlayer)

	c := app.GetGoController()

	guestList := make([]structs.TMGuestListEntry, 0)
	loop := 0
	for true {
		gl, err := c.Server.Client.GetGuestList(100, loop)
		if err != nil {
			go c.ChatError("Error getting guest list", err, login)
			return
		}

		guestList = append(guestList, gl...)
		if len(gl) < 100 {
			break
		}
		loop++
	}

	isGuest := false
	for _, gl := range guestList {
		if gl.Login == player.Login {
			isGuest = true
			break
		}
	}

	if isGuest {
		if err := c.Server.Client.RemoveGuest(player.Login); err != nil {
			go c.ChatError("Error removing guest", err, login)
			return
		}
		c.Chat(fmt.Sprintf("#Primary#Removed #White#%s #Primary#from guest list", player.NickName), login)
		zap.L().Debug("Removed player from guest list", zap.String("admin", login), zap.String("player", player.Login))
	} else {
		if err := c.Server.Client.AddGuest(player.Login); err != nil {
			go c.ChatError("Error adding guest", err, login)
			return
		}
		c.Chat(fmt.Sprintf("#Primary#Added #White#%s #Primary#to guest list", player.NickName), login)
		zap.L().Debug("Added player to guest list", zap.String("admin", login), zap.String("player", player.Login))
	}

	// Update item in list
	for _, item := range plw.Items {
		if item[1] == player.Login {
			if guest, ok := item[5].(models.Toggle); ok {
				if isGuest {
					guest.Color = "Red"
				} else {
					guest.Color = "Green"
				}
				item[5] = guest
			}
		}
	}

	plw.Refresh()
}
