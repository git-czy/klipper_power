package protocol

import (
	"encoding/binary"
	"encoding/json"
	"klipper_power/pkg/crypto"
)

type MiioPacket struct {
	magic    uint16
	length   uint16
	unknown  uint32
	deviceID uint32
	stamp    uint32
	checksum []byte
	data     []byte
}

type MiioPacketData struct {
	RequestId  int16       `json:"id"`
	MethodName string      `json:"method"`
	Args       interface{} `json:"params"`
}

var HelloPacket = []byte{
	0x21, 0x31, 0x00, 0x20, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

func NewRequestData(method string, args interface{}) []byte {
	data, err := json.Marshal(MiioPacketData{
		RequestId:  1,
		MethodName: method,
		Args:       args,
	})
	if err != nil {
		return nil
	}
	return data
}

func NewRequestHead(token []byte, data []byte) []byte {
	buffer := make([]byte, 32)
	packet := DefaultPacket(data)
	binary.BigEndian.PutUint16(buffer[:2], packet.magic)
	binary.BigEndian.PutUint16(buffer[2:], uint16(len(packet.data)+32))
	binary.BigEndian.PutUint32(buffer[4:], packet.unknown)
	binary.BigEndian.PutUint32(buffer[8:], packet.deviceID)
	binary.BigEndian.PutUint32(buffer[12:], packet.stamp)

	checksum := crypto.Md5Byte(buffer[:16], token[:], data)
	copy(buffer[16:], checksum[:])
	return buffer
}

func DefaultPacket(data []byte) MiioPacket {
	return MiioPacket{
		magic:    0x2131,
		length:   0x0000,
		unknown:  0x0000,
		deviceID: 0x0000,
		stamp:    0x0000,
		checksum: nil,
		data:     data,
	}
}
