package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func GetDevices(userId int) []string {
	core.Log.Infoln("Getting devices...")
	filter := bson.M{"user_id": userId}
	var result bson.M
	err := core.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return []string{}
	}
	core.Log.Infoln("The user's devices are", result["devices"])
	devices := result["devices"].([]string)
	return devices
}

func RemoveDevice(userId int, device string) error {
	core.Log.Infoln("Removing device...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$pull": bson.M{"devices": device}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}
