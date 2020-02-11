package main

import (
	// "flag"
	"fmt"

	"github.com/spectrex02/proto/net"
	"github.com/spectrex02/proto/util"
)

// var (
// name string
// addr string
// )

// func init() {
// flag.StringVar(&name, "dev", "eth0", "device name")
// flag.StringVar(&addr, "addr", "", "interface address")
// }

func main() {
	// flag.Parse()
	name := "client_veth0"
	addr := "192.168.0.3/24"
	fmt.Printf("device[%s] start\n", name)
	err := util.DisableIPForward()
	if err != nil {
		panic(err)
	}
	// ctx, err := net.Run(name, addr)
	ctx, dev, err := net.SetUp(name, addr)
	dev.DeviceInfo()
	// defer dev.Close()
	// dev.Handle()
	// go dev.Next()
	// if err != nil {
	// panic(err)
	// }
	// conn, err := net.Dial("udp", "192.168.0.2:8888")
	// conn, err := net.DialUDP(ctx, "192.168.0.2:8888")
	conn, err := net.DialUDP(ctx, "192.168.0.2:50000")
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("my udp client start")
		data := []byte("Hello from my udp client\n")
		buf := make([]byte, 50)
		for {
			l, err := conn.Write(data)
			if err != nil {
				fmt.Println("[error]: write: ", err)
			} else {
				fmt.Printf("[info]: write: %dbytes\n", l)
			}
			l, err = conn.Read(buf)
			if err != nil {
				fmt.Printf("[error]: read: %v\n", err)
			} else {
				fmt.Printf("[info]: read: %dbytes\n>> %s\n", l, string(buf))
			}
		}
	}()
	go dev.Handle()
	dev.Next()
}
