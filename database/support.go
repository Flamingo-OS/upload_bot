package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson"
)

func AddSupportGroup(userId int, supportGroup string) error {
	core.Log.Infoln("Adding support group")
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{"support_group": supportGroup}}
	_, err := core.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}
