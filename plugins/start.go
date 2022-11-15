package plugins

import (
	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func startHelper(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /start")
	_, err := b.SendMessage(chat.Id, "Hey there! I am alive", &gotgbot.SendMessageOpts{})
	return err
}
