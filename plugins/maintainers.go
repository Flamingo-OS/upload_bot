package plugins

import (
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func addMaintainerHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /add")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, e := b.SendMessage(chat.Id, "Reply to a message from the user you want to add as maintainer.", &gotgbot.SendMessageOpts{})
		return e
	}
	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	if len(args) == 0 {
		_, e := b.SendMessage(chat.Id, "Please provide atleast one devide", &gotgbot.SendMessageOpts{})
		return e
	}

	msg, err := b.SendMessage(chat.Id, "Adding a new maintainer", &gotgbot.SendMessageOpts{})

	userName := replyMessage.From.FirstName
	userId := replyMessage.From.Id

	e := database.AddMaintainer(userId, userName, args)
	if e != nil {
		msg.EditText(b, "Error adding maintainer. Please try again later.", &gotgbot.EditMessageTextOpts{})
	}

	return err
}
