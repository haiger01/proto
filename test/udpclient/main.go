package main

import (
	"fmt"
	"net"
)

func main() {
	conn, _ := net.Dial("udp", "localhost:9999")
	defer conn.Close()
	fmt.Println("send message to server.")
	conn.Write([]byte("Hello From Client."))

	fmt.Println("receive reply from server.")
	buffer := make([]byte, 1500)
	length, _ := conn.Read(buffer)
	fmt.Printf("Receive: %s \n", string(buffer[:length]))
}
