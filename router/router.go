package router

import (
	"fmt"

	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/net"
)

type Router struct {
	Device1  net.Device
	Device2  net.Device
	Nexthops ip.IPAddress
	EndFlag  bool
}

func setUp() {

}

func NewRouter(name1, name2, addr1, addr2, next string) (*Router, error) {
	dev1, err := setUpDevice(name1, addr1)
	if err != nil {
		return nil, fmt.Errorf("failed to set up %v: %v\n", name1, err)
	}
	dev2, err := setUpDevice(name2, addr2)
	if err != nil {
		return nil, fmt.Errorf("failed to set up %v: %v\n", name1, err)
	}
	nextAddr, err := ip.StrintToIPAddress(next)
	if err != nil {
		return nil, err
	}
	return &Router{
		Device1:  dev1,
		Device2:  dev2,
		Nexthops: *nextAddr,
	}, nil
}

func setUpDevice(name string, addr string) (net.Device, error) {
	address, err := ip.StrintToIPAddress(addr)
	if err != nil {
		return nil, err
	}
	dev, err := net.NewDevicePFPacket(name, 1500)
	if err != nil {
		return nil, err
	}
	// link := net.NewEthernet(dev)
	arp := net.NewARP(dev)
	// ip := net.NewIP(*address, link)
	// icmp := net.NewICMP()
	err = dev.RegisterProtocol(arp)
	if err != nil {
		return nil, err
	}
	// err = ip.RegisterProtocol(icmp)
	// if err != nil {
	// return nil, err
	// }
	dev.RegisterIPAddress(*address)
	return dev, nil
}

func (r *Router) Handle() {
	for {

	}
}

func (r *Router) Run() {

}
