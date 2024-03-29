package core

import (
	"net/http"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v2"
)

// various publically accisible vars. Make sure all of them are initialized in init()
var Config *BotConfig            // contains various configs loaded
var Log *zap.SugaredLogger       // logger
var CancelTasks *CancelCmds      // map to store and manage cancellable tasks
var Collection *mongo.Collection // mongodb collection
var DriveService *drive.Service  // Gdrive service
var DriveClient *http.Client     // Gdrive client
var Mut sync.Mutex

const DeviceOrg = "FlamingoOS-Devices"
const MainOrg = "Flamingo-OS"
const Branch = "A13"
const BaseUrl = "https://downloads.e11z.net/d/flamingo/"
const UpdateChannelId = -1001655954786
