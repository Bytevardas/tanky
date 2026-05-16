package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("staring client")

	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 4)

	buf[0] = 9

	for {
		conn.Write(buf)
	}
}
