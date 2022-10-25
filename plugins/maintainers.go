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

	msg.EditText(b, "Successfully added a maintainer", &gotgbot.EditMessageTextOpts{})
	return err
}

func removeMaintainerHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /add")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, e := b.SendMessage(chat.Id, "Reply to a message from the user you want to add as maintainer.", &gotgbot.SendMessageOpts{})
		return e
	}
	msg, err := b.SendMessage(chat.Id, "Removing the maintainer", &gotgbot.SendMessageOpts{})
	userId := replyMessage.From.Id

	e := database.RemoveMaintainer(userId)
	if e != nil {
		msg.EditText(b, "Error adding maintainer. Please try again later.", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully removed the maintainer", &gotgbot.EditMessageTextOpts{})

	return err
}

func removeDevicesHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /add")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	userId := ctx.EffectiveUser.Id
	if replyMessage != nil {
		userId = replyMessage.From.Id
	}
	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	if len(args) == 0 {
		_, e := b.SendMessage(chat.Id, "Please provide atleast one devide", &gotgbot.SendMessageOpts{})
		return e
	}

	msg, err := b.SendMessage(chat.Id, "Removing the device(s)", &gotgbot.SendMessageOpts{})

	e := database.RemoveDevice(userId, args)
	if e != nil {
		msg.EditText(b, "Something went wrong while removing device", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully removed the device(s)", &gotgbot.EditMessageTextOpts{})

	return err
}
