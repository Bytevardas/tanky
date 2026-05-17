package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

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
