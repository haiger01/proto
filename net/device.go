package net

import (
	"fmt"
	"log"
	"syscall"

	"github.com/spectrex02/router-shakyo-go/arp"
	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/raw"
)

type Device struct {
	Name    string
	Fd      int
	Address ethernet.HardwareAddress
	MTU     int
}

type Packet struct {
	Device *Device
}

func NewDevice(name string, MTU int) (*Device, error) {
	dev, err := raw.NewPFPacket(name)
	if err != nil {
		return nil, err
	}
	addr, err := ethernet.Address(dev.Address())
	if err != nil {
		return nil, err
	}
	return &Device{
		Name:    name,
		Fd:      dev.Fd,
		Address: *addr,
		MTU:     MTU,
	}, nil
}

func (dev *Device) Read(data []byte) (int, error) {
	return syscall.Read(dev.Fd, data)
}

func (dev *Device) Write(data []byte) (int, error) {
	return syscall.Write(dev.Fd, data)
}

func (dev *Device) Close() error {
	return syscall.Close(dev.Fd)
}

func (dev *Device) DeviceInfo() {
	fmt.Println("----------device info----------")
	fmt.Println("name: ", dev.Name)
	fmt.Println("fd = ", dev.Fd)
	fmt.Println("hardware address = ", dev.Address)
}

func (dev *Device) Handle() {
	buffer := make([]byte, dev.MTU)
	for {
		_, err := dev.Read(buffer)
		if err != nil {
			log.Printf("%v error (read): %v", dev.Name, err)
		}
		etherFrame, err := ethernet.NewEthernet(buffer)
		if err != nil {
			log.Printf("%v error (read): %v", dev.Name, err)
		}
		switch etherFrame.Type() {
		case ethernet.ETHER_TYPE_ARP:
			arp, err := arp.NewARP(etherFrame.Payload())
			if err != nil {
				log.Printf("%v failed to encode to arp apcket: %v\n", dev.Name, err)
				continue
			}
			arp.Handle()
		case ethernet.ETHER_TYPE_IP:
			ip, err := ip.NewIP(etherFrame.Payload())
			if err != nil {
				log.Printf("%v failed to encode to ip packet: %v\n", dev.Name, err)
				continue
			}
			ip.Handle()
		}
	}
}
