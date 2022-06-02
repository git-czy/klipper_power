package plug

import (
	"klipper_power/pkg/crypto"
	"klipper_power/pkg/protocol"
)

const DefaultMethod = "set_properties"

type DID struct {
	Did   string `json:"did"`
	Siid  uint8  `json:"siid"`
	Piid  uint8  `json:"piid"`
	Value bool   `json:"value"`
}

func NewPowerRequest(token []byte, t string) []byte {
	var args interface{}
	switch t {
	case "on":
		args = onArgs()
	case "off":
		args = offArgs()
	case "hello":
		return protocol.HelloPacket
	}
	requestData := protocol.NewRequestData(DefaultMethod, args)
	requestHead := protocol.NewRequestHead(token, requestData)

	deviceKey := crypto.DeviceKeyFromToken(token)
	encryptData := deviceKey.Encrypt(requestData)
	request := append(requestHead, encryptData...)

	return request
}

func offArgs() []DID {
	return []DID{{Did: "MYDID", Siid: 2, Piid: 1, Value: false}}
}

func onArgs() []DID {
	return []DID{{Did: "MYDID", Siid: 2, Piid: 1, Value: true}}
}
