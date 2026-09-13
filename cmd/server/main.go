package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"strconv"
	"sync"

	"tanky/internal/protocol"
)

var (
	roomsMap = make(map[string]chan net.Conn)
	mu       sync.Mutex
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("listening on", listener.Addr())
	for _, a := range reachableAddresses(listener.Addr()) {
		fmt.Println("clients on this network can connect with -server", a)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed accept connection")
			continue
		}
		go handleConnection(conn)
	}
}

func reachableAddresses(listen net.Addr) []string {
	tcpAddr, ok := listen.(*net.TCPAddr)
	if !ok {
		return nil
	}
	if !tcpAddr.IP.IsUnspecified() {
		return []string{tcpAddr.String()}
	}
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	port := strconv.Itoa(tcpAddr.Port)
	var addrs []string
	for _, a := range interfaceAddrs {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
			continue
		}
		addrs = append(addrs, net.JoinHostPort(ipNet.IP.String(), port))
	}
	return addrs
}

func handleConnection(conn net.Conn) {
	fmt.Println("inside the go routine")

	b, err := protocol.ReadMessage(conn)
	if err != nil {
		fmt.Println("failed to read all bytes")
		return
	}
	if len(b) == 0 {
		fmt.Println("received empty message")
		return
	}

	switch b[0] {
	case protocol.CommandHost:
		room, err := generateRoomId()
		if err != nil {
			fmt.Println(err)
			return
		}
		mu.Lock()
		roomsMap[room] = make(chan net.Conn)
		mu.Unlock()
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgRoomCode, []byte(room)))

		joinerConn := <-roomsMap[room]
		defer conn.Close()
		defer joinerConn.Close()

		mu.Lock()
		delete(roomsMap, room)
		mu.Unlock()

		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgStart, nil))
		protocol.WriteMessage(joinerConn, protocol.EncodeCommand(protocol.MsgStart, nil))

		remaining := make(chan net.Conn, 2)
		go relay(conn, joinerConn, remaining)
		go relay(joinerConn, conn, remaining)
		protocol.WriteMessage(<-remaining, protocol.EncodeCommand(protocol.MsgPeerLeft, nil))

	case protocol.CommandJoin:
		if len(b) < 2 {
			protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgError, []byte("missing code")))
			conn.Close()
			return
		}
		mu.Lock()
		ch, ok := roomsMap[string(b[1:])]
		mu.Unlock()
		if !ok {
			protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgError, []byte("Room does not exist")))
			conn.Close()
			return
		}

		ch <- conn

	default:
		protocol.WriteMessage(conn, protocol.EncodeCommand(protocol.MsgError, []byte("unknown command")))
		conn.Close()
		return
	}
}

func relay(dst, src net.Conn, remaining chan<- net.Conn) {
	io.Copy(dst, src)
	remaining <- dst
}

func generateRoomId() (string, error) {
	b := make([]byte, protocol.RoomCodeLength)
	charLength := big.NewInt(int64(len(protocol.RoomCodeChars)))

	for i := range b {
		index, err := rand.Int(rand.Reader, charLength)
		if err != nil {
			return "", err
		}
		b[i] = protocol.RoomCodeChars[index.Int64()]
	}

	return string(b), nil
}
