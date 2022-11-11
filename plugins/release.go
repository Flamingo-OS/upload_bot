package plugins

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/Flamingo-OS/upload-bot/database"
	"github.com/Flamingo-OS/upload-bot/documents"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const bannerLink = "https://sourceforge.net/projects/kosp/files/banners/banner-01.png/download"

// validates the release is indeed a flamingo OS file
// also creates and pushes OTA file
func validateRelease(fileNames []string) (core.DeviceInfo, error) {
	deviceInfo, fullOtaFile, incrementalOtaFile, err := core.ParseDeviceInfo(fileNames)
	if err != nil {
		return core.DeviceInfo{}, err
	}
	core.Log.Info("Parsed device info:", deviceInfo)
	err = core.CreateOTACommit(deviceInfo, fullOtaFile, incrementalOtaFile)
	if err != nil {
		return deviceInfo, err
	}
	return deviceInfo, err
}

func CreateReleaseText(deviceInfo core.DeviceInfo, urls []string, maintainers []database.Maintainers, deviceSupport string) (string, error) {
	formatTime := time.Now().Format("2006-01-02")
	buildTime, _ := time.Parse("20060102", deviceInfo.BuildDate)
	av, err := strconv.ParseFloat(deviceInfo.Version, 64)
	if err != nil {
		core.Log.Error("Error while parsing android version: %s", err)
		return "", err
	}
	androidVersion := int(av) + 11
	changeLogUrl := fmt.Sprintf("https://raw.githubusercontent.com/FlamingoOS-Devices/OTA/main/%s/%s/changelog_%s", deviceInfo.DeviceName, deviceInfo.Flavour, strings.ReplaceAll(formatTime, "-", "_"))
	msgTxt := fmt.Sprintf(`FlamingoOS %s | Android %d | OFFICIAL | %s
	
		Maintainers: `, deviceInfo.Version, androidVersion, deviceInfo.Flavour)

	for _, maintainer := range maintainers {
		msgTxt += fmt.Sprintf(` [%s](tg://user?id=%d) `, maintainer.MaintainerName, maintainer.UserId)
	}

	msgTxt += fmt.Sprintf(`
		
		Device: %s
		Date: %s

		Downloads`, deviceInfo.DeviceName, buildTime.Format("02-01-2006"))

	for _, url := range urls {
		if strings.Contains(url, "-full") {
			msgTxt += fmt.Sprintf(`
		-	[Full](%s)	`, url)
		} else if strings.Contains(url, "-incremental") {
			msgTxt += fmt.Sprintf(`
		-	[Incremental](%s)	`, url)
		} else if strings.Contains(url, "-fastboot") {
			msgTxt += fmt.Sprintf(`
		-	[Fastboot](%s)	`, url)
		} else if strings.Contains(url, "-boot") {
			msgTxt += fmt.Sprintf(`
		-	[Boot](%s)	`, url)
		}
	}
	msgTxt += fmt.Sprintf(`

		[Changelog](%s)
	
		Support Groups
		-	[Common](https://t.me/flamingo_common)
		-	[Updates](https://t.me/flamingo_updates)
	`, changeLogUrl)

	if deviceSupport != "" {
		msgTxt += fmt.Sprintf(`
		-	[Device Support](%s)`, deviceSupport)
	}

	return msgTxt, nil
}

func releaseHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /release")

	// create and delete dir after completion
	os.Mkdir(core.DumpPath, 0755)
	defer os.RemoveAll(core.DumpPath)

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
		b.SendMessage(chat.Id, "Something went wrong with cancel tasks?!", &gotgbot.SendMessageOpts{})
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
	defer b.DeleteMessage(chat.Id, m.MessageId, &gotgbot.DeleteMessageOpts{})

	var filePaths []string // stores the downloaded file paths
	msgTxt = fmt.Sprintf("Initialising download...\nThis might take a while\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	for _, url := range args {
		// download the file
		core.Log.Info("Downloading file with url:", url)
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

	msgTxt = fmt.Sprintf("Finished Downloading ....\nValidating release...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})

	deviceInfo, err := validateRelease(filePaths)
	if err != nil {
		core.Log.Errorln(err)
		b.SendMessage(chat.Id, fmt.Sprintf("Release failed due to: %s", err), &gotgbot.SendMessageOpts{})
		return err
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
		m.EditText(b, msgTxt, &gotgbot.EditMessageTextOpts{ParseMode: "markdown"})
	}

	msgTxt = "Uploaded file to `"
	for _, url := range uploadUrls {
		msgTxt += url + "`, `"
	}
	msgTxt = strings.Trim(msgTxt, ", `")
	core.Log.Info(msgTxt)

	maintainers, e := database.GetMaintainer(deviceInfo.DeviceName)

	if e != nil {
		core.Log.Errorln("Couldn't fetch maintainers", e)
	}

	supportGroup, e := database.GetSupportGroup(maintainers[0].UserId)
	if e != nil {
		core.Log.Errorln("Couldn't fetch support group", e)
	}

	msgTxt, e = CreateReleaseText(deviceInfo, uploadUrls, maintainers, supportGroup)
	if e != nil {
		core.Log.Errorln("Couldn't create release text", e)
		b.SendMessage(chat.Id, "Something went wrong while creating release", &gotgbot.SendMessageOpts{})
	}
	_, _ = b.SendPhoto(chat.Id, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	_, e = b.SendPhoto(core.UpdateChannelId, bannerLink, &gotgbot.SendPhotoOpts{
		Caption:   msgTxt,
		ParseMode: "Markdown",
	})
	return e
}
