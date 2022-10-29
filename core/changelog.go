package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func findLastDate(device string, isVanilla bool) (time.Time, error) {
	Log.Info("Finding last date for device %s", device)
	buildType := "GApps"
	if isVanilla {
		buildType = "Vanilla"
	}
	apiUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/ota/main/%s/%s/ota.json", DeviceOrg, device, buildType)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		Log.Error("Error while creating request: %s", err)
		return time.Time{}, err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: %s", err)
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		fmt.Println("Status:", resp.StatusCode)
		return time.Time{}, err
	}
	var ota OTA
	err = json.NewDecoder(resp.Body).Decode(&ota)
	if err != nil {
		Log.Error("Error while decoding json: %s", err)
		return time.Time{}, err
	}
	timestamp, _ := strconv.ParseInt(ota.Date, 10, 64)
	date := time.Unix(timestamp/1000, 0)
	Log.Info("Found last date for device %s: %s", device, date)
	return date, nil
}
