package net

import (
	"github.com/spectrex02/proto/ethernet"
	"github.com/spectrex02/proto/ip"
)

type Device interface {
	Read(data []byte) (int, error)
	Write(data []byte) (int, error)
	Close() error
	Address() ethernet.HardwareAddress
	// ProtocolAddressIP() ip.IPAddress
	Name() string
	NetInfo() ip.IPSubnetMask
	IPAddress() ip.IPAddress
	Subnet() ip.IPAddress
	Netmask() ip.IPAddress
	RegisterNetInfo(info string) error
	RegisterProtocol(protocol LinkNetProtocol) error
	RegisteredProtocol() []LinkNetProtocol
	DeviceInfo()
	Handle()
	Next()
	Buffer() chan *ethernet.EthernetFrame
}

type Buffer struct {
	Data []byte
	Src  ethernet.HardwareAddress
	Dst  ethernet.HardwareAddress
}
