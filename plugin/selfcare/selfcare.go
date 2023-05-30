package selfcare

import (
	"mirko/plugin"
	"time"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

var p plugin.InjectHandler
var Skin string = "pb-1"
var OP bool = true
var CurrentGM Gamemode = Creative
var e = false

func Start() {

}

func Inject(h plugin.InjectHandler) {
	p = h
	skinChange(p.Client.Name)
	if !e {
		ticker := time.NewTicker(200 * time.Millisecond)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					if !OP {
						h.Chat("/op @s[type=player]")
					}
					if CurrentGM != Creative {
						h.Chat("/gamemode creative")
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}
	e = true
}

func OnPacket(pk packet.Packet) {
	if pk.ID == int32(packetid.ClientboundEntityEvent) {
		var (
			entityId packet.Int
			status   packet.Byte
		)
		pk.Scan(&entityId, &status)
		switch status {
		case 24:
			OP = false
		case 25:
			OP = false
		case 26:
			OP = false
		case 27:
			OP = false
		case 28:
			OP = true
		}
	}
	if pk.ID == int32(packetid.ClientboundGameEvent) {
		var (
			event packet.UnsignedByte
			value packet.Float
		)
		pk.Scan(&event, &value)
		if event == 3 {
			switch value {
			case 0:
				CurrentGM = Survival
			case 1:
				CurrentGM = Creative
			case 2:
				CurrentGM = Adventure
			case 3:
				CurrentGM = Spectator
			}
		}
	}

}
func OnSystemChat(msg chat.Message) {
	clear := msg.ClearString()
	if clear == "Successfully removed your skin" {
		skinChange(p.Client.Name)
	}
	if len(clear) > len("Successfully set your skin to ") && clear[:30] == "Successfully set your skin to " {
		skinChange(clear[30:][:len(clear)-32])
	}

}

func skinChange(newSkin string) {
	p.L.Info("&7Skin change: &a" + newSkin)
}
