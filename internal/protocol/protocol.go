package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	CommandHost byte = 0x01
	CommandJoin byte = 0x02
)

func EncodeCommand(t byte, payload []byte) []byte {
	b := make([]byte, 1+len(payload))
	b[0] = t
	copy(b[1:], payload)
	return b
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
