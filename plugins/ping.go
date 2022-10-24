package plugins

import (
	"fmt"
	"time"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func pingHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /ping")
	start := time.Now()
	_, e := b.SendMessage(chat.Id, "Pinging.......", &gotgbot.SendMessageOpts{})
	if e != nil {
		core.Log.Fatalf("Failed to send message")
		return e
	}
	diff := time.Since(start).Milliseconds()

	_, e = b.SendMessage(chat.Id, fmt.Sprintf("Pong! %dms", diff), &gotgbot.SendMessageOpts{})

	return e
}