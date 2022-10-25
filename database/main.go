package database

import (
	"context"

	"github.com/Flamingo-OS/upload-bot/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "kosp"
const colName = "maintainers"

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
