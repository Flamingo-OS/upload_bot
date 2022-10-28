package core

import (
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func BotInit(config *BotConfig) (b *gotgbot.Bot, updater ext.Updater, err error) {
	b, err = gotgbot.NewBot(config.BotToken, &gotgbot.BotOpts{
		Client: http.Client{},
	})
	if err != nil {
		Log.Fatalf("Failed to create new bot due to %s\n", err.Error())
		return nil, updater, err
	}
	updater = ext.NewUpdater(nil)
	err = updater.StartPolling(b, &ext.PollingOpts{})
	if err != nil {
		Log.Errorf("Failed to start polling due to %s\n", err.Error())
	}
	Log.Info("The upload bot is up and running")
	return b, updater, nil
}
