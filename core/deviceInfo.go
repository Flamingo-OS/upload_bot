package core

import (
	"fmt"
	"strings"
)

type DeviceInfo struct {
	DeviceName string `json:"device_name"`
	Version    string `json:"version"`
	BuildDate  string `json:"build_date"`
	BuildType  string `json:"build_type"` // user, userdebug, eng
	Flavour    string `json:"flavour"`
}

func ParseDeviceInfo(fileName []string) (DeviceInfo, string, string, error) {
	Log.Info("Parsing device info from ", fileName)
	fullOtaFile := ""
	incrementalOtaFile := ""
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
		Version:    version,
		DeviceName: deviceDets[2],
		BuildType:  deviceDets[3],
		Flavour:    deviceDets[5],
		BuildDate:  deviceDets[6],
	}, fullOtaFile, incrementalOtaFile, nil
}
