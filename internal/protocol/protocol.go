package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	CommandHost byte = 0x01
	CommandJoin byte = 0x02
)

type Command struct {
	Type    byte
	Payload []byte
}

func EncodeCommand(cmd Command) []byte {
	b := make([]byte, 1+len(cmd.Payload))
	b[0] = cmd.Type
	copy(b[1:], cmd.Payload)

	return b
}

func DecodeCommand(b []byte) (Command, error) {
	if len(b) == 0 {
		return Command{}, errors.New("empty command message")
	}

	return Command{Type: b[0], Payload: b[1:]}, nil
}

func ReadMessage(conn net.Conn) ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return []byte{}, err
	}
	length := binary.BigEndian.Uint32(header)
	msg := make([]byte, length)
	_, err = io.ReadFull(conn, msg)
	if err != nil {
		return []byte{}, err
	}

	return msg, nil
}

func WriteMessage(conn net.Conn, b []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(b)))

	buf := net.Buffers{header, b}
	_, err := buf.WriteTo(conn)
	return err
}
