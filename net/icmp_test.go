package net

import (
	"testing"

	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/util"
)

func TestICMPHandle(t *testing.T) {
	dev, err := NewDevice("eth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.ProtocolAddressIP = ip.IPAddress{172, 22, 0, 3}
	link := NewEthernet(dev)
	arp := newARP(dev)
	err = dev.RegisterProtocol(arp)
	if err != nil {
		t.Fatal(err)
	}
	ip := newIP(ip.IPAddress{172, 22, 0, 3}, link)
	icmp := newICMP()
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
}
