package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/MRegterschot/GoController/config"
	"github.com/MRegterschot/GoController/models"
	"go.uber.org/zap"
)

type NadeoTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type NadeoAPI struct {
	ProdUrl string
	LiveUrl string
	MeetUrl string

	Tokens NadeoTokens
}

func newNadeoAPI() *NadeoAPI {
	n := &NadeoAPI{
		ProdUrl: "https://prod.trackmania.core.nadeo.online",
		LiveUrl: "https://live-services.trackmania.nadeo.live",
		MeetUrl: "https://meet.trackmania.nadeo.club",
	}

	n.loginDedicated("NadeoServices")

	return n
}

var (
	nadeoAPIInstance *NadeoAPI
	nadeoAPIOnce     sync.Once
)

func GetNadeoAPI() *NadeoAPI {
	nadeoAPIOnce.Do(func() {
		nadeoAPIInstance = newNadeoAPI()
	})
	return nadeoAPIInstance
}

func (api *NadeoAPI) loginDedicated(audience string) {
	login := config.AppEnv.ServerLogin
	pass := config.AppEnv.ServerPass
	contact := config.AppEnv.Contact

	auth := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
	authUrl := api.ProdUrl + "/v2/authentication/token/basic"

	body := map[string]string{"audience": audience}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		zap.L().Error("Failed to marshal JSON", zap.Error(err))
		return
	}
	
	req, err := http.NewRequest("POST", authUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		zap.L().Error("Failed to create request", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", contact)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("Failed to send request", zap.Error(err))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		zap.L().Error("Failed to authenticate", zap.String("status", res.Status))
		return
	}
	
	var tokens NadeoTokens
	err = json.NewDecoder(res.Body).Decode(&tokens)
	if err != nil {
		zap.L().Error("Failed to decode response", zap.Error(err))
		return
	}
	
	api.Tokens = tokens
	zap.L().Info("Authenticated with Nadeo")
}

func (api *NadeoAPI) GetMapsInfo(mapUids []string) ([]models.MapInfo, error) {
	url := api.ProdUrl + "/maps/?mapUidList=" + strings.Join(mapUids, ",")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	res, err := api.doRequest(req)
	if err != nil {
		zap.L().Error("Failed to send request", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		zap.L().Error("Failed to get maps info", zap.String("status", res.Status))
		return nil, err
	}

	var mapsInfo []models.MapInfo
	err = json.NewDecoder(res.Body).Decode(&mapsInfo)
	if err != nil {
		zap.L().Error("Failed to decode response", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully retrieved maps info", zap.Int("count", len(mapsInfo)))
	return mapsInfo, nil
}

// doRequest sends a request to the Nadeo API and returns the response.
// It handles authentication and token refresh if necessary.
// This function should be called for every API request.
func (api *NadeoAPI) doRequest(req *http.Request) (*http.Response, error) {
	if api.Tokens.AccessToken == "" {
		api.loginDedicated("NadeoServices")
	}

	req.Header.Set("Authorization", "nadeo_v1 t="+api.Tokens.AccessToken)
	req.Header.Set("User-Agent", config.AppEnv.Contact)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		api.loginDedicated("NadeoServices")
		req.Header.Set("Authorization", "nadeo_v1 t="+api.Tokens.AccessToken)
		res, err = client.Do(req)
	}

	return res, err
}