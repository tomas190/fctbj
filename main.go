package main

import (
	"fctbj/conf"
	"fctbj/game"
	"fctbj/gate"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"log"
)

func main() {
	log.Println("start:", )
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
	)
}
