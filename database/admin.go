package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func PromoteAdmin(userId int) error {
	core.Log.Infoln("Promoting maintainer...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"is_admin": true}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func DemoteAdmin(userId int) error {
	core.Log.Infoln("Demoting maintainer...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"is_admin": false}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func IsAdmin(userId int) bool {
	core.Log.Infoln("Checking if user is an admin...")
	filter := bson.M{"user_id": userId}
	var result bson.M
	err := core.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false
	}
	core.Log.Infoln("The user's admin status is", result["is_admin"])
	return result["is_admin"].(bool)
}
