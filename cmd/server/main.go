package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"sync"
)

var (
	roomsMap = make(map[string]chan net.Conn)
	mu       sync.Mutex
)

func main() {
	fmt.Println("starting to listen to 8080")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("failed to list to the port")
	}
	defer listener.Close()

	fmt.Println("waiting for connections")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed accept connection")
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("inside the go routine")
	b, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("failed to read all bytes")
		conn.Close()
	}

	fmt.Println(string(b))
}

const chars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func generateRoomId() (string, error) {
	b := make([]byte, 6)
	charLength := big.NewInt(int64(len(chars)))

	for i := range b {
		index, err := rand.Int(rand.Reader, charLength)
		if err != nil {
			return "", err
		}
		b[i] = chars[index.Int64()]
	}

	return string(b), nil
}
