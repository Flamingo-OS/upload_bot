package plugins

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func releaseHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	core.Log.Infoln("Recieved request to handle /release")

	taskId, err := rand.Int(rand.Reader, big.NewInt(1000000))

	if err != nil {
		core.Log.Errorln(err)
		return err
	}

	core.CancelTasks.Insert(taskId.Uint64())
	defer core.CancelTasks.Remove(taskId.Uint64())

	fmt.Println(core.CancelTasks)

	msgTxt := fmt.Sprintf("Starting release process...\nYou can cancel using `/cancel %d`", taskId.Uint64())
	_, e := b.SendMessage(chat.Id, msgTxt, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return e
}
