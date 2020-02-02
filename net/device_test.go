package net

import (
	"testing"

	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
)

func TestNewDevice(t *testing.T) {
	dev, err := NewDevicePFPacket("eth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
}

func TestRead(t *testing.T) {
	dev, err := NewDeviceTun("test0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
	defer dev.Close()
	buffer := make([]byte, 1500)
	for {
		_, err := dev.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}
		eth, err := ethernet.NewEthernet(buffer)
		if err != nil {
			t.Fatal(err)
		}
		eth.Header.PrintEthernetHeader()
	}
}

func TestHandle(t *testing.T) {
	dev, err := NewDeviceTun("test0", 1500)
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
	err = dev.RegisterProtocol(ip)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
	defer dev.Close()
	dev.Handle()
}
