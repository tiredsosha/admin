package protocols

import (
	"github.com/linde12/gowol"
	"github.com/tiredsosha/admin/tools/logger"
)

func SendWOL(mac string) {

	packet, err := gowol.NewMagicPacket(mac)
	if err == nil {
		packet.Send("255.255.255.255") // send to broadcast
		logger.Info.Printf("mac sent - %q\n", mac)
	} else {
		logger.Error.Printf("error sending mac %v", err)
	}
}
