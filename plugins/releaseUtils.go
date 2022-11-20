package plugins

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
)

// validates the release is indeed a flamingo OS file
// also creates and pushes OTA file
func validateRelease(fileNames []string, dumpPath string) (core.DeviceInfo, error) {
	deviceInfo, fullOtaFile, incrementalOtaFile, err := core.ParseDeviceInfo(fileNames)
	if err != nil {
		return core.DeviceInfo{}, err
	}
	core.Log.Info("Parsed device info:", deviceInfo)
	err = core.CreateOTACommit(deviceInfo, fullOtaFile, incrementalOtaFile, dumpPath)
	if err != nil {
		return deviceInfo, err
	}
	return deviceInfo, err
}

func CreateReleaseText(deviceInfo core.DeviceInfo, urls []string, maintainers []database.Maintainers, deviceSupport string, notes string) (string, error) {
	formatTime := time.Now().Format("2006-01-02")
	buildTime, _ := time.Parse("20060102", deviceInfo.BuildDate)
	av, err := strconv.ParseFloat(deviceInfo.Version, 64)
	if err != nil {
		core.Log.Error("Error while parsing android version: %s", err)
		return "", err
	}
	androidVersion := int(av) + 11
	changeLogUrl := fmt.Sprintf("https://raw.githubusercontent.com/FlamingoOS-Devices/OTA/main/%s/%s/changelog_%s", deviceInfo.DeviceName, deviceInfo.Flavour, strings.ReplaceAll(formatTime, "-", "_"))
	msgTxt := fmt.Sprintf(`FlamingoOS %s | Android %d | OFFICIAL | %s
	
		Maintainers: `, deviceInfo.Version, androidVersion, deviceInfo.Flavour)

	for _, maintainer := range maintainers {
		msgTxt += fmt.Sprintf(` [%s](tg://user?id=%d) `, maintainer.MaintainerName, maintainer.UserId)
	}

	msgTxt += fmt.Sprintf(`
		
		Device: %s
		Date: %s

		Downloads`, deviceInfo.DeviceName, buildTime.Format("02-01-2006"))

	for format, url := range deviceInfo.BuildFormat {
		msgTxt += fmt.Sprintf(`
		-	[%s](%s)	`, format, url)
	}
	msgTxt += fmt.Sprintf(`

		[Changelog](%s)
	
		Support Groups
		-	[Common](https://t.me/flamingo_common)
		-	[Updates](https://t.me/flamingo_updates)
	`, changeLogUrl)

	if deviceSupport != "" {
		msgTxt += fmt.Sprintf(`
		-	[Device Support](%s)`, deviceSupport)
	}

	if notes != "" {
		msgTxt += fmt.Sprintf(`
		Notes: %s`, notes)
	}

	return msgTxt, nil
}
