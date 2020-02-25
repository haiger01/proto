package net

import "testing"

func TestTCPHandler(t *testing.T) {
	_, dev, err := SetUp("server_veth0", "192.168.0.2/24")
	if err != nil {
		panic(err)
	}
	go dev.Handle()
	dev.Next()
}
