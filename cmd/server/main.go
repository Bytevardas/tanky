package main

import (
	"fmt"
	"log"
	"net"
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
			fmt.Println("failed to list to the port")
			listener.Close()
		}
		fmt.Println(conn.RemoteAddr())
		buf := make([]byte, 4)
		conn.Read(buf)
		fmt.Println(buf)
	}
}
