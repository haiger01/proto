package net

import (
	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
)

type Device interface {
	Read(data []byte) (int, error)
	Write(data []byte) (int, error)
	Close() error
	Address() ethernet.HardwareAddress
	ProtocolAddressIP() ip.IPAddress
	Name() string
	RegisterIPAddress(addr ip.IPAddress)
	RegisterProtocol(protocol LinkNetProtocol) error
	DeviceInfo()
	Handle()
}
