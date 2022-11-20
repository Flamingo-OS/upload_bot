package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNotes(userId int64) (string, error) {
	core.Log.Infoln("Getting notes...")
	filter := bson.M{"user_id": userId}
	var result Maintainers
	err := core.Collection.FindOne(context.Background(), filter).Decode(&result)
	return result.Notes, err
}

func SetNotes(userId int64, notes string) error {
	core.Log.Infoln("Setting notes...")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"notes": notes}}
	_, err := core.Collection.UpdateOne(context.Background(), filter, update)
	return err
}
