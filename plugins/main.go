package plugins

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Main(d *ext.Dispatcher) {
	d.AddHandler(handlers.NewCommand("start", startHelper))
	d.AddHandler(handlers.NewCommand("help", helpHandler))
}
