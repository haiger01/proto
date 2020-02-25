package net

import (
	"context"
)

func Run(name, addr string) (context.Context, error) {
	ctx, dev, err := SetUp(name, addr)
	if err != nil {
		return nil, err
	}
	dev.DeviceInfo()
	defer dev.Close()
	go dev.Handle()
	go dev.Next()
	return ctx, nil
}

func SetUp(name, addr string) (context.Context, Device, error) {
	// ctx := context.Background()
	dev, err := NewDevicePFPacket(name, 1500)
	if err != nil {
		return nil, nil, err
	}
	link := NewEthernet(dev)
	arp := NewARP(dev)
	err = dev.RegisterProtocol(arp)
	if err != nil {
		return nil, nil, err
	}
	err = dev.RegisterNetInfo(addr)
	if err != nil {
		return nil, nil, err
	}
	ip := NewIP(addr, link)
	icmp := NewICMP()
	ip.RegisterProtocol(icmp)
	table := newTcbTable()
	udp := NewUDP(table, *ip)
	ip.RegisterProtocol(udp)
	tcp := NewTCP(table, *ip)
	ip.RegisterProtocol(tcp)
	err = dev.RegisterProtocol(ip)
	if err != nil {
		return nil, nil, err
	}
	ctx := ip.WithValue(context.Background())
	return ctx, dev, nil
}
