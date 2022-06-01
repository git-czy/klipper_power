package plug

import (
	"klipper_power/pkg/protocol"
)

func OffRequest() {

}

func DiscoverRequest(token [16]byte) []byte {
	//deviceKey := crypto.DeviceKeyFromToken(token)
	helloPacket := protocol.NewHelloPacket()
	requestHead := protocol.NewRequestHead(*helloPacket, token)
	return requestHead
}
