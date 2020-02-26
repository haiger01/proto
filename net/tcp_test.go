package net

import (
	"testing"

	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/util"
)

func TestTCPHandler(t *testing.T) {
	err := util.DisableIPForward()
	if err != nil {
		t.Fatal(err)
	}
	ctx, dev, err := SetUp("server_veth0", "192.168.0.2/24")
	if err != nil {
		panic(err)
	}
	addr := NewAddress(ip.NewIPAddress([]byte{192, 168, 0, 2}), 8888)
	ListenTCP(ctx, addr)
	go dev.Handle()
	go dev.Next()
	for {
	}
}
