package plugins

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
	"github.com/Flamingo-OS/upload-bot/documents"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func releaseHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /release")
	// sanity checks. should have a download link, should be a maintainer
	args := ctx.Args()[1:]
	if len(args) == 0 {
		_, err := b.SendMessage(chat.Id, "Please provide a valid URL", &gotgbot.SendMessageOpts{})
		return err
	}
	userId := ctx.EffectiveUser.Id
	if ctx.Message.ReplyToMessage != nil && database.IsAdmin(userId) {
		userId = ctx.Message.ReplyToMessage.From.Id
	}
	if database.IsMaintainer(userId) {
		_, err := b.SendMessage(chat.Id, "You are not a maintainer", &gotgbot.SendMessageOpts{})
		return err
	}

	taskId, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		core.Log.Errorln(err)
		return err
	}
	core.CancelTasks.Insert(taskId.Uint64())
	defer core.CancelTasks.Remove(taskId.Uint64())

	msgTxt := fmt.Sprintf("Starting release process...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m, e := b.SendMessage(chat.Id, msgTxt, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	if e != nil {
		core.Log.Error("Something went wrong")
	}
	f, e := documents.DocumentFactory(args[0])
	if e != nil {
		core.Log.Errorln(e)
		m.EditText(b, "Download failed. Please try again or ask darknanobot", &gotgbot.EditMessageTextOpts{})
		return e
	}
	msgTxt = fmt.Sprintf("Downloaded file at %v", f)
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{})

	return e
}
