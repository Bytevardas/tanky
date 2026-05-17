package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"tanky/internal/protocol"
)

var availableCommands = []string{"host", "join", "help"}

func main() {
	fmt.Println("staring client")

	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if len(os.Args) < 2 {
		log.Fatal("expecting command to be passed in: host or join <code>")
	}

	switch os.Args[1] {
	case "host":
		host := []byte("host")
		protocol.WriteMessage(conn, host)
	case "join":
		if len(os.Args) < 3 {
			log.Fatal("join command requires room id")
		}
		join := []byte("join " + os.Args[2])
		protocol.WriteMessage(conn, join)
	}

	for {
		b, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
