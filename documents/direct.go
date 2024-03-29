package documents

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
)

func downloadSaver(resp *http.Response, fileName string, dumpPath string) (string, error) {
	core.Log.Info("Downloading the file")
	filePath := dumpPath + fileName
	fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return "", err
	}
	defer fileHandle.Close()
	val, err := io.Copy(fileHandle, resp.Body)
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return "", err
	}
	core.Log.Info("Downloaded bytes:", val)
	return filePath, nil
}

func DirectDownloader(url string, dumpPath string) (string, error) {

	var fileName string = ""
	for _, item := range strings.Split(url, "/") {
		if strings.Contains(item, "Flamingo") && strings.Contains(item, "Official") {
			fileName = item
		}
	}
	if fileName == "" {
		core.Log.Errorf("Unable to find file name")
		return "", fmt.Errorf("unable to find file name")
	}

	core.Log.Info("Fetching requests")
	resp, err := http.Get(url)
	if err != nil {
		core.Log.Errorf("Unable to fetch requests: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	return downloadSaver(resp, fileName, dumpPath)
}
