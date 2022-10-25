package core

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v3"
)

// various publically accisible vars. Make sure all of them are initialized in init()
var Config *BotConfig            // contains various configs loaded
var Log *zap.SugaredLogger       // logger
var CancelTasks *CancelCmds      // map to store and manage cancellable tasks
var Collection *mongo.Collection // mongodb collection
var DriveService *drive.Service  // Gdrive client
