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
	if !database.IsMaintainer(userId) {
		_, err := b.SendMessage(chat.Id, "You are not a maintainer", &gotgbot.SendMessageOpts{})
		return err
	}

	// introduce our cancel tasks
	taskId, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		core.Log.Errorln(err)
		return err
	}
	core.CancelTasks.Insert(taskId.Uint64())
	defer core.CancelTasks.Remove(taskId.Uint64())

	// Actual start of the release process
	msgTxt := fmt.Sprintf("Starting release process...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m, e := b.SendMessage(chat.Id, msgTxt, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	if e != nil {
		core.Log.Error("Something went wrong")
		return e
	}

	var filePaths []string // stores the downloaded file paths
	for _, url := range args {
		// download the file
		msgTxt = fmt.Sprintf("Downloading file %v...\nThis might take a while\nYou can cancel using `/cancel %d`", url, taskId.Uint64())
		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
		f, e := documents.DocumentFactory(url)
		if e != nil {
			core.Log.Errorln(e)
			b.SendMessage(chat.Id, "Download failed. Please try again or ask darknanobot", &gotgbot.SendMessageOpts{})
			return e
		}
		filePaths = append(filePaths, f)
		msgTxt = fmt.Sprintf("Downloaded file to %s. Have downloaded %v files.\nYou can cancel using `/cancel %d`", f, len(filePaths), taskId.Uint64())
		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})

		if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
			core.Log.Infoln("Release cancelled by user")
			b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
			return nil
		}
	}

	// upload the files
	// TODO: upload it to a specific dir instead of some random dir
	msgTxt = fmt.Sprintf("Uploading files...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	for _, f := range filePaths {
		err := documents.OneDriveUploader(f, "YAY1")
		if err != nil {
			core.Log.Errorln(err)
			b.SendMessage(chat.Id, "Upload failed. Please try again or ask darknanobot", &gotgbot.SendMessageOpts{})
			return err
		}
		msgTxt = fmt.Sprintf("Uploaded file %s\nYou can cancel using `/cancel %d`", f, taskId.Uint64())
		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	}

	return e
}
