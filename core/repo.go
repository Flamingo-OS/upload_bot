package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CreateOTACommit(deviceInfo DeviceInfo, fullOtaFile string, incrementalOtaFile string, dumpPath string) error {
	Log.Info("Creating OTA commits")
	branch := "main"
	clonePath := dumpPath + "OTA/"
	otaPath := deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "ota.json"
	incrementalOtaPath := deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "incremental_ota.json"
	changelogPath := ""

	r, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "npv12",
			Password: Config.GithubToken,
		},
		URL:           "https://github.com/FlamingoOS-Devices/ota",
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})

	if err != nil {
		Log.Error("Error while cloning repo: %s", err)
		return err
	}

	formatTime := time.Now().Format("2006-01-02")

	changeLog, _ := CreateChangelog(deviceInfo.DeviceName, false)
	changelogPath = deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "changelog_" + strings.ReplaceAll(formatTime, "-", "_")
	writeToFile(clonePath+changelogPath, changeLog)
	ota, _ := CreateOtaJson(fullOtaFile, deviceInfo, dumpPath)
	jsonData, err := json.MarshalIndent(ota, "", "  ")
	if err != nil {
		Log.Error("Error while marshalling json: %s", err)
		return err
	}
	writeToFile(clonePath+otaPath, string(jsonData))
	if incrementalOtaFile != "" {
		ota, _ := CreateOtaJson(incrementalOtaFile, deviceInfo, dumpPath)
		jsonData, err := json.MarshalIndent(ota, "", "  ")
		if err != nil {
			Log.Error("Error while marshalling json: %s", err)
			return err
		}

		writeToFile(clonePath+incrementalOtaPath, string(jsonData))
	}

	w, err := r.Worktree()
	if err != nil {
		Log.Error("Error while getting worktree: %s", err)
		return err
	}

	w.Add(changelogPath)
	w.Add(otaPath)
	w.Add(incrementalOtaPath)

	commitMsg := fmt.Sprintf("%s: update %s", deviceInfo.DeviceName, formatTime)
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "github-actions[bot]",
			Email: "github-actions[bot]@users.noreply.github.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		Log.Errorf("Error: %v", err)
		return err
	}

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "npv12",
			Password: Config.GithubToken,
		}})

	if err != nil {
		Log.Errorf("Error: %v", err)
		return err
	}

	return nil
}
