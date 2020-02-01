package net

import (
	"fmt"
	"log"
	"syscall"

	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/raw"
)

type PFPacket struct {
	fd                 int
	name               string
	address            ethernet.HardwareAddress
	protocolAddressIP  ip.IPAddress
	registeredProtocol []LinkNetProtocol
	MTU                int
}

func NewDevicePFPacket(name string, mtu int) (*PFPacket, error) {
	p, err := raw.NewPFPacket(name)
	if err != nil {
		return nil, err
	}
	addr, err := ethernet.Address(p.Address())
	if err != nil {
		return nil, err
	}
	return &PFPacket{
		fd:      p.Fd,
		name:    p.Name,
		address: *addr,
		MTU:     mtu,
	}, nil
}

func (p *PFPacket) Name() string {
	return p.name
}

func (p *PFPacket) Read(data []byte) (int, error) {
	return syscall.Read(p.fd, data)
}

func (p *PFPacket) Write(data []byte) (int, error) {
	return syscall.Write(p.fd, data)
}

func (p *PFPacket) Close() error {
	return syscall.Close(p.fd)
}

func (p *PFPacket) Address() ethernet.HardwareAddress {
	return p.address
}

func (p *PFPacket) ProtocolAddressIP() ip.IPAddress {
	return p.protocolAddressIP
}

func (p *PFPacket) RegisterIPAddress(addr ip.IPAddress) {
	p.protocolAddressIP = addr
}

func (p *PFPacket) DeviceInfo() {
	fmt.Println("----------device info----------")
	fmt.Println("name: ", p.name)
	fmt.Println("fd = ", p.fd)
	fmt.Println("hardware address = ", p.address)
}

func (p *PFPacket) RegisterProtocol(protocol LinkNetProtocol) error {
	p.registeredProtocol = append(p.registeredProtocol, protocol)
	return nil
}

func (p *PFPacket) Handle() {
	buffer := make([]byte, p.MTU)
	for {
		_, err := p.Read(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", p.name, err)
		}
		frame, err := ethernet.NewEthernet(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", p.name, err)
		}
		for _, protocol := range p.registeredProtocol {
			if protocol.Type() == frame.Header.Type {
				err := protocol.Handle(frame.Payload())
				if err != nil {
					log.Printf("%v error: %v\n", p.name, err)
				}
			}
		}
	}
}
