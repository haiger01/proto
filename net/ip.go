package net

import (
	"fmt"

	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
)

type IP struct {
	RegisteredProtocol []NetTransProtocol
	NetInfo            ip.IPSubnetMask
	HardwareType       ethernet.EtherType
	Link               *Ethernet
}

func NewIP(ninfo string, link *Ethernet) *IP {
	info, err := ip.NewIPSubnetMask(ninfo)
	if err != nil {
		panic(err)
	}
	return &IP{
		NetInfo:      *info,
		HardwareType: ethernet.ETHER_TYPE_IP,
		Link:         link,
	}
}

func (i *IP) Type() ethernet.EtherType {
	return i.HardwareType
}

func (i *IP) RegisterProtocol(protocol NetTransProtocol) error {
	i.RegisteredProtocol = append(i.RegisteredProtocol, protocol)
	return nil
}

func (i *IP) Handle(data []byte) error {
	packet, err := ip.NewIPPacket(data)
	if err != nil {
		return fmt.Errorf("failed to create ip packet")
	}
	// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	// packet.Header.PrintIPHeader()
	if i.RegisteredProtocol == nil {
		return fmt.Errorf("next protocols is not registered")
	}
	for _, protocol := range i.RegisteredProtocol {
		err := protocol.Handle(packet.Header.Src.Bytes(), i, packet.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *IP) Write(dst []byte, protocol interface{}, data []byte) (int, error) {
	d, err := ip.Address(dst)
	if err != nil {
		return 0, err
	}
	packet := ip.BuildIPPacket(i.NetInfo.Address, *d, protocol.(ip.IPProtocol), data)
	buf, err := packet.Serialize()
	if err != nil {
		return 0, err
	}
	return i.Link.Write(i.Link.Dev.Address().Bytes(), ethernet.ETHER_TYPE_IP, buf)
}
