package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

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
	core.Log.Info("Getting access token")
	apiUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", core.Config.OneDriveTenantId)
	reqData := url.Values{}
	reqData.Set("client_id", core.Config.OneDriveClientId)
	reqData.Set("scope", "offline_access files.readwrite")
	reqData.Set("client_secret", core.Config.OneDriveClientSecret)
	reqData.Set("grant_type", "refresh_token")
	reqData.Set("redirect_uri", "http://localhost:8080")
	reqData.Set("refresh_token", core.Config.OneDriveRefreshToken)
	client := &http.Client{}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(reqData.Encode()))
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

	var tokResp tokenResponse
	_ = json.Unmarshal(body, &tokResp)
	core.Log.Info("Access token was successfully retrieved")
	return tokResp.AccessToken, nil
}

func listDir(accessToken string, fileId string) (map[string]string, error) {
	core.Log.Info("Listing directory")
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

	if reflect.TypeOf(data.(map[string]interface{})["value"]) != reflect.TypeOf([]interface{}{}) {
		return nil, fmt.Errorf("invalid response body")
	}

	// huge hax to prevent me from actually fully defining the response.
	// I have no idea what the json is and need a lot of help correcting in Go.
	// If you know how to do this, please help me.
	var files map[string]string = make(map[string]string)
	for _, dir := range data.(map[string]interface{})["value"].([]interface{}) {
		fileId := dir.(map[string]interface{})["id"].(string)
		fileName := dir.(map[string]interface{})["name"].(string)
		files[fileId] = fileName
	}
	core.Log.Info("Directory was successfully listed")
	return files, nil
}

func makeFolder(accessToken string, fileName string, parentFileId string) (string, error) {
	core.Log.Info("Making folder")
	apiUrl := "https://graph.microsoft.com/v1.0/me/drive/"
	if parentFileId != "" {
		apiUrl += "items/" + parentFileId + "/children"
	} else {
		apiUrl += "root/children"
	}

	// check if file already exists
	files, err := listDir(accessToken, parentFileId)
	if err != nil {
		return "", err
	}

	for fileId, name := range files {
		if name == fileName {
			return fileId, nil
		}
	}

	reqBody := map[string]any{
		"name":                              fileName,
		"folder":                            map[string]any{},
		"@microsoft.graph.conflictBehavior": "rename",
	}
	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(string(jsonBody)))
	if err != nil {
		core.Log.Fatal("Error creating request: ", err)
		return "", err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		core.Log.Fatal("Error sending request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)
	if string(resBody) == "" {
		return "", fmt.Errorf("failed to create folder")
	}

	var resData interface{}
	json.Unmarshal(resBody, &resData)
	fileId := resData.(map[string]interface{})["id"].(string)
	core.Log.Info("Folder was successfully created with fileId: ", fileId)
	return fileId, nil
}

func uploadFile(accessToken string, filePath string, parentFileId string) error {
	core.Log.Info("Uploading file")
	var chunkSize int64 = 1024 * 1024 * 40 // 40MB
	apiUrl := "https://graph.microsoft.com/v1.0/me/drive/"
	if parentFileId != "" {
		apiUrl += "items/" + parentFileId + ":/"
	} else {
		apiUrl += "root:/"
	}

	fileStat, err := os.Stat(filePath)

	if err != nil {
		core.Log.Fatal("Error getting file stat: ", err)
		return err
	}

	destinationFileName := fileStat.Name()
	apiUrl += url.PathEscape(destinationFileName) + ":/createUploadSession"
	modFileDate := fileStat.ModTime().In(time.UTC).Format(time.RFC3339)
	body := map[string]any{
		"item": map[string]any{
			"@microsoft.graph.conflictBehavior": "replace",
			"name":                              destinationFileName,
			"fileSystemInfo": map[string]any{
				"createdDateTime":      modFileDate,
				"lastModifiedDateTime": modFileDate,
			},
		},
	}

	data, _ := json.Marshal(body)

	client := &http.Client{}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(string(data)))
	if err != nil {
		core.Log.Fatal("Error creating request: ", err)
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		core.Log.Fatal("Error sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)
	var resData interface{}
	json.Unmarshal(resBody, &resData)
	uploadUrl := resData.(map[string]interface{})["uploadUrl"].(string)

	contentSize := fileStat.Size()

	fileData, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		core.Log.Fatal("Error opening file: ", err)
		return err
	}
	defer fileData.Close()

	var currChunk int64 = 0
	for currChunk < contentSize {
		content := make([]byte, int(math.Min(float64(chunkSize), float64(contentSize-currChunk))))
		ret, _ := fileData.Read(content)

		req, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(content))
		if err != nil {
			core.Log.Error("Error creating request: ", err)
			return err
		}
		contentRange := fmt.Sprintf("bytes %v-%v/%v", currChunk, currChunk+int64(ret)-1, contentSize)
		req.Header.Set("Content-Range", contentRange)
		resp, err := client.Do(req)
		if err != nil {
			core.Log.Error("Error sending request: ", err)
			return err
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			core.Log.Error("Error reading response body: ", err)
			return err
		}
		core.Log.Info("Uploaded chunk ", currChunk, " - ", currChunk+int64(ret)-1, " of ", contentSize)
		currChunk += int64(ret)
	}
	core.Log.Info("File was successfully uploaded")
	return nil
}

func OneDriveUploader(filePath string, dirPath string) error {
	core.Log.Info("Uploading file to OneDrive")
	accessToken, err := getAccessToken()
	if err != nil || accessToken == "" {
		core.Log.Error("Error getting access token: ", err)
		return err
	}

	dirPaths := strings.Split(dirPath, "/")

	fileId, err := makeFolder(accessToken, dirPaths[0], "")
	if err != nil {
		core.Log.Error("Error creating folder: ", err)
		return err
	}

	// traverse through dir till we reach the required folder
	// This is the only way graph api allows to get the required folder
	for i := 1; i < len(dirPaths); i++ {
		fileId, err = makeFolder(accessToken, dirPaths[i], fileId)
		if err != nil {
			core.Log.Error("Error creating folder: ", err)
			return err
		}
	}

	return uploadFile(accessToken, filePath, fileId)
}
