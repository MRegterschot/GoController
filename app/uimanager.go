package app

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UIModule struct {
	ID       string     `json:"id"`
	Position [2]float64 `json:"position"`
	Scale    int        `json:"scale"`
	Visible  bool       `json:"visible"`
}

type UIManager struct {
	Templates        *jet.Set
	Actions          map[string]models.ManialinkAction
	PublicManiaLinks map[string]*models.Manialink
	PlayerManiaLinks map[string][]*models.Manialink
	Modules          []UIModule
	ScriptCalls      []string
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

	uim.Templates = jet.NewSet(jet.NewOSFileSystemLoader("./templates"))

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

func (uim *UIManager) DisplayManialink(ml *models.Manialink) {
	zap.L().Info("Displaying manialink", zap.String("id", ml.ID))
}

func (uim *UIManager) RefreshManialink(ml *models.Manialink) {
	zap.L().Info("Refreshing manialink", zap.String("id", ml.ID))
}

func (uim *UIManager) getUIProperties() {
	uuid := uuid.NewString()
	uim.ScriptCalls = append(uim.ScriptCalls, uuid)
	GetClient().TriggerModeScriptEventArray("Common.UIModules.GetProperties", []string{uuid})
}

func (uim *UIManager) onManialinkAnswer(manialinkAnswerEvent events.PlayerManialinkPageAnswerEventArgs) {
	fmt.Println(manialinkAnswerEvent)
}

func (uim *UIManager) onPlayerConnect(playerConnectEvent events.PlayerConnectEventArgs) {
	fmt.Println(playerConnectEvent)
}

func (uim *UIManager) onPlayerDisconnect(playerDisconnectEvent events.PlayerDisconnectEventArgs) {
	fmt.Println(playerDisconnectEvent)
}

func (uim *UIManager) onUIModulesProperties(event interface{}) {
	// Ensure event is a slice
	outerArray, ok := event.([]interface{})
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
