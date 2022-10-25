package plugins

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Main(d *ext.Dispatcher) {
	d.AddHandler(handlers.NewCommand("add", addMaintainerHandler))
	d.AddHandler(handlers.NewCommand("cancel", cancelHandler))
	d.AddHandler(handlers.NewCommand("demote", demoteAdminHandler))
	d.AddHandler(handlers.NewCommand("dropDevice", removeDevicesHandler))
	d.AddHandler(handlers.NewCommand("remove", removeMaintainerHandler))
	d.AddHandler(handlers.NewCommand("help", helpHandler))
	d.AddHandler(handlers.NewCommand("ping", pingHandler))
	d.AddHandler(handlers.NewCommand("promote", promoteAdminHandler))
	d.AddHandler(handlers.NewCommand("release", releaseHandler))
	d.AddHandler(handlers.NewCommand("start", startHelper))
}
