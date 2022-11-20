package plugins

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

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
	if ctx.Message.ReplyToMessage != nil && database.IsAdmin(userId) {
		userId = ctx.Message.ReplyToMessage.From.Id // switch to replied user if admin
	}

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

	deviceInfo, err := validateRelease(filePaths, dumpPath)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, fmt.Sprintf("Release failed due to: %s", err), &gotgbot.SendMessageOpts{})
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
	var uploadUrls []string
	for _, f := range filePaths {
		uploadFolder := core.Branch + "/" + deviceInfo.DeviceName + "/" + deviceInfo.Flavour
		err := documents.OneDriveUploader(f, "flamingo"+"/"+uploadFolder)
		if err != nil {
			core.Log.Errorln(err)
			b.SendMessage(chat.Id, "Upload failed. Please try again or ask darknanobot", &gotgbot.SendMessageOpts{})
			return err
		}
		fileName := strings.Split(f, "/")[len(strings.Split(f, "/"))-1]
		uploadUrl := core.BaseUrl + uploadFolder + "/" + fileName
		uploadUrls = append(uploadUrls, uploadUrl)
		msgTxt = fmt.Sprintf("Uploaded file %s\nYou can cancel using `/cancel %d`", f, taskId.Uint64())

		if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
			core.Log.Infoln("Release cancelled by user")
			b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
			return nil
		}

		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	}

	msgTxt = "Uploaded file to `"
	for _, url := range uploadUrls {
		msgTxt += url + "`, `"
	}
	msgTxt = strings.Trim(msgTxt, ", `")
	core.Log.Info(msgTxt)

	maintainers, err := database.GetMaintainer(deviceInfo.DeviceName)

	if err != nil {
		core.Log.Errorln("Couldn't fetch maintainers", err)
	}

	supportGroup, err := database.GetSupportGroup(maintainers[0].UserId)
	if err != nil {
		core.Log.Errorln("Couldn't fetch support group", err)
	}

	notes, err := database.GetNotes(userId)
	if err != nil {
		core.Log.Errorln("Couldn't fetch notes", err)
	}

	if core.CancelTasks.GetCancelStatus(taskId.Uint64()) {
		core.Log.Infoln("Release cancelled by user")
		b.SendMessage(chat.Id, "Release cancelled by user", &gotgbot.SendMessageOpts{})
		return nil
	}

	msgTxt, err = CreateReleaseText(deviceInfo, uploadUrls, maintainers, supportGroup, notes)
	if err != nil {
		core.Log.Errorln("Couldn't create release text", err)
		b.SendMessage(chat.Id, "Something went wrong while creating release", &gotgbot.SendMessageOpts{})
		return err
	}
	_, _ = b.SendPhoto(chat.Id, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	_, err = b.SendPhoto(core.UpdateChannelId, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	return err
}
