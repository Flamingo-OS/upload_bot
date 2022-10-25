package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func addMaintainer(maintainerName string, userId int64, devices []string) error {
	_, err := core.Collection.InsertOne(context.TODO(), Maintainers{
		MaintainerName: maintainerName,
		UserId:         userId,
		Devices:        devices,
		IsMaintainer:   true,
		IsAdmin:        false,
		SupportGroup:   "",
	})
	return err
}

func addDevices(userId int64, devices []string) error {
	originalDevices := GetDevices(userId)
	for _, origDevice := range originalDevices {
		for idx, newDevice := range devices {
			if origDevice == newDevice {
				devices = append(devices[:idx], devices[idx+1:]...)
			}
		}
	}

	if len(devices) == 0 {
		return nil
	}

	filter := bson.M{"user_id": userId}
	update := bson.M{"$push": bson.M{"devices": bson.M{"$each": devices}}}
	_, err := core.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func AddMaintainer(userId int64, maintainerName string, devices []string) error {
	var err error
	if IsMaintainer(userId) {
		core.Log.Infoln("User is already a maintainer. Adding device to maintainer's list")
		err = addDevices(userId, devices)
	} else {
		core.Log.Infoln("Adding maintainer...")
		err = addMaintainer(maintainerName, userId, devices)
	}
	return err
}

func IsMaintainer(userId int64) bool {
	core.Log.Infoln("Checking if user is a maintainer...")
	filter := bson.M{"user_id": userId}
	var result bson.M
	err := core.Collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return false
	}
	core.Log.Infoln("The user's maintainer status is", result["is_maintainer"])
	return result["is_maintainer"].(bool)
}

func RemoveMaintainer(userId int64) error {
	core.Log.Infoln("Removing maintainer...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"is_maintainer": false}}
	_, err := core.Collection.UpdateOne(context.Background(), filter, update)
	return err
}
