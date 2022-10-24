package plugins

import (
	"errors"

	"github.com/Flamingo-OS/upload-bot/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func startHelper(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat

	core.Log.Infoln("Just testing errors out")

	b.SendMessage(chat.Id, "Hey there! I am alive", &gotgbot.SendMessageOpts{})
	return errors.New("hey that works")
}
