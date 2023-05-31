package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"mirko/command"
	"mirko/plugin"
	"mirko/plugin/selfcare"
	"mirko/utils"
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
		// bcraw for now
		core("bcraw &#5f5f5fp&#5b5b5ba&#565656r&#525252a&#4e4e4eb&#494949o&#454545l&#414141o&#3c3c3ci&#383838d&8 \u203a [&binfo&8] - &b" + message)
	},
	Warn: func(message string) {
		// bcraw for now
		core("bcraw &#5f5f5fp&#5b5b5ba&#565656r&#525252a&#4e4e4eb&#494949o&#454545l&#414141o&#3c3c3ci&#383838d&8 \u203a [&6warn&8] - &e" + message)
	},
	Err: func(message string) {
		// bcraw for now
		core("bcraw &#5f5f5fp&#5b5b5ba&#565656r&#525252a&#4e4e4eb&#494949o&#454545l&#414141o&#3c3c3ci&#383838d&8 \u203a [&cerror&8] - &c" + message)
	},
}

func send(msg string) {
	if msg == "" {
		return
	}
	chatqueue = append(chatqueue, msg)
}

func tellraw(message chat.Message) {
	s, e := message.MarshalJSON()
	if e != nil {
		L.Err("Error while turning chat.Message into JSON: " + e.Error())
		return
	}
	core("tellraw @a " + string(s))
}

func refillCore() {
	send(fmt.Sprintf("/fill %d %d %d %d %d %d command_block{CustomName: '\"chuba core\"'} destroy", corePos.X, corePos.Y, corePos.Z, corePos.X+16, corePos.Y+16, corePos.Z+16))
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
		Execute: func(c command.Context) *chat.Message {
			s := chat.Text("")
			for i, e := range command.Commands {
				s = s.Append(chat.Text(e.Name).SetColor("green"), chat.Text(": ").SetColor("dark_gray"), chat.Text(e.Description).SetColor("gray"))
				if i < len(command.Commands)-1 {
					s = s.Append(chat.Text("\n"))
				}
			}
			tellraw(s)
			return nil
		},
	})
	client.Events.AddListener(
		bot.PacketHandler{
			ID:       0x38,
			Priority: 2387489027890,
			F: func(p pk.Packet) error {
				var (
					x               pk.Double
					y               pk.Double
					z               pk.Double
					yaw             pk.Float
					pitch           pk.Float
					flags           pk.Byte
					teleportId      pk.VarInt
					dismountViechle pk.Boolean
				)
				p.Scan(&x, &y, &z, &yaw, &pitch, &flags, &teleportId, &dismountViechle)
				corePos = utils.Vec3{
					X: int64(x),
					Y: 255,
					Z: int64(z),
				}
				refillCore()
				return nil
			},
		},
	)
	client.Events.AddGeneric(
		bot.PacketHandler{
			ID:       0x69420,
			Priority: 2387489027890,
			F: func(p pk.Packet) error {
				selfcare.OnPacket(p)
				return nil
			},
		},
	)
	wrld = world.NewWorld(client, player, world.EventsListener{
		LoadChunk: func(pos level.ChunkPos) error {
			return nil
		},
	})
	pl = playerlist.New(client)
	msgHandler = msg.New(client, player, pl, handler)

	queueChatHandler()
	err := client.JoinServer("kaboom.pw")
	println(err)
	if err = client.HandleGame(); err == nil {
		panic("HandleGame never return nil")
	}
	selfcare.Start()
}

var handler = msg.EventsHandler{
	SystemChat: func(msg chat.Message, overlay bool) error {
		if msg.Translate == "%s %s \u203a %s" {
			onChat(msg.With[0].ClearString(), msg.With[1].ClearString(), msg.With[2].ClearString())
		}
		if msg.Translate == "[%s] %s \u203a %s" {
			onChat("["+msg.With[0].ClearString()+"]", msg.With[1].ClearString(), msg.With[2].ClearString())
		}
		selfcare.OnSystemChat(msg)
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
var relativeCorePos = utils.Vec3{
	X: 0,
	Y: 0,
	Z: 0,
}

var corePos = utils.Vec3{
	X: 0,
	Y: 0,
	Z: 0,
}

func core(command string) {
	for relativeCorePos.X > 16 {
		relativeCorePos.X -= 16
		relativeCorePos.Z += 1
	}
	for relativeCorePos.Z > 16 {
		relativeCorePos.Z -= 16
		relativeCorePos.Y += 1
	}
	if relativeCorePos.Y > 16 {
		relativeCorePos.X = 0
		relativeCorePos.Y = 0
		relativeCorePos.Z = 0
	}
	client.Conn.WritePacket(
		pk.Marshal(
			packetid.ServerboundSetCommandBlock,
			pk.Long((((relativeCorePos.X+corePos.X)&0x3FFFFFF)<<38)|(((relativeCorePos.Z+corePos.Z)&0x3FFFFFF)<<12)|((relativeCorePos.Y+corePos.Y)&0xFFF)),
			pk.String(command),
			pk.VarInt(1),
			pk.Byte(0x04),
		),
	)
	relativeCorePos.X++
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
				element.Execute(command.Context{
					L:              L,
					Args:           args,
					Executor:       username,
					ExecutorPrefix: rank,
					Core:           core,
					Chat:           send,
					Client:         *client,
					Tellraw:        tellraw,
				})
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
						var salt int64
						if err := binary.Read(rand.Reader, binary.BigEndian, &salt); err != nil {
						} else {
							if x[0] == '/' {

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

							} else {
								client.Conn.WritePacket(pk.Marshal(
									packetid.ServerboundChat,
									pk.String(x),
									pk.Long(time.Now().UnixMilli()),
									pk.Long(salt),
									pk.Boolean(false),
									sign.HistoryUpdate{
										Acknowledged: pk.NewFixedBitSet(20),
									},
								))
							}
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
	if connected {
		fmt.Println("connected")
	}
	send("/nick &#5f5f5fp&#5b5b5ba&#565656r&#525252a&#4e4e4eb&#494949o&#454545l&#414141o&#3c3c3ci&#383838d")
	send("/rank &8|")
	if disconnectedAtleastOnce {
		send("&ereconnected to target server")
	} else {
		send("&econnected to target server")
	}
	send("&7prefix is &c`")
	selfcare.Inject(plugin.InjectHandler{
		Core:    core,
		Chat:    send,
		Client:  *client,
		L:       L,
		PL:      *pl,
		Tellraw: tellraw,
	})
	return nil
}
func disconnect(reason chat.Message) error {
	connected = false
	disconnectedAtleastOnce = true
	fmt.Println("Disconnected: " + reason.String())
	time.Sleep(5 * time.Second)
	err := client.JoinServer("kaboom.pw")
	println(err)
	if err = client.HandleGame(); err == nil {
		panic("HandleGame never return nil")
	}
	return nil
}
