package core

import (
	"fmt"
	"strings"
)

type DeviceInfo struct {
	DeviceName  string   `json:"device_name"`
	Version     string   `json:"version"`
	BuildDate   string   `json:"build_date"`
	BuildType   string   `json:"build_type"` // user, userdebug, eng
	Flavour     string   `json:"flavour"`
	BuildFormat []string `json:"build_format"` // full, incremental, etc
}

func ParseDeviceInfo(fileName []string) (DeviceInfo, string, string, error) {
	Log.Info("Parsing device info from ", fileName)
	buildFormats := []string{}
	fullOtaFile := ""
	incrementalOtaFile := ""
	for _, file := range fileName {
		if strings.Contains(file, "full") {
			buildFormats = append(buildFormats, "full")
			fullOtaFile = file
		} else if strings.Contains(file, "incremental") {
			buildFormats = append(buildFormats, "incremental")
			incrementalOtaFile = file
		} else if strings.Contains(file, "fastboot") {
			buildFormats = append(buildFormats, "fastboot")
		} else if strings.Contains(file, "boot") {
			buildFormats = append(buildFormats, "boot")
		}
	}

	if fullOtaFile == "" {
		return DeviceInfo{}, "", "", fmt.Errorf("full build not found")
	}

	deviceDets := strings.Split(fileName[0], "-")
	if len(deviceDets) < 9 {
		return DeviceInfo{}, "", "", fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}
	if strings.Contains(deviceDets[0], "FlamingoOS") && deviceDets[4] != "Official" {
		return DeviceInfo{}, "", "", fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}

	version := strings.Replace(deviceDets[1], "v", "", 1)

	return DeviceInfo{
		Version:     version,
		DeviceName:  deviceDets[2],
		BuildType:   deviceDets[3],
		Flavour:     deviceDets[5],
		BuildDate:   deviceDets[6],
		BuildFormat: buildFormats,
	}, fullOtaFile, incrementalOtaFile, nil
}
