package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

// stores all repos from where we can get the commits and the corresponding changelogs
// includes repos from flamingoOS org.
var repos map[string]string
var deviceRepos map[string]string

// find the last update date
func findLastDate(device string, flavour string) (time.Time, error) {
	Log.Info("Finding last date for device", device)
	apiUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/ota/main/%s/%s/ota.json", DeviceOrg, device, flavour)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		Log.Error("Error while creating request: ", err)
		return time.Time{}, err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: ", err)
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Error("Status:", resp.StatusCode)
		return time.Time{}, fmt.Errorf("could not find last date")
	}
	var ota OTA
	err = json.NewDecoder(resp.Body).Decode(&ota)
	if err != nil {
		Log.Error("Error while decoding json: ", err)
		return time.Time{}, err
	}
	timestamp, _ := strconv.ParseInt(ota.Date, 10, 64)
	date := time.Unix(timestamp/1000, 0)
	Log.Infof("Found last date for device %s: %s", device, date)
	return date, nil
}

// given the next url, extract url for next page
func findNextPage(nextPosUrl string) string {
	for _, link := range strings.Split(nextPosUrl, ",") {
		if strings.Contains(link, "rel=\"next\"") {
			return strings.Trim(strings.Trim(strings.Trim(strings.Split(link, ";")[0], " "), "<"), ">")
		}
	}
	return ""
}

// find all repos from Flamingo-OS org
func findRepoUrls(url string, endDate time.Time) error {
	Log.Info("Finding repo urls for ", url)
	var blacklist = []string{"vendor_prebuilts", "upload_bot"}
	if url == "" {
		Log.Warn("Empty url")
		return nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Error while creating request: ", err)
		return err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Error("Status:", resp.StatusCode)
		return err
	}

	var resBody []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	if err != nil {
		Log.Error("Error while decoding json: ", err)
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
		repos[commitsUrl] = ""
	}

	findRepoUrls(findNextPage(nextLink), endDate)

	return nil
}

// Find the device repo given device name
func findDeviceRepo(device string) (string, error) {
	Log.Info("Finding repo for device ", device)
	apiUrl := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+user:%s+in:name+fork:true", device, DeviceOrg)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		Log.Error("Error while creating request: ", err)
		return "", err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: ", err)
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
		Log.Error("Error while decoding json: ", err)
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

// Find the required commtis
func findCommits(url string, changelog *string, endDate time.Time) error {
	Log.Info("Finding commits for ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Error while creating request: ", err)
		return err
	}

	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Error("Status: ", resp.StatusCode)
		return err
	}

	var commits []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		Log.Error("Error while decoding json: ", err)
		return err
	}

	for _, commit := range commits {
		commitDate, _ := time.Parse(time.RFC3339, commit["commit"].(map[string]interface{})["committer"].(map[string]interface{})["date"].(string))
		if endDate.UTC().After(commitDate.UTC()) {
			return nil
		}
		commitMessage := strings.Split(commit["commit"].(map[string]interface{})["message"].(string), "\n")[0]
		*changelog += commitMessage + "\n"
	}

	nextLink := resp.Header.Get("Link")
	findCommits(findNextPage(nextLink), changelog, endDate)
	return nil
}

// find all deps for a device
func findDependencies(device string) error {
	Log.Info("Finding dependencies for device ", device)
	apiUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/flamingo.dependencies", DeviceOrg, device, Branch)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		Log.Error("Error while creating request: ", err)
		return err
	}
	req.Header.Set("Authorization", "token "+Config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("Error while sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) || (err != nil) {
		Log.Info("Status: ", resp.StatusCode)
		return err
	}

	var dependencies []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dependencies)
	if err != nil {
		Log.Error("Error while decoding json: ", err)
		return err
	}

	for _, dependency := range dependencies {
		depName := dependency["repository"].(string)
		if dependency["branch"] != nil {
			continue // skip non-main branches
		}
		depRemote := DeviceOrg
		if dependency["remote"] != nil {
			if dependency["remote"].(string) == "flamingo" {
				depRemote = MainOrg
			} else {
				continue // skip outside remotes
			}
		}

		findDependencies(depName)

		dep := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", depRemote, depName)
		deviceRepos[dep] = ""
	}

	return nil
}

// create changelog for a single repo
// it will all be writen to a single changelog string
// Fully multi threaded and async
func createChangelogs(r map[string]string, repo string, date time.Time, mut *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	var changelog string
	e := findCommits(repo, &changelog, date)
	if e != nil {
		Log.Error("Error while finding commits: ", e)
		return
	}
	if changelog != "" {
		mut.Lock()
		r[repo] = changelog
		mut.Unlock()
	}
}

func genChangelog(r map[string]string, ch *string) {
	keys := make([]string, 0, len(r))

	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, repoUrl := range keys {
		if r[repoUrl] != "" {
			splitUrl := strings.Split(repoUrl, "/")
			*ch += fmt.Sprintf("# %s #\n", splitUrl[len(splitUrl)-2])
			*ch += r[repoUrl]
			*ch += "\n"
		}
	}
}

// Main handler.
func CreateChangelog(deviceName string, flavour string) (string, error) {
	Log.Info("Creating changelog")
	var wg sync.WaitGroup

	date, err := findLastDate(deviceName, flavour)
	if err != nil {
		Log.Error("Error while finding last date: ", err)
		return "Initial release", err
	}
	repos = map[string]string{}
	deviceRepos = map[string]string{}
	findRepoUrls(fmt.Sprintf("https://api.github.com/orgs/%s/repos?type=all", MainOrg), date)
	deviceRepo, err := findDeviceRepo(deviceName)
	if err != nil {
		return "", err
	}
	deviceRepos[fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", DeviceOrg, deviceRepo)] = ""
	findDependencies(deviceRepo)

	for repoUrl := range repos {
		go createChangelogs(repos, repoUrl, date, &Mut, &wg)
		wg.Add(1)
	}
	for repoUrl := range deviceRepos {
		go createChangelogs(deviceRepos, repoUrl, date, &Mut, &wg)
		wg.Add(1)
	}
	wg.Wait()
	var changelogs string = "### Source side changes ###\n\n"
	genChangelog(repos, &changelogs)
	changelogs += "\n### Device side changes ###\n\n"
	genChangelog(deviceRepos, &changelogs)
	return changelogs, nil
}
