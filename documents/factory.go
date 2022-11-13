package documents

import (
	"fmt"
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
)

func DocumentFactory(url string, dumpPath string) (string, error) {
	if strings.Contains(url, "drive.google.com") {
		core.Log.Info("GDrive detected. Using GDrive module.")
		return gDriveDownloader(url, dumpPath)
	} else if strings.Contains(url, "mega.nz") {
		core.Log.Info("Mega detected. Aborting Download")
		return "", fmt.Errorf("mega download big no no")
	}
	core.Log.Infoln("Unknown URL. Attempting to download as is.")
	return DirectDownloader(url, dumpPath)
}
