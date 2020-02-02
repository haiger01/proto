package router

import (
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

func NewRouter(dev1, dev2 string, next [4]byte) (*Router, error) {
	device1, err := net.NewDevice(dev1, 1500)
	if err != nil {
		return nil, err
	}
	device2, err := net.NewDevice(dev2, 1500)
	if err != nil {
		return nil, err
	}
	return &Router{
		Device1:  device1,
		Device2:  device2,
		Nexthops: ip.IPAddress(next),
	}, nil
}

func (r *Router) Handle() {
	for {
		if r.EndFlag {
			return
		}

	}
}

func (r *Router) Run() {

}
