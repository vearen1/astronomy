package protocol

import (
	"fmt"
	"io"
)

func ReadVarInt(r io.Reader) (int32, error) {
	var value int32
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return 0, err
		}
		value |= int32(b[0]&0x7F) << shift
		if b[0]&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, fmt.Errorf("VarInt too big")
		}
	}
	return value, nil
}

func WriteVarInt(w io.Writer, value int32) error {
	for {
		b := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			b |= 0x80
		}
		if _, err := w.Write([]byte{b}); err != nil {
			return err
		}
		if value == 0 {
			break
		}
	}
	return nil
}

func ReadString(r io.Reader) (string, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteString(w io.Writer, s string) error {
	if err := WriteVarInt(w, int32(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}
