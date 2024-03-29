package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func AddSupportGroup(userId int64, supportGroup string) error {
	core.Log.Infoln("Adding support group")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"support_group": supportGroup}}
	_, err := core.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func GetSupportGroup(userId int64) (string, error) {
	core.Log.Infoln("Getting support group")
	var user Maintainers
	filter := bson.M{"user_id": userId}
	err := core.Collection.FindOne(context.Background(), filter).Decode(&user)
	return user.SupportGroup, err
}
