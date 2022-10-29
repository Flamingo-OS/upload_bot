package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CreateOTACommit(deviceInfo DeviceInfo, fullOtaFile string, incrementalOtaFile string) error {
	branch := "dev"
	clonePath := DumpPath + "OTA/"
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
	writeToFile(clonePath+"ota/"+deviceInfo.DeviceName+"/"+deviceInfo.Flavour+"/"+"changelog"+formatTime, changeLog)
	ota, _ := CreateOtaJson(fullOtaFile)
	jsonData, err := json.MarshalIndent(ota, "", "  ")
	if err != nil {
		Log.Error("Error while marshalling json: %s", err)
		return err
	}
	writeToFile(clonePath+"ota/"+deviceInfo.DeviceName+"/"+deviceInfo.Flavour+"/"+"ota.json", string(jsonData))
	if incrementalOtaFile != "" {
		ota, _ := CreateOtaJson(incrementalOtaFile)
		jsonData, err := json.MarshalIndent(ota, "", "  ")
		if err != nil {
			Log.Error("Error while marshalling json: %s", err)
			return err
		}

		writeToFile(clonePath+"ota/"+deviceInfo.DeviceName+"/"+deviceInfo.Flavour+"/"+"incremental_ota.json", string(jsonData))
	}

	w, err := r.Worktree()
	if err != nil {
		Log.Error("Error while getting worktree: %s", err)
		return err
	}

	_, err = w.Status()
	if err != nil {
		Log.Error("Error while getting status: %s", err)
		return err
	}

	commitMsg := fmt.Sprintf("%s: update %s", deviceInfo.DeviceName, formatTime)
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "github-actions[bot]",
			Email: "github-actions[bot]@users.noreply.github.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "npv12",
			Password: Config.GithubToken,
		}})

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	return nil
}
