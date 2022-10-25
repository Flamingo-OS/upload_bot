package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "kosp"
const colName = "maintainers"

type Maintainers struct {
	Id             primitive.ObjectID `bson:"_id"`
	MaintainerName string             `bson:"name"`
	UserId         int64              `bson:"user_id"`
	Devices        []string           `bson:"devices"`
	IsMaintainer   bool               `bson:"is_maintainer"`
	IsAdmin        bool               `bson:"is_admin"`
	SupportGroup   string             `bson:"support_group"`
}

// connect to the database
func Init() {
	core.Log.Infoln("Connecting to database...")
	clientOption := options.Client().ApplyURI(core.Config.MongoDbConnectionString)

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		core.Log.DPanicln("Error while connecting to database", err)
		panic(err)
	}

	core.Log.Infoln("Connected to MongoDB!")
	core.Collection = client.Database(dbName).Collection(colName)
}
