package packet

import (
	"astronomy/astronomy/internal/protocol"
	"bytes"
	"encoding/binary"
	"io"
)

type PingRequestPacket struct {
	Payload int64
}

func (p *PingRequestPacket) ID() PacketID {
	return PingRequestPacketID
}

func (p *PingRequestPacket) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := protocol.WriteVarInt(buf, int32(p.ID())); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.Payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PingRequestPacket) Decode(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &p.Payload)
}

type PingResponsePacket struct {
	Payload int64
}

func (p *PingResponsePacket) ID() PacketID {
	return PingResponsePacketID
}

func (p *PingResponsePacket) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := protocol.WriteVarInt(buf, int32(p.ID())); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.Payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PingResponsePacket) Decode(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &p.Payload)
}
