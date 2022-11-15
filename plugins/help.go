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
	/cancel <ID> - Cancel an upload
	/help - Show this message
	/ping - Check how slow I am today ;)
	/release - Upload a release
			It also has certain flags like below
			-n or --notes - addes extra notes to release posts
	/start - Check if I am alive? 

	These below are specific to maintainers. Alternatively admins can use by replying to a user
	/addGroup <link> - add a support group link
	/dropDevice <devices separated by space> - drop device from maintainence
	/remove - to drop yourself from maintainer status

	These below are specific to admins
	/add <devices> - add a new maintainer. Reply to the user
	/demote - demote the user being replied to an admin
	/promote - promote the user being replied to an admin
	`
	_, err := b.SendMessage(chat.Id, helpTxt, &gotgbot.SendMessageOpts{})
	return err
}
