package packet

import (
	"io"
)

type PacketID byte

const (
	HandshakePacketID      PacketID = 0x00
	StatusRequestPacketID  PacketID = 0x00
	StatusResponsePacketID PacketID = 0x00
	PingRequestPacketID    PacketID = 0x01
	PingResponsePacketID   PacketID = 0x01
)

type Packet interface {
	ID() PacketID
	Encode() ([]byte, error)
	Decode(r io.Reader) error
}
