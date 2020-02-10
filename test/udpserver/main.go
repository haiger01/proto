package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server is Running at localhost:8888")
	conn, _ := net.ListenPacket("udp", "192.168.0.2:8888")
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		// 通信読込 + 接続相手アドレス情報が受取
		length, remoteAddr, _ := conn.ReadFrom(buffer)
		fmt.Printf("Received from %v: %v\n", remoteAddr, string(buffer[:length]))
		i, err := conn.WriteTo([]byte("Hello, World !"), remoteAddr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[info] %dbytes to %v\n", i, remoteAddr)
	}
}
