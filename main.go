package main

import (
	"snap_chat_server/cmd"
	"snap_chat_server/config"
	"snap_chat_server/logger"
)

func main() {
	config.Init()
	logger.Register()
	cmd.Execute()
}
