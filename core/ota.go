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

func CreateOtaJson(zipFilePath string) (OTA, error) {
	fileStat, err := os.Stat(zipFilePath)
	if err != nil {
		Log.Error("Error while getting file stats: %s", err)
		return OTA{}, err
	}
	a, err := UnzipFile(zipFilePath, fileStat.Name())
	if err != nil {
		Log.Error("Error while unzipping file: %s", err)
		return OTA{}, err
	}
	file, err := os.Open(fmt.Sprintf("%s/META-INF/com/android/metadata", a))
	if err != nil {
		Log.Error("Error while opening file: %s", err)
		return OTA{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	postTimestamp := ""
	preBuildIncremental := ""
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "post-timestamp=") {
			postTimestamp = strings.Replace(text, "post-timestamp=", "", -1)
		} else if strings.Contains(text, "pre-build-incremental=") {
			preBuildIncremental = strings.Replace(text, "pre-build-incremental=", "", -1)
		}
	}

	version := strings.Replace(strings.Split(fileStat.Name(), "-")[1], "v", "", 1)
	sha_512, err := FindShaSum(zipFilePath)
	if err != nil {
		Log.Error("Error while finding sha512: %s", err)
		return OTA{}, err
	}

	ota := OTA{
		Version:             version,
		Date:                postTimestamp,
		DownloadSources:     map[string]string{"OneDrive": "https://sourceforge.net/projects/"},
		FileName:            fileStat.Name(),
		FileSize:            fmt.Sprint(fileStat.Size()),
		ShaSum:              sha_512,
		PreBuildIncremental: preBuildIncremental,
	}

	if err := scanner.Err(); err != nil {
		Log.Error("Error while scanning file: %s", err)
		return OTA{}, err
	}

	return ota, nil
}
