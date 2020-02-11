package net

import (
	"fmt"
	"log"
	"syscall"

	"github.com/spectrex02/proto/ethernet"
	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/raw"
)

type PFPacket struct {
	fd                 int
	name               string
	address            ethernet.HardwareAddress
	netInfo            ip.IPSubnetMask
	registeredProtocol []LinkNetProtocol
	MTU                int
	buffer             chan *ethernet.EthernetFrame
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
	buf := make(chan *ethernet.EthernetFrame)
	return &PFPacket{
		fd:      p.Fd,
		name:    p.Name,
		address: *addr,
		MTU:     mtu,
		buffer:  buf,
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

func (p *PFPacket) NetInfo() ip.IPSubnetMask {
	return p.netInfo
}

func (p *PFPacket) IPAddress() ip.IPAddress {
	return p.netInfo.Address
}

func (p *PFPacket) Subnet() ip.IPAddress {
	return p.netInfo.Subnet
}

func (p *PFPacket) Netmask() ip.IPAddress {
	return p.netInfo.Netmask
}

func (p *PFPacket) RegisterNetInfo(info string) error {
	nInfo, err := ip.NewIPSubnetMask(info)
	if err != nil {
		return err
	}
	p.netInfo = *nInfo
	return nil
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
	// fmt.Printf("%v start handling packet", p.name)
	fmt.Println("packet handling start")
	for {
		_, err := p.Read(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", p.name, err)
		}
		frame, err := ethernet.NewEthernet(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", p.name, err)
		}
		p.buffer <- frame
	}
}

func (p *PFPacket) Next() {
	for {
		if p.registeredProtocol == nil {
			panic("next layer protocol is not registered")
		}
		// frame := <-p.buffer
		frame := <-p.Buffer()
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

// これは絶対スループット下がる
func (p *PFPacket) Buffer() chan *ethernet.EthernetFrame {
	return p.buffer
}

func (p *PFPacket) testListen() {
	fmt.Println("----------test listen----------")
	frame := <-p.buffer
	fmt.Println("src:", frame.Header.Src)
	fmt.Println("dst:", frame.Header.Dst)
}

func (p *PFPacket) RegisteredProtocol() []LinkNetProtocol {
	return p.registeredProtocol
}
