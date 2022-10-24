package main

import (
	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/plugins"
)

// various constants. Might me needed to change in future
const FILENAME string = "config.json"

func main() {
	// enable logger
	core.InitLogger()

	// connect to tg
	core.Config = core.NewBotConfig()
	core.Config.ReadConfig(FILENAME)
	_, updater, err := core.BotInit(core.Config)
	if err != nil {
		core.Log.Errorln(err)
		return
	}

	// Load up our custom handlers
	plugins.Main(updater.Dispatcher)
	updater.Idle()
}
