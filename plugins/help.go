package plugins

import (
	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func helpHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /help")
	var helpTxt = `
	Hey, I'm Flamingo upload bot. I upload OTA files from your drives to hosted drive
	/help - Show this message
	/ping - Check how slow I am today ;)
	/start - Check if I am alive? 
	`
	_, e := b.SendMessage(chat.Id, helpTxt, &gotgbot.SendMessageOpts{})
	return e
}
