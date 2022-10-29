package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type OTA struct {
	Version             string            `json:"version"`
	Date                string            `json:"date"`
	DownloadSources     map[string]string `json:"download_sources"`
	FileName            string            `json:"file_name"`
	FileSize            string            `json:"file_size"`
	ShaSum              string            `json:"sha_512"`
	PreBuildIncremental string            `json:"pre_build_incremental,omitempty"`
}

func CreateOtaJson(zipFilePath string, deviceInfo DeviceInfo) (OTA, error) {
	Log.Info("Creating OTA json for: ", zipFilePath)
	fileStat, err := os.Stat(zipFilePath)
	if err != nil {
		Log.Error("Error while getting file stats:", err)
		return OTA{}, err
	}
	a, err := UnzipFile(zipFilePath, strings.Trim(fileStat.Name(), ".zip"))
	if err != nil {
		Log.Error("Error while unzipping file:", err)
		return OTA{}, err
	}
	file, err := os.Open(fmt.Sprintf("%s/META-INF/com/android/metadata", a))
	if err != nil {
		Log.Error("Error while opening file:", err)
		return OTA{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	postTimestamp := ""
	preBuildIncremental := ""
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "post-timestamp=") {
			postTimestamp = strings.Replace(text, "post-timestamp=", "", -1) + "000"
		} else if strings.Contains(text, "pre-build-incremental=") {
			preBuildIncremental = strings.Replace(text, "pre-build-incremental=", "", -1)
		}
	}

	sha_512, err := FindShaSum(zipFilePath)
	if err != nil {
		Log.Error("Error while finding sha512:", err)
		return OTA{}, err
	}

	url := BaseUrl + Branch + "/" + deviceInfo.DeviceName + "/" + deviceInfo.Flavour + "/" + fileStat.Name()

	ota := OTA{
		Version:             deviceInfo.Version,
		Date:                postTimestamp,
		DownloadSources:     map[string]string{"OneDrive": url},
		FileName:            fileStat.Name(),
		FileSize:            fmt.Sprint(fileStat.Size()),
		ShaSum:              sha_512,
		PreBuildIncremental: preBuildIncremental,
	}

	if err := scanner.Err(); err != nil {
		Log.Error("Error while scanning file:", err)
		return OTA{}, err
	}

	Log.Info("Ota json was created as:", ota)

	return ota, nil
}
