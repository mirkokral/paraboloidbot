package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"mirko/command"
	"strings"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/bot/world"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/chat/sign"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/level"
	pk "github.com/Tnze/go-mc/net/packet"
)

var client = bot.NewClient()
var msgHandler *msg.Manager
var player *basic.Player
var pl *playerlist.PlayerList
var cqtick *time.Ticker
var wrld *world.World
var chatqueue = []string{}
var connected = false
var disconnectedAtleastOnce = false

var L = command.Log{
	Info: func(message string) {
		send("&8[&binfo&8] - &b" + message)
	},
	Warn: func(message string) {
		send("&8[&6warn&8] - &e" + message)
	},
	Err: func(message string) {
		send("&8[&cerror&8] - &c" + message)
	},
}

func send(msg string) {
	if msg == "" {
		return
	}
	chatqueue = append(chatqueue, msg)
}

func main() {

	fmt.Println("parabloid v0.1")
	client.Auth.Name = "pb-1"
	player = basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{
		GameStart:  connect,
		Disconnect: disconnect,
	})
	command.Commands = append(command.Commands, command.Command{
		Name:        "help",
		Description: "The command that got you here...",
		Execute: func(l command.Log, args []string, executor string) string {
			for _, e := range command.Commands {
				send("&a" + e.Name + "&8: &7" + e.Description)
			}
			return ""
		},
	})
	pl = playerlist.New(client)
	msgHandler = msg.New(client, player, pl, handler)
	err := client.JoinServer("kaboom.fusselig.xyz:25565")
	wrld = world.NewWorld(client, player, world.EventsListener{
		LoadChunk: func(pos level.ChunkPos) error {
			return nil
		},
	})
	queueChatHandler()

	if err = client.HandleGame(); err == nil {
		panic("HandleGame never return nil")
	}
	client.Events.AddListener(
		bot.PacketHandler{
			ID:       packetid.ClientboundSystemChat,
			Priority: 12,
			F: func(p pk.Packet) error {

				return nil
			},
		},
	)

}

var handler = msg.EventsHandler{
	SystemChat: func(msg chat.Message, overlay bool) error {
		if msg.Translate == "%s %s \u203a %s" {
			onChat(msg.With[0].ClearString(), msg.With[1].ClearString(), msg.With[2].ClearString())
		}
		if msg.Translate == "[%s] %s \u203a %s" {
			onChat("["+msg.With[0].ClearString()+"]", msg.With[1].ClearString(), msg.With[2].ClearString())
		}
		return nil
	},
	DisguisedChat: func(msg chat.Message) error {
		if msg.Translate == "%s" {
			if len(msg.With[0].Extra) != 0 {
				s1 := msg.With[0].Extra[0].ClearString()
				// -- BEGIN STACKOVERFLOW STEAL ZONE --
				if last := len(s1) - 1; last >= 0 && s1[last] == ' ' {
					s1 = s1[:last]
				}
				// --  END STACKOVERFLOW STEAL ZONE  --
				onChat(
					s1,
					msg.With[0].Extra[1].ClearString(),
					msg.With[0].Extra[4].ClearString(),
				)
			}
		}
		return nil
	},
}

func onChat(rank string, username string, message string) {
	fmt.Printf("%s %s > %s\n", rank, username, message)
	if len(message) > 0 && message[0] == '`' {
		args := strings.Split(message[1:], " ")
		cmd := args[0]
		args = args[1:]
		executed := false
		for _, element := range command.Commands {
			if element.Name == cmd {
				send(element.Execute(L, args, username))
				executed = true
				break
			}
		}
		if !executed {
			L.Err("Command not found!")
		}
	}
}

func queueChatHandler() {
	ticker := time.NewTicker(200 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if connected {
					if len(chatqueue) > 0 {
						x := chatqueue[0]
						chatqueue = chatqueue[1:]

						if x[0] == '/' {
							var salt int64
							if err := binary.Read(rand.Reader, binary.BigEndian, &salt); err != nil {
							} else {

								client.Conn.WritePacket(pk.Marshal(
									packetid.ServerboundChatCommand,
									pk.String(x[1:]),
									pk.Long(time.Now().UnixMilli()),
									pk.Long(salt),
									pk.Boolean(false),
									sign.HistoryUpdate{
										Acknowledged: pk.NewFixedBitSet(20),
									},
								))
							}
						} else {
							msgHandler.SendMessage(x)
						}
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func connect() error {
	connected = true
	send("/nick &#5f5f5fp&#5b5b5ba&#565656r&#525252a&#4e4e4eb&#494949o&#454545l&#414141o&#3c3c3ci&#383838d")
	send("/rank &8|")
	if disconnectedAtleastOnce {
		send("&ereconnected to target server")
	} else {
		send("&aconnected to target server")
	}
	send("&7prefix is &c`")
	return nil
}
func disconnect(reason chat.Message) error {
	connected = false
	disconnectedAtleastOnce = true
	fmt.Println("Disconnected: " + reason.String())
	err := client.JoinServer("kaboom.fusselig.xyz:25565")
	if err = client.HandleGame(); err == nil {
		panic("HandleGame never return nil")
	}
	return nil
}
