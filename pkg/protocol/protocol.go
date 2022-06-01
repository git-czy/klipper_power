package protocol

import (
	"encoding/binary"
)

type MiioPacket struct {
	magic    uint16
	length   uint16
	unknown  uint32
	deviceID uint32
	stamp    uint32
	checksum int
	data     []byte
}

type MiioPacketData struct {
}

func NewHelloPacket() *MiioPacket {
	return &MiioPacket{
		magic:    0x2131,
		length:   0x0020,
		unknown:  0xffffffff,
		deviceID: 0xffffffff,
		stamp:    0xffffffff,
		checksum: 0xffffffffffffffffffffffffffffff,
		data:     nil,
	}
}

func NewRequestHead(packet MiioPacket, token [16]byte) []byte {
	var buffer []byte
	binary.BigEndian.PutUint16(buffer[:], packet.magic)
	if packet.data == nil && packet.length != 0x00 {
		binary.BigEndian.PutUint16(buffer[2:], packet.length)
	} else {
		binary.BigEndian.PutUint16(buffer[2:], uint16(len(packet.data)+32))
	}
	binary.BigEndian.PutUint32(buffer[4:], packet.unknown)
	binary.BigEndian.PutUint32(buffer[8:], packet.deviceID)
	binary.BigEndian.PutUint32(buffer[12:], packet.stamp)
	copy(buffer[16:], token[:])
	return buffer
}
