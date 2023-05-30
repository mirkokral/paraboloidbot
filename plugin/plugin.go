package plugin

import (
	"mirko/command"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/chat"
)

type InjectHandler struct {
	Core    func(command string)
	Chat    func(command string)
	Tellraw func(msg chat.Message)
	Client  bot.Client
	L       command.Log
	PL      playerlist.PlayerList
}
