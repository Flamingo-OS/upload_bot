package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

var repos []string

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
		Log.Error("Status: %d", resp.StatusCode)
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

func findNextPage(nextPosUrl string) string {
	for _, link := range strings.Split(nextPosUrl, ",") {
		if strings.Contains(link, "rel=\"next\"") {
			return strings.Trim(strings.Trim(strings.Trim(strings.Split(link, ";")[0], " "), "<"), ">")
		}
	}
	return ""
}

func findRepoUrls(url string, endDate time.Time) error {
	Log.Info("Finding repo urls for %s", url)
	var blacklist = []string{"vendor_prebuilts"}
	if url == "" {
		Log.Warn("Empty url")
		return nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Error while creating request: %s", err)
		return err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: %s", err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Error("Status: %d", resp.StatusCode)
		return err
	}

	var resBody []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	if err != nil {
		Log.Error("Error while decoding json: %s", err)
		return err
	}

	nextLink := resp.Header.Get("Link")
	if nextLink == "" {
		Log.Info("No more pages")
		return fmt.Errorf("no next link found")
	}

	for _, repo := range resBody {
		repoName := repo["name"].(string)
		repoPushed, _ := time.Parse(time.RFC3339, repo["pushed_at"].(string))
		if slices.Contains(blacklist, repoName) || (endDate.UTC().After(repoPushed.UTC())) {
			continue
		}
		commitsUrl := strings.Replace(repo["commits_url"].(string), "{/sha}", "", 1)
		repos = append(repos, commitsUrl)
	}

	findRepoUrls(findNextPage(nextLink), endDate)

	return nil
}

func findDeviceRepo(device string) (string, error) {
	Log.Info("Finding repo for device %s", device)
	apiUrl := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+user:%s+in:name+fork:true", device, DeviceOrg)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		Log.Error("Error while creating request: %s", err)
		return "", err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Info("Status: %d", resp.StatusCode)
		return "", err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		Log.Error("Error while decoding json: %s", err)
		return "", err
	}

	for _, deviceRepo := range result["items"].([]interface{}) {
		repoUrl := deviceRepo.(map[string]interface{})["name"].(string)
		splitRepoUrl := strings.Split(repoUrl, "_")
		if slices.Contains(splitRepoUrl, "device") && slices.Contains(splitRepoUrl, device) {
			return repoUrl, nil
		}
	}

	return "", fmt.Errorf("no device repo found")
}

func findCommits(url string, changelog *string, endDate time.Time) error {
	Log.Info("Finding commits for %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Error while creating request: %s", err)
		return err
	}

	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: %s", err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Error("Status: %d", resp.StatusCode)
		return err
	}

	var commits []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		Log.Error("Error while decoding json: %s", err)
		return err
	}

	for _, commit := range commits {
		commitDate, _ := time.Parse(time.RFC3339, commit["commit"].(map[string]interface{})["author"].(map[string]interface{})["date"].(string))
		if endDate.UTC().After(commitDate.UTC()) {
			break
		}
		commitMessage := strings.Split(commit["commit"].(map[string]interface{})["message"].(string), "\n")[0]
		*changelog += commitMessage + "\n"
	}

	return nil
}
