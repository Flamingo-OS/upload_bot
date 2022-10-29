package core

import (
	"fmt"
	"strings"
)

type DeviceInfo struct {
	DeviceName  string `json:"device_name"`
	Version     string `json:"version"`
	BuildDate   string `json:"build_date"`
	BuildType   string `json:"build_type"` // user, userdebug, eng
	Flavour     string `json:"flavour"`
	BuildFormat string `json:"build_format"` // full, incremental, etc
}

func ParseDeviceInfo(fileName string) (DeviceInfo, error) {
	deviceDets := strings.Split(fileName, "-")
	if len(deviceDets) < 9 {
		return DeviceInfo{}, fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}
	if deviceDets[0] != "FlamingoOS" && deviceDets[4] != "Official" {
		return DeviceInfo{}, fmt.Errorf("invalid file. This isn't a flamingoOS file")
	}

	version := strings.Replace(deviceDets[1], "v", "", 1)

	return DeviceInfo{
		Version:     version,
		DeviceName:  deviceDets[2],
		BuildType:   deviceDets[3],
		Flavour:     deviceDets[5],
		BuildDate:   deviceDets[6],
		BuildFormat: deviceDets[8],
	}, nil
}
