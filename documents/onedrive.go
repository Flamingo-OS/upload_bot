package documents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
)

type tokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func getAccessToken() (string, error) {
	apiUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", core.Config.OneDriveTenantId)
	jsonReqData := map[string]string{
		"client_id":     core.Config.OneDriveClientId,
		"client_secret": core.Config.OneDriveClientSecret,
		"refresh_token": core.Config.OneDriveRefreshToken,
		"grant_type":    "refresh_token",
		"scope":         "offline_access Files.ReadWrite.All",
		"redirect_uri":  "http://localhost:8080",
	}
	reqData, err := json.Marshal(jsonReqData)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(string(reqData)))
	if err != nil {
		core.Log.Fatal("Error creating request: ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		core.Log.Fatal("Error sending request: ", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) == "" {
		return "", fmt.Errorf("empty response body")
	}

	var tokResp tokenResponse
	_ = json.Unmarshal(body, &tokResp)
	return tokResp.AccessToken, nil
}

func listDir(accessToken string, fileId string) (map[string]string, error) {
	apiUrl := "https://graph.microsoft.com/v1.0/me/drive/"
	if fileId != "" {
		apiUrl += "items/" + fileId + "/children"
	} else {
		apiUrl += "root/children"
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		core.Log.Fatal("Error creating request: ", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		core.Log.Fatal("Error sending request: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) == "" {
		return nil, fmt.Errorf("empty response body")
	}

	var data interface{}
	json.Unmarshal(body, &data)

	// huge hax to prevent me from actually fully defining the response.
	// I have no idea what the json is and need a lot of help correcting in Go.
	// If you know how to do this, please help me.
	var files map[string]string = make(map[string]string)
	for _, dir := range data.(map[string]interface{})["value"].([]interface{}) {
		fileId := dir.(map[string]interface{})["id"].(string)
		fileName := dir.(map[string]interface{})["name"].(string)
		files[fileId] = fileName
	}
	return files, nil
}
