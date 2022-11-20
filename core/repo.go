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

func pushOTARepo(deviceInfo DeviceInfo, clonePath string, otaPath string, otaData string, changelogPath string, changelogData string, incrementalOtaPath string, incrementalOtaData string) error {
	branch := "main"
	formatTime := time.Now().Format("2006-01-02")
	Mut.Lock()
	defer Mut.Unlock()

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

	w, err := r.Worktree()
	if err != nil {
		Log.Error("Error while getting worktree: %s", err)
		return err
	}

	Log.Infof("Pulling latest changes for device: %s", deviceInfo.DeviceName)
	w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})
	Log.Infof("Staging new changes for : %s", deviceInfo.DeviceName)
	writeToFile(clonePath+changelogPath, changelogData)
	w.Add(changelogPath)
	writeToFile(clonePath+otaPath, otaData)
	w.Add(otaPath)
	if incrementalOtaData != "" {
		writeToFile(clonePath+incrementalOtaPath, incrementalOtaData)
		w.Add(incrementalOtaPath)
	}

	commitMsg := fmt.Sprintf("%s: update %s", deviceInfo.DeviceName, formatTime)
	Log.Infof("Pushing OTA for : %s", deviceInfo.DeviceName)
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

	return r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "npv12",
			Password: Config.GithubToken,
		}})
}

func CreateOTACommit(deviceInfo DeviceInfo, dumpPath string) error {
	Log.Info("Creating OTA commits")

	clonePath := dumpPath + "OTA/"
	otaPath := deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "ota.json"
	incrementalOtaPath := deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "incremental_ota.json"
	changelogPath := ""

	formatTime := time.Now().Format("2006-01-02")

	changeLog, _ := CreateChangelog(deviceInfo.DeviceName, false)
	changelogPath = deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + "changelog_" + strings.ReplaceAll(formatTime, "-", "_")
	ota, _ := CreateOtaJson(deviceInfo.BuildFormat["full"], deviceInfo, dumpPath)
	otaJson, err := json.MarshalIndent(ota, "", "  ")
	if err != nil {
		Log.Error("Error while marshalling json: %s", err)
		return err
	}
	otaData := string(otaJson)
	incrementalOtaData := ""
	if deviceInfo.BuildFormat["incremental"] != "" {
		ota, _ := CreateOtaJson(deviceInfo.BuildFormat["incremental"], deviceInfo, dumpPath)
		incrementalOtaJson, err := json.MarshalIndent(ota, "", "  ")
		if err != nil {
			Log.Error("Error while marshalling json: %s", err)
			return err
		}
		incrementalOtaData = string(incrementalOtaJson)
	}
	return pushOTARepo(deviceInfo, clonePath, otaPath, otaData, changelogPath, changeLog, incrementalOtaPath, incrementalOtaData)
}
