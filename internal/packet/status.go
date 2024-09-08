package packet

import (
	"astronomy/astronomy/internal/protocol"
	"bytes"
	"io"
)

type StatusResponsePacket struct {
	Response string
}

func (p *StatusResponsePacket) ID() PacketID {
	return StatusResponsePacketID
}

func (p *StatusResponsePacket) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := protocol.WriteVarInt(buf, int32(p.ID())); err != nil {
		return nil, err
	}
	if err := protocol.WriteString(buf, p.Response); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *StatusResponsePacket) Decode(r io.Reader) error {
	var err error
	if p.Response, err = protocol.ReadString(r); err != nil {
		return err
	}
	return nil
}
