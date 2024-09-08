package packet

import (
	"astronomy/astronomy/internal/protocol"
	"bytes"
	"encoding/binary"
	"io"
)

type HandshakePacket struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (p *HandshakePacket) ID() PacketID {
	return HandshakePacketID
}

func (p *HandshakePacket) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := protocol.WriteVarInt(buf, int32(p.ID())); err != nil {
		return nil, err
	}
	if err := protocol.WriteVarInt(buf, p.ProtocolVersion); err != nil {
		return nil, err
	}
	if err := protocol.WriteString(buf, p.ServerAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.ServerPort); err != nil {
		return nil, err
	}
	if err := protocol.WriteVarInt(buf, p.NextState); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *HandshakePacket) Decode(r io.Reader) error {
	var err error
	if p.ProtocolVersion, err = protocol.ReadVarInt(r); err != nil {
		return err
	}
	if p.ServerAddress, err = protocol.ReadString(r); err != nil {
		return err
	}
	if err = binary.Read(r, binary.BigEndian, &p.ServerPort); err != nil {
		return err
	}
	if p.NextState, err = protocol.ReadVarInt(r); err != nil {
		return err
	}
	return nil
}
