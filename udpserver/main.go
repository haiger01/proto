package main

import (
	"fmt"

	"github.com/spectrex02/proto/net"
	"github.com/spectrex02/proto/util"
)

func main() {
	name := "server_veth0"
	addr := "192.168.0.2/24"
	fmt.Printf("device[%s] start\n", name)
	err := util.DisableIPForward()
	if err != nil {
		panic(err)
	}
	ctx, dev, err := net.SetUp(name, addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP(ctx, "192.168.0.2:50000")
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("my udp server start")
		data := []byte("Hell client from my udp server\n")
		queue := make(chan []byte)
		go func() {
			buf := make([]byte, 48)
			for {
				l, remote, _ := conn.ReadFrom(buf)
				if l > 0 {
					select {
					case queue <- buf:
					}
				}
				// if err != nil {
				// fmt.Printf("[error] read: %v\n", err)
				// } else {
				// fmt.Printf("[info] read: %dbytes\n", l)
				// }

				l, err = conn.WriteTo(data, *remote)
				if err != nil {
					fmt.Printf("[error] write (%s:%v): %v\n", remote.IPAddress.String(), remote.Port, err)
				} else {
					fmt.Printf("[info] write to %s:%v %dbytes\n", remote.IPAddress.String(), remote.Port, l)
				}
			}
		}()
		for {
			b := <-queue
			fmt.Println("[info]received: ", string(b))
		}
	}()
	go dev.Handle()
	dev.Next()
}
