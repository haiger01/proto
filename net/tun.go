package net

import (
	"fmt"
	"io"
	"log"

	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/raw"
)

type Tun struct {
	// file               *os.File
	file               io.ReadWriteCloser
	name               string
	address            ethernet.HardwareAddress
	protocolAddressIP  ip.IPAddress
	registeredProtocol []LinkNetProtocol
	MTU                int
	buffer             chan *ethernet.EthernetFrame
}

func NewDeviceTun(name string, mtu int) (*Tun, error) {
	t, err := raw.NewTunDevice(name)
	if err != nil {
		return nil, err
	}
	addr, err := ethernet.Address(t.Address())
	if err != nil {
		return nil, err
	}
	return &Tun{
		file:    t.File,
		name:    t.Name,
		address: *addr,
		MTU:     mtu,
		buffer:  make(chan *ethernet.EthernetFrame),
	}, nil
}

func (t *Tun) Read(data []byte) (int, error) {
	return t.file.Read(data)
}

func (t *Tun) Write(data []byte) (int, error) {
	return t.file.Write(data)
}

func (t *Tun) Close() error {
	return t.file.Close()
}

func (t *Tun) RegisterProtocol(protocol LinkNetProtocol) error {
	t.registeredProtocol = append(t.registeredProtocol, protocol)
	return nil
}

func (t *Tun) Address() ethernet.HardwareAddress {
	return t.address
}

func (t *Tun) Name() string {
	return t.name
}

func (t *Tun) ProtocolAddressIP() ip.IPAddress {
	return t.protocolAddressIP
}

func (t *Tun) RegisterIPAddress(addr ip.IPAddress) {
	t.protocolAddressIP = addr
}

func (t *Tun) DeviceInfo() {
	fmt.Println("----------device info----------")
	fmt.Println("name: ", t.name)
	fmt.Println("hardware address: ", t.address)
}

func (t *Tun) Handle() {
	buffer := make([]byte, t.MTU)
	for {
		_, err := t.Read(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", t.name, err)
		}
		frame, err := ethernet.NewEthernet(buffer)
		if err != nil {
			log.Printf("%v error (read): %v\n", t.name, err)
		}
		t.buffer <- frame
	}
}

func (t *Tun) Next() {
	for {
		if t.registeredProtocol == nil {
			panic("next leyer protocol is not registered")
		}
		frame := <-t.buffer
		for _, protocol := range t.registeredProtocol {
			if protocol.Type() == frame.Header.Type {
				err := protocol.Handle(frame.Payload())
				if err != nil {
					log.Printf("%v error: %v\n", t.name, err)
				}
			}
		}
	}
}
