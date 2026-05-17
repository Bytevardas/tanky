package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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

	command := os.Args[1]
	if strings.ToLower(command) == "host" {
		host := []byte("host")
		protocol.WriteMessage(conn, host)
	}

	for {
		b, err := protocol.ReadMessage(conn)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}
}
