package plugins

import (
	"strconv"
	"strings"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func cancelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	args := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]
	core.Log.Infoln("Recieved request to handle /cancel")

	if len(args) == 0 {
		_, e := b.SendMessage(chat.Id, "Please provide a task id", &gotgbot.SendMessageOpts{})
		return e
	}

	taskId, err := strconv.ParseUint(args[0], 0, 64)
	if err != nil {
		_, e := b.SendMessage(chat.Id, "Please provide a valid task id", &gotgbot.SendMessageOpts{})
		return e
	}

	core.CancelTasks.Cancel(taskId)

	_, e := b.SendMessage(chat.Id, "Queued request to cancel the task", &gotgbot.SendMessageOpts{})
	return e
}
