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
		_, err := b.SendMessage(chat.Id, "Reply to a message from the user you want to add as maintainer.", &gotgbot.SendMessageOpts{})
		return err
	}

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(chat.Id, "Ask an admin to add you as a maintainer", &gotgbot.SendMessageOpts{})
		return err
	}

	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	if len(args) == 0 {
		_, err := b.SendMessage(chat.Id, "Please provide atleast one devide", &gotgbot.SendMessageOpts{})
		return err
	}

	msg, err := b.SendMessage(chat.Id, "Adding a new maintainer", &gotgbot.SendMessageOpts{})
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	userName := replyMessage.From.FirstName
	userId := replyMessage.From.Id

	err = database.AddMaintainer(userId, userName, args)
	if err != nil {
		msg.EditText(b, "Error adding maintainer. Please try again later.", &gotgbot.EditMessageTextOpts{})
		return err
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
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	err = database.RemoveMaintainer(userId)
	if err != nil {
		msg.EditText(b, "Error removing maintainer. Please try again later.", &gotgbot.EditMessageTextOpts{})
		return err
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
		_, err := b.SendMessage(chat.Id, "Please provide atleast one devide", &gotgbot.SendMessageOpts{})
		return err
	}

	msg, err := b.SendMessage(chat.Id, "Removing the device(s)", &gotgbot.SendMessageOpts{})
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	err = database.RemoveDevice(userId, args)
	if err != nil {
		msg.EditText(b, "Something went wrong while removing device", &gotgbot.EditMessageTextOpts{})
		return err
	}

	msg.EditText(b, "Successfully removed the device(s)", &gotgbot.EditMessageTextOpts{})

	return err
}

func promoteAdminHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /promote")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, err := b.SendMessage(chat.Id, "Reply to a message from the user you want to add as an admin.", &gotgbot.SendMessageOpts{})
		return err
	}

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(chat.Id, "Ask an admin to promote you", &gotgbot.SendMessageOpts{})
		return err
	}

	msg, err := b.SendMessage(chat.Id, "Promoting user", &gotgbot.SendMessageOpts{})

	if err != nil {
		core.Log.Errorln("Error sending message", err)
		return err
	}

	err = database.PromoteAdmin(replyMessage.From.Id)
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	_, _, err = msg.EditText(b, "Successfully promoted user", &gotgbot.EditMessageTextOpts{})

	return err
}

func demoteAdminHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /demote")

	replyMessage := ctx.EffectiveMessage.ReplyToMessage
	if replyMessage == nil {
		_, err := b.SendMessage(chat.Id, "Reply to a message from the user you want to remove as an admin.", &gotgbot.SendMessageOpts{})
		return err
	}

	if !database.IsAdmin(ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(chat.Id, "You aren't an admin?!", &gotgbot.SendMessageOpts{})
		return err
	}

	msg, err := b.SendMessage(chat.Id, "Demoting user", &gotgbot.SendMessageOpts{})
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	err = database.DemoteAdmin(replyMessage.From.Id)
	if err != nil {
		msg.EditText(b, "Something went wrong while demoting user", &gotgbot.EditMessageTextOpts{})
		return err
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
		_, err := b.SendMessage(chat.Id, "Please provide a support group", &gotgbot.SendMessageOpts{})
		return err
	}

	supportGroup := args[0]

	msg, err := b.SendMessage(chat.Id, "Adding the support group", &gotgbot.SendMessageOpts{})
	if err != nil {
		msg.EditText(b, "Something went wrong while promoting user", &gotgbot.EditMessageTextOpts{})
		return err
	}

	err = database.AddSupportGroup(userId, supportGroup)
	if err != nil {
		msg.EditText(b, "Something went wrong while adding support group", &gotgbot.EditMessageTextOpts{})
		return err
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

func setNotesHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	core.Log.Infoln("Recieved request to handle /setNotes")
	chat := ctx.EffectiveChat
	userId := ctx.EffectiveUser.Id
	if ctx.Message.ReplyToMessage != nil && database.IsAdmin(userId) {
		userId = ctx.Message.ReplyToMessage.From.Id // switch to replied user if admin
	}

	if !database.IsMaintainer(ctx.EffectiveUser.Id) {
		_, err := b.SendMessage(chat.Id, "This is maintainers only", &gotgbot.SendMessageOpts{})
		return err
	}
	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	note := strings.Join(args, " ")
	if len(args) == 0 {
		core.Log.Infoln("No notes provided. Clearing it")
		note = ""
	}
	err := database.SetNotes(userId, note)
	if err != nil {
		_, err := b.SendMessage(chat.Id, "Something went wrong while setting notes", &gotgbot.SendMessageOpts{})
		return err
	}
	_, err = b.SendMessage(chat.Id, "Successfully updated the notes", &gotgbot.SendMessageOpts{})
	return err
}
