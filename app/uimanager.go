package app

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/config"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UIModule struct {
	ID             string     `json:"id"`
	Position       [2]float64 `json:"position"`
	PositionUpdate bool       `json:"position_update"`
	Scale          int        `json:"scale"`
	ScaleUpdate    bool       `json:"scale_update"`
	Visible        bool       `json:"visible"`
	VisibleUpdate  bool       `json:"visible_update"`
}

type UIManager struct {
	Templates        *jet.Set
	Actions          map[string]ManialinkAction
	PublicManialinks map[string]*Manialink
	PlayerManialinks map[string]map[string]*Manialink
	Modules          []UIModule
	ScriptCalls      []string
	Theme            models.Theme
}

var (
	uiInstance *UIManager
	uiOnce     sync.Once
)

func GetUIManager() *UIManager {
	uiOnce.Do(func() {
		uiInstance = &UIManager{
			Actions:          make(map[string]ManialinkAction, 0),
			PublicManialinks: make(map[string]*Manialink, 0),
			PlayerManialinks: make(map[string]map[string]*Manialink, 0),
			Modules:          make([]UIModule, 0),
			ScriptCalls:      make([]string, 0),
			Theme:            models.Theme{},
		}
	})
	return uiInstance
}

func (uim *UIManager) Init() {
	zap.L().Info("Initializing UIManager")

	uim.Templates = jet.NewSet(jet.NewOSFileSystemLoader("./ui/templates"))
	uim.loadTheme()

	// Add global functions to the template engine
	uim.Templates.AddGlobalFunc("floor", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("floor", 1, 1)
		if args.Get(0).Kind() == reflect.Float64 {
			return reflect.ValueOf(math.Floor(args.Get(0).Float()))
		} else if args.Get(0).Kind() == reflect.Int64 {
			return reflect.ValueOf(math.Floor(float64(args.Get(0).Int())))
		}
		return reflect.ValueOf("")
	})

	uim.Templates.AddGlobalFunc("ceil", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("ceil", 1, 1)
		if args.Get(0).Kind() == reflect.Float64 {
			return reflect.ValueOf(math.Ceil(args.Get(0).Float()))
		} else if args.Get(0).Kind() == reflect.Int64 {
			return reflect.ValueOf(math.Floor(float64(args.Get(0).Int())))
		}
		return reflect.ValueOf("")
	})

	uim.Templates.AddGlobalFunc("formatDate", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("formatDate", 1, 1)
		value := args.Get(0)
		if value.Type() == reflect.TypeOf(time.Time{}) {
			extractedTime := value.Interface().(time.Time)
			return reflect.ValueOf(extractedTime.Format("02 January, 15:04"))
		}
		return reflect.ValueOf("")
	})

	uim.Templates.AddGlobalFunc("formatTime", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("formatTime", 1, 1)
		value := args.Get(0)
		time := value.Interface().(int)
		if time < 60000 {
			return reflect.ValueOf(fmt.Sprintf("%.3f", float64(time)/1000))
		}
		return reflect.ValueOf(fmt.Sprintf("%d:%06.3f", time/60000, float64(time%60000)/1000))
	})

	// Add callbacks
	GetClient().OnPlayerManialinkPageAnswer = append(GetClient().OnPlayerManialinkPageAnswer, gbxclient.GbxCallbackStruct[events.PlayerManialinkPageAnswerEventArgs]{
		Key:  "uimanager",
		Call: uim.onManialinkAnswer,
	})

	GetClient().OnPlayerConnect = append(GetClient().OnPlayerConnect, gbxclient.GbxCallbackStruct[events.PlayerConnectEventArgs]{
		Key:  "uimanager",
		Call: uim.onPlayerConnect,
	})

	GetClient().OnPlayerDisconnect = append(GetClient().OnPlayerDisconnect, gbxclient.GbxCallbackStruct[events.PlayerDisconnectEventArgs]{
		Key:  "uimanager",
		Call: uim.onPlayerDisconnect,
	})

	GetClient().AddScriptCallback("Common.UIModules.Properties", "uimanager", uim.onUIModulesProperties)
	uim.getUIProperties()

	zap.L().Info("UIManager initialized")
}

func (uim *UIManager) AfterInit() {
	uim.setUIProperty("Race_RespawnHelper", "Visible", false)
	uim.sendUIProperties()
}

func (uim *UIManager) loadTheme() {
	uim.Templates.AddGlobal("Styling", config.Theme.Styling)
	uim.Templates.AddGlobal("Fonts", config.Theme.Fonts)
	uim.Templates.AddGlobal("Icons", config.Theme.Icons)
	uim.Theme = config.Theme
}

func (uim *UIManager) getUIProperties() {
	uuid := uuid.NewString()
	uim.ScriptCalls = append(uim.ScriptCalls, uuid)
	GetClient().TriggerModeScriptEventArray("Common.UIModules.GetProperties", []string{uuid})
}

func (uim *UIManager) setUIProperty(ID string, property string, value any) {
	for i := range uim.Modules {
		if uim.Modules[i].ID == ID {
			switch property {
			case "Position":
				uim.Modules[i].Position = value.([2]float64)
				uim.Modules[i].PositionUpdate = true
			case "Scale":
				uim.Modules[i].Scale = value.(int)
				uim.Modules[i].ScaleUpdate = true
			case "Visible":
				uim.Modules[i].Visible = value.(bool)
				uim.Modules[i].VisibleUpdate = true
			}

			return
		}
	}

	zap.L().Error("Module not found", zap.String("ID", ID))
}

func (uim *UIManager) sendUIProperties() {
	var moduleProperties struct {
		UIModules []UIModule `json:"uimodules"`
	}
	moduleProperties.UIModules = uim.Modules

	jsonBytes, err := json.Marshal(moduleProperties)
	if err != nil {
		zap.L().Error("Error marshalling module properties", zap.Error(err))
		return
	}

	GetClient().TriggerModeScriptEventArray("Common.UIModules.SetProperties", []string{string(jsonBytes)})
}

func (uim *UIManager) onUIModulesProperties(event any) {
	// Ensure event is a slice
	outerArray, ok := event.([]any)
	if !ok {
		zap.L().Error("Error: event is not a JSON array")
		return
	}

	if len(outerArray) == 0 {
		zap.L().Error("Error: No data found in JSON array")
		return
	}

	// Extract the first element (which is expected to be a JSON string)
	innerJSONString, ok := outerArray[0].(string)
	if !ok {
		zap.L().Error("Error: First element is not a JSON string")
		return
	}

	// Define the target struct
	var moduleProperties struct {
		ResponseID string     `json:"responseid"`
		UIModules  []UIModule `json:"uimodules"`
	}

	// Unmarshal the extracted JSON string
	err := json.Unmarshal([]byte(innerJSONString), &moduleProperties)
	if err != nil {
		zap.L().Error("Error unmarshalling module properties", zap.Error(err))
		return
	}

	uim.ScriptCalls, ok = utils.Remove(uim.ScriptCalls, moduleProperties.ResponseID)
	if !ok {
		return
	}

	uim.Modules = moduleProperties.UIModules

	var reset []string
	for _, module := range uim.Modules {
		reset = append(reset, module.ID)
	}

	var resetRequest struct {
		UIModules []string `json:"uimodules"`
	}
	resetRequest.UIModules = reset

	jsonBytes, err := json.Marshal(resetRequest)
	if err != nil {
		zap.L().Error("Error marshalling reset request", zap.Error(err))
		return
	}

	GetClient().TriggerModeScriptEventArray("Common.UIModules.ResetProperties", []string{string(jsonBytes)})
}

func (uim *UIManager) AddAction(callback func(string, any, any), data any) string {
	uuid := uuid.NewString()
	uim.Actions[uuid] = ManialinkAction{
		Callback: callback,
		Data:     data,
	}

	return uuid
}

func (uim *UIManager) RemoveAction(uuid string) {
	delete(uim.Actions, uuid)
}

func (uim *UIManager) onManialinkAnswer(manialinkAnswerEvent events.PlayerManialinkPageAnswerEventArgs) {
	if action, exists := uim.Actions[manialinkAnswerEvent.Answer]; exists {
		action.Callback(manialinkAnswerEvent.Login, action.Data, manialinkAnswerEvent.Entries)
	}
}

func (uim *UIManager) onPlayerConnect(playerConnectEvent events.PlayerConnectEventArgs) {
	for _, manialink := range uim.PublicManialinks {
		render, err := manialink.Render()
		if err != nil {
			zap.L().Error("Error rendering manialink", zap.Error(err))
			continue
		}
		xml := fmt.Sprintf("<manialinks>%s</manialinks>", render)
		GetClient().SendDisplayManialinkPageToLogin(playerConnectEvent.Login, gbxclient.CData(xml), 0, false)
	}
}

func (uim *UIManager) onPlayerDisconnect(playerDisconnectEvent events.PlayerDisconnectEventArgs) {
	for _, manialink := range uim.PlayerManialinks[playerDisconnectEvent.Login] {
		uim.DestroyManialink(manialink)
	}
}

func (uim *UIManager) sendManialink(ml *Manialink) {
	render, err := ml.Render()
	if err != nil {
		zap.L().Error("Error rendering manialink", zap.Error(err))
		return
	}

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><manialinks>%s</manialinks>`, render)

	if ml.Recipient == nil {
		GetClient().SendDisplayManialinkPage(gbxclient.CData(xml), 0, false)
	} else {
		GetClient().SendDisplayManialinkPageToLogin(*ml.Recipient, gbxclient.CData(xml), 0, false)
	}
}

func (uim *UIManager) DisplayManialink(ml *Manialink) {
	if ml.Recipient == nil {
		uim.PublicManialinks[ml.ID] = ml
	} else {
		if _, ok := uim.PlayerManialinks[*ml.Recipient]; !ok {
			uim.PlayerManialinks[*ml.Recipient] = make(map[string]*Manialink, 0)
		}
		uim.PlayerManialinks[*ml.Recipient][ml.ID] = ml
	}

	uim.sendManialink(ml)
}

func (uim *UIManager) RefreshManialink(ml *Manialink) {
	uim.sendManialink(ml)
}

func (uim *UIManager) HideManialink(ml *Manialink) {
	zap.L().Debug("Hiding manialink", zap.String("id", ml.ID))
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<manialinks><manialink id="%s"></manialink></manialinks>`, ml.ID)

	if ml.Recipient == nil {
		GetClient().SendDisplayManialinkPage(gbxclient.CData(xml), 0, false)
	} else {
		GetClient().SendDisplayManialinkPageToLogin(*ml.Recipient, gbxclient.CData(xml), 0, false)
	}
}

func (uim *UIManager) DestroyManialink(ml *Manialink) {
	zap.L().Debug("Destroying manialink", zap.String("id", ml.ID))

	uim.HideManialink(ml)

	// Remove actions
	for _, value := range ml.Actions {
		uim.RemoveAction(value)
	}

	ml.Data = nil
	if ml.Recipient != nil {
		for key, value := range uim.PlayerManialinks[*ml.Recipient] {
			if value.ID == ml.ID {
				delete(uim.PlayerManialinks[*ml.Recipient], key)
			}
		}
	} else {
		delete(uim.PublicManialinks, ml.ID)
	}
}
