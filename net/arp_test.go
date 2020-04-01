package net

import (
	"testing"
)

func TestARPHandler(t *testing.T) {
	dev, err := NewDevicePFPacket("host1_veth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.RegisterNetInfo("192.168.0.2/24")
	arp := NewARP(dev)
	err = dev.RegisterProtocol(arp)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
	defer dev.Close()
	go dev.Handle()
	dev.Next()
}
