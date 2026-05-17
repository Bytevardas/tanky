package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var availableCommands = []string{"host", "join", "help"}

func main() {
	fmt.Println("staring client")

	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[1]
	if strings.ToLower(command) != "host" {
		host := []byte("host")
		conn.Write(host)
	}
}
