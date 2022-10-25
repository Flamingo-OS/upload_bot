package main

import (
	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
	"github.com/Flamingo-OS/upload-bot/documents"
	"github.com/Flamingo-OS/upload-bot/plugins"
)

// various constants. Might me needed to change in future
const FILENAME string = "config.json"

func init() {
	// init the logger
	core.InitLogger()

	// extract config
	core.Config = core.NewBotConfig(FILENAME)

	// init a map to store and manage cancellable tasks
	core.CancelTasks = core.NewCancelCmd()

	// connect to db
	database.Init()

	// GDrive
	err := documents.NewGdrive()
	if err != nil {
		core.Log.Errorln(err)
		return
	}

}

func main() {
	_, updater, err := core.BotInit(core.Config)
	if err != nil {
		core.Log.Errorln(err)
		return
	}

	// Load up our custom handlers
	plugins.Main(updater.Dispatcher)
	updater.Idle()
}
