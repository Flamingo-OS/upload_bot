package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func GetDevices(userId int64) []string {
	core.Log.Infoln("Getting devices...")
	filter := bson.M{"user_id": userId}
	var result Maintainers
	err := core.Collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		core.Log.Errorln(err)
		return []string{}
	}
	return result.Devices
}

func RemoveDevice(userId int64, devices []string) error {
	core.Log.Infoln("Removing device...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$pull": bson.M{"devices": bson.M{"$in": devices}}}
	_, err := core.Collection.UpdateOne(context.Background(), filter, update)
	return err
}
