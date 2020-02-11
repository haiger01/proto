package net

import (
	"context"
	"fmt"

	"github.com/spectrex02/proto/ethernet"
	"github.com/spectrex02/proto/ip"
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

func (i *IP) WithValue(ctx context.Context) context.Context {
	newCtx := ctx
	for _, protocol := range i.RegisteredProtocol {
		typ := protocol.Type().String()
		fmt.Println("add contex with value:", typ)
		newCtx = context.WithValue(newCtx, typ, protocol)
	}
	return newCtx
}

func (i *IP) Handle(data []byte) error {
	packet, err := ip.NewIPPacket(data)
	if err != nil {
		return fmt.Errorf("failed to create ip packet")
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	packet.Header.PrintIPHeader()
	if i.RegisteredProtocol == nil {
		return fmt.Errorf("next protocols is not registered")
	}
	for _, protocol := range i.RegisteredProtocol {
		if protocol.Type() == packet.Header.Protocol {
			err := protocol.Handle(packet.Header.Src.Bytes(), packet.Header.Dst.Bytes(), i, packet.Data)
			if err != nil {
				return err
			}
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
	err = packet.ReCalculateChecksum()
	if err != nil {
		return 0, err
	}
	buf, err := packet.Serialize()
	if err != nil {
		return 0, err
	}
	packet.Header.PrintIPHeader()
	return i.Link.Write(i.Link.Dev.Address().Bytes(), ethernet.ETHER_TYPE_IP, buf)
}

func (i *IP) Address() ip.IPAddress {
	return i.NetInfo.Address
}
