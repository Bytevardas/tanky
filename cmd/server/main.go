package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"sync"

	"tanky/internal/protocol"
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

	b, err := protocol.ReadMessage(conn)
	if err != nil {
		fmt.Println("failed to read all bytes")
		return
	}

	if string(b) == "host" {
		room, err := generateRoomId()
		if err != nil {
			fmt.Println(err)
			return
		}
		mu.Lock()
		roomsMap[room] = make(chan net.Conn)
		mu.Unlock()
		protocol.WriteMessage(conn, []byte(room))

		joinerConn := <-roomsMap[room]
		defer conn.Close()
		defer joinerConn.Close()

		done := make(chan struct{})

		go func() {
			io.Copy(conn, joinerConn)
			close(done)
		}()
		io.Copy(joinerConn, conn)

		<-done
	}

	if string(b) == "join" && len(b) == 10 {
		code := string(b[4:])
		mu.Lock()
		ch, ok := roomsMap[code]
		mu.Unlock()
		if !ok {
			protocol.WriteMessage(conn, []byte("Room does not exist"))
			conn.Close()
			return
		}
		ch <- conn
	}
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
