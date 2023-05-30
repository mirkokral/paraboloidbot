package selfcare

import (
	"mirko/plugin"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

var p plugin.InjectHandler
var Skin string
var OP bool
var CurrentGM Gamemode = Creative

func Inject(h plugin.InjectHandler) {
	p = h
	skinChange(p.Client.Name)
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
			p.L.Info("&7OP Change: &cDeOP")
			OP = false
		case 25:
			p.L.Info("&7OP Change: &cDeOP")
			OP = false
		case 26:
			p.L.Info("&7OP Change: &cDeOP")
			OP = false
		case 27:
			p.L.Info("&7OP Change: &cDeOP")
			OP = false
		case 28:
			p.L.Info("&7OP Change: &aOP")
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
				p.L.Info("&7Bot's gamemode: &cSurvival")
				CurrentGM = Survival
			case 1:
				p.L.Info("&7Bot's gamemode: &aCreative")
				CurrentGM = Creative
			case 2:
				p.L.Info("&7Bot's gamemode: &eAdventure")
				CurrentGM = Adventure
			case 3:
				p.L.Info("&7Bot's gamemode: &dSpectator")
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
