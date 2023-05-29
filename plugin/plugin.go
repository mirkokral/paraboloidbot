package plugin

import (
	"mirko/command"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/playerlist"
)

type InjectHandler struct {
	Core   func(command string)
	Chat   func(command string)
	Client bot.Client
	L      command.Log
	PL     playerlist.PlayerList
}
