package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func addMaintainer(maintainerName string, userId int, device string) error {
	_, err := core.Collection.InsertOne(context.TODO(), bson.M{
		"name":          maintainerName,
		"user_id":       userId,
		"devices":       []string{device},
		"is_maintainer": true,
		"is_admin":      false,
		"support_group": "",
	})
	return err
}

func addDevices(userId int, device string) error {
	filter := bson.M{"user_id": userId}
	update := bson.M{"$push": bson.M{"devices": device}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func AddMaintainer(maintainerName string, userId int, device string) error {
	var err error
	if IsMaintainer(userId) {
		core.Log.Warnln("User is already a maintainer. Adding device to maintainer's list")
		err = addDevices(userId, device)
	} else {
		core.Log.Infoln("Adding maintainer...")
		err = addMaintainer(maintainerName, userId, device)
	}
	return err
}

func IsMaintainer(userId int) bool {
	core.Log.Infoln("Checking if user is a maintainer...")
	filter := bson.M{"user_id": userId}
	var result bson.M
	err := core.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false
	}
	core.Log.Infoln("The user's maintainer status is", result["is_maintainer"])
	return result["is_maintainer"].(bool)
}

func RemoveMaintainer(userId int) error {
	core.Log.Infoln("Removing maintainer...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"is_maintainer": false}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}
