package core

import (
	"fmt"
	"strings"
)

type DeviceInfo struct {
	DeviceName  string            `json:"device_name"`
	Version     string            `json:"version"`
	BuildDate   string            `json:"build_date"`
	BuildType   string            `json:"build_type"` // user, userdebug, eng
	Flavour     string            `json:"flavour"`
	BuildFormat map[string]string `json:"build_format"` // full, incremental, etc
}

func ParseDeviceInfo(fileName []string) (DeviceInfo, error) {
	Log.Info("Parsing device info from ", fileName)
	var deviceInfo DeviceInfo = DeviceInfo{
		BuildFormat: map[string]string{},
	}
	buildFormats := []string{"full", "incremental", "boot", "fastboot", "recovery"}
	for _, file := range fileName {
		for _, buildFormat := range buildFormats {
			if strings.Contains(file, buildFormat) {
				deviceInfo.BuildFormat[buildFormat] = file
			}
		}
	}

	if deviceInfo.BuildFormat["full"] == "" {
		return deviceInfo, fmt.Errorf("full build not found")
	}

	deviceDets := strings.Split(fileName[0], "-")
	if len(deviceDets) < 9 {
		return deviceInfo, fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}
	if strings.Contains(deviceDets[0], "FlamingoOS") && deviceDets[4] != "Official" {
		return deviceInfo, fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}

	deviceInfo.Version = strings.Replace(deviceDets[1], "v", "", 1)
	deviceInfo.DeviceName = deviceDets[2]
	deviceInfo.BuildType = deviceDets[3]
	deviceInfo.Flavour = deviceDets[5]
	deviceInfo.BuildDate = deviceDets[6]

	return deviceInfo, nil
}
