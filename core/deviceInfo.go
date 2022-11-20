package core

import (
	"fmt"
	"strings"
)

type DeviceInfo struct {
	DeviceName         string            `json:"device_name"`
	Version            string            `json:"version"`
	BuildDate          string            `json:"build_date"`
	BuildType          string            `json:"build_type"` // user, userdebug, eng
	Flavour            string            `json:"flavour"`
	BuildFormat        map[string]string `json:"build_format"` // full, incremental, etc
	fullOtaPath        string
	incrementalOtaPath string
}

func ParseDeviceInfo(files []string) (DeviceInfo, error) {
	Log.Info("Parsing device info from ", files)
	var deviceInfo DeviceInfo = DeviceInfo{
		BuildFormat: map[string]string{},
	}

	deviceDets := strings.Split(files[0], "-")
	if len(deviceDets) < 9 || !strings.Contains(deviceDets[0], "FlamingoOS") || deviceDets[4] != "Official" {
		return deviceInfo, fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}

	deviceInfo.Version = strings.Replace(deviceDets[1], "v", "", 1)
	deviceInfo.DeviceName = deviceDets[2]
	deviceInfo.BuildType = deviceDets[3]
	deviceInfo.Flavour = deviceDets[5]
	deviceInfo.BuildDate = deviceDets[6]

	if deviceInfo.BuildType != "user" {
		return deviceInfo, fmt.Errorf("official builds needs user build type")
	}

	buildFormats := []string{"full", "incremental", "boot", "fastboot", "recovery"}
	for _, file := range files {
		for _, buildFormat := range buildFormats {
			if strings.Contains(file, buildFormat) {
				if !strings.Contains(file, "FlamingoOS") || !strings.Contains(file, "Official") || !strings.Contains(file, "user") {
					return deviceInfo, fmt.Errorf("invalid file. %s isn't a flamingoOS file", file)
				}
				fileName := strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
				uploadFolder := Branch + "/" + deviceInfo.DeviceName + "/" + deviceInfo.Flavour
				uploadUrl := BaseUrl + uploadFolder + "/" + fileName
				deviceInfo.BuildFormat[buildFormat] = uploadUrl
				if buildFormat == "full" {
					deviceInfo.fullOtaPath = file
				} else if buildFormat == "incremental" {
					deviceInfo.incrementalOtaPath = file
				}
			}
		}
	}

	if deviceInfo.BuildFormat["full"] == "" {
		return deviceInfo, fmt.Errorf("full build not found")
	}

	return deviceInfo, nil
}
