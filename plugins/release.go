package plugins

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
	"github.com/Flamingo-OS/upload-bot/documents"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func releaseHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	core.Log.Infoln("Recieved request to handle /release")
	const bannerLink = "https://sourceforge.net/projects/kosp/files/banners/banner-01.png/download"
	chat := ctx.EffectiveChat
	links := ctx.Args()[1:]
	userId := ctx.EffectiveUser.Id

	if !database.IsMaintainer(userId) {
		_, err := b.SendMessage(chat.Id, "You are not a maintainer", &gotgbot.SendMessageOpts{})
		return err
	}

	// introduce our cancel tasks
	taskId, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, "Something went wrong with cancel tasks?!", &gotgbot.SendMessageOpts{})
		return err
	}
	core.CancelTasks.Insert(taskId.Uint64())
	defer core.CancelTasks.Remove(taskId.Uint64())

	// create and delete dir after completion
	var dumpPath string = fmt.Sprintf("Dumpster/%s/", taskId.String())
	core.Log.Infoln("Creating directory:", dumpPath)
	err = os.MkdirAll(dumpPath, os.ModePerm)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, "Something went wrong with creating directory", &gotgbot.SendMessageOpts{})
		return err
	}
	defer os.RemoveAll(dumpPath)

	// Actual start of the release process
	msgTxt := fmt.Sprintf("Starting release process...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m, err := b.SendMessage(chat.Id, msgTxt, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	if err != nil {
		core.Log.Error("Something went wrong")
		return err
	}
	defer b.DeleteMessage(chat.Id, m.MessageId, &gotgbot.DeleteMessageOpts{})

	var filePaths []string // stores the downloaded file paths
	msgTxt = fmt.Sprintf("Initialising download...\nThis might take a while\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	for _, url := range links {
		// download the file
		core.Log.Info("Downloading file with url:", url)
		f, err := documents.DocumentFactory(url, dumpPath)
		if err != nil {
			core.Log.Errorln(err)
			b.SendMessage(chat.Id, "Download failed. Please try again or ask darknanobot", &gotgbot.SendMessageOpts{})
			return err
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

	msgTxt = fmt.Sprintf("Finished Downloading ....\nValidating release...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})

	deviceInfo, err := parseRelease(filePaths, dumpPath)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, fmt.Sprintf("Release failed due to: %s", err), &gotgbot.SendMessageOpts{})
		return err
	}
	maintainers, err := database.GetMaintainer(deviceInfo.DeviceName)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, "Something went wrong", &gotgbot.SendMessageOpts{})
		return err
	}
	err = validateRelease(maintainers, userId)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, fmt.Sprintf("Something went wrong: %v", err), &gotgbot.SendMessageOpts{})
		return err
	}

	if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
		core.Log.Infoln("Release cancelled by user")
		b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
		return nil
	}

	// upload the files
	msgTxt = fmt.Sprintf("Uploading files...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	for _, f := range filePaths {
		uploadFolder := core.Branch + "/" + deviceInfo.DeviceName + "/" + deviceInfo.Flavour
		err := documents.OneDriveUploader(f, "flamingo"+"/"+uploadFolder)
		if err != nil {
			core.Log.Errorln(err)
			b.SendMessage(chat.Id, "Upload failed. Please try again or ask darknanobot", &gotgbot.SendMessageOpts{})
			return err
		}
		msgTxt = fmt.Sprintf("Uploaded file %s\nYou can cancel using `/cancel %d`", f, taskId.Uint64())

		if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
			core.Log.Infoln("Release cancelled by user")
			b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
			return nil
		}

		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	}
	msgTxt = fmt.Sprintf("Done uploading\nCreating and pushing ota\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})

	if err != nil {
		core.Log.Errorln("Couldn't fetch maintainers", err)
	}

	supportGroup, err := database.GetSupportGroup(maintainers[0].UserId)
	if err != nil {
		core.Log.Errorln("Couldn't fetch support group", err)
	}

	notes, err := database.GetNotes(maintainers[0].UserId)
	if err != nil {
		core.Log.Errorln("Couldn't fetch notes", err)
	}

	if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
		core.Log.Infoln("Release cancelled by user")
		b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
		return nil
	}

	msgTxt, err = CreateReleaseText(deviceInfo, maintainers, supportGroup, notes)
	if err != nil {
		core.Log.Errorln("Couldn't create release text", err)
		b.SendMessage(chat.Id, "Something went wrong while creating release", &gotgbot.SendMessageOpts{})
		return err
	}

	err = core.CreateOTACommit(deviceInfo, dumpPath)
	if err != nil {
		b.SendMessage(chat.Id, "Something went wrong while pushing ota", &gotgbot.SendMessageOpts{})
		return err
	}
	core.Log.Infoln("Sending in the final release post")

	_, err = b.SendPhoto(chat.Id, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	if err != nil {
		core.Log.Errorln("Something went wrong due to ", err)
		b.SendMessage(chat.Id, "Failed to create release post", &gotgbot.SendMessageOpts{})
		return err
	}
	_, err = b.SendPhoto(core.UpdateChannelId, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	return err
}
