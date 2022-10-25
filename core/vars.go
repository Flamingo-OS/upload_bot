package core

import "go.uber.org/zap"

// various publically accisible vars.
var Config *BotConfig      // contains various configs loaded
var Log *zap.SugaredLogger // logger
var CancelTasks *CancelCmds
