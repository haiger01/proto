package net

import (
	"testing"

	"github.com/spectrex02/router-shakyo-go/ethernet"
)

func TestNewDevice(t *testing.T) {
	dev, err := NewDevice("eth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
}

func TestRead(t *testing.T) {
	dev, err := NewDevice("eth0", 1500)
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
	dev, err := NewDevice("eth0", 1500)
	if err != nil {
		t.Fatal(err)
	}
	dev.DeviceInfo()
	defer dev.Close()
	dev.Handle()
}
