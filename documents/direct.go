package documents

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
)

func directDownloader(url string, fileName string, t http.RoundTripper) (string, error) {
	core.Log.Infof("Fetching requests")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		core.Log.Errorf("error occurred: %v\n", err)
		return "", err
	}
	core.Log.Infof("Fetching response body")
	resp, err := t.RoundTrip(req)
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	core.Log.Infof("Downloading the file")
	fileHandle, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
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
	core.Log.Infof("Downloaded %v bytes\n", val)
	return fileName, nil
}

func DirectDownloader(url string) (string, error) {

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

	return directDownloader(url, fileName, http.DefaultTransport)
}
