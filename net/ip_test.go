package net

import (
	"testing"

	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/util"
)

func TestIPHandle(t *testing.T) {
	dev, err := NewDevicePFPacket("eth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.RegisterIPAddress(ip.IPAddress{172, 22, 0, 3})
	link := NewEthernet(dev)
	arp := NewARP(dev)
	err = dev.RegisterProtocol(arp)
	if err != nil {
		t.Fatal(err)
	}
	ip := NewIP(ip.IPAddress{172, 22, 0, 3}, link)
	icmp := NewICMP()
	ip.RegisterProtocol(icmp)
	err = dev.RegisterProtocol(ip)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
	defer dev.Close()
	err = util.DisableIPForward()
	if err != nil {
		t.Fatal(err)
	}
	dev.Handle()
	dev.Next()
}
