package plugins

import (
	"fmt"
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

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, e := b.SendMessage(chat.Id, "Ask an admin to add you as a maintainer", &gotgbot.SendMessageOpts{})
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
	core.Log.Infoln("Recieved request to handle /remove")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	userId := ctx.EffectiveUser.Id
	if replyMessage != nil {
		userId = replyMessage.From.Id
	}
	msg, err := b.SendMessage(chat.Id, "Removing the maintainer", &gotgbot.SendMessageOpts{})

	e := database.RemoveMaintainer(userId)
	if e != nil {
		msg.EditText(b, "Error adding maintainer. Please try again later.", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully removed the maintainer", &gotgbot.EditMessageTextOpts{})

	return err
}

func removeDevicesHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /dropDevice")

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

func promoteAdminHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /promote")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, e := b.SendMessage(chat.Id, "Reply to a message from the user you want to add as an admin.", &gotgbot.SendMessageOpts{})
		return e
	}

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, e := b.SendMessage(chat.Id, "Ask an admin to promote you", &gotgbot.SendMessageOpts{})
		return e
	}

	msg, err := b.SendMessage(chat.Id, "Promoting user", &gotgbot.SendMessageOpts{})

	e := database.PromoteAdmin(replyMessage.From.Id)
	if e != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully promoted user", &gotgbot.EditMessageTextOpts{})

	return err
}

func demoteAdminHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /demote")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, e := b.SendMessage(chat.Id, "Reply to a message from the user you want to remove as an admin.", &gotgbot.SendMessageOpts{})
		return e
	}

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, e := b.SendMessage(chat.Id, "You aren't an admin?!", &gotgbot.SendMessageOpts{})
		return e
	}

	msg, err := b.SendMessage(chat.Id, "Demoting user", &gotgbot.SendMessageOpts{})

	e := database.DemoteAdmin(replyMessage.From.Id)
	if e != nil {
		msg.EditText(b, "Something went wrong while demoting user", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully demoted user", &gotgbot.EditMessageTextOpts{})

	return err
}

func addSupportGroupHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /addGroup")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	userId := ctx.EffectiveUser.Id
	if replyMessage != nil {
		userId = replyMessage.From.Id
	}
	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	if len(args) == 0 {
		_, e := b.SendMessage(chat.Id, "Please provide a support group", &gotgbot.SendMessageOpts{})
		return e
	}

	supportGroup := args[0]

	msg, err := b.SendMessage(chat.Id, "Adding the support group", &gotgbot.SendMessageOpts{})

	e := database.AddSupportGroup(userId, supportGroup)
	if e != nil {
		msg.EditText(b, "Something went wrong while adding support group", &gotgbot.EditMessageTextOpts{})
	}

	msg.EditText(b, "Successfully added support group", &gotgbot.EditMessageTextOpts{})

	return err
}

func getAllMaintainers(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /getMaintainers")

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(chat.Id, "This is admin only", &gotgbot.SendMessageOpts{})
		return err
	}
	maintainers := database.GetAllMaintainers()
	msgTxt := fmt.Sprintf("The maintainers are %#v", maintainers)
	_, err := b.SendMessage(chat.Id, msgTxt, &gotgbot.SendMessageOpts{})
	return err
}
