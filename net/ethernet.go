package net

import "github.com/spectrex02/router-shakyo-go/ethernet"

type Ethernet struct {
	Dev *Device
	// Address ethernet.HardwareAddress
}

func NewEthernet(dev *Device) *Ethernet {
	return &Ethernet{
		Dev: dev,
	}
}

func (e *Ethernet) Write(dst []byte, typ ethernet.EtherType, data []byte) (int, error) {
	d, err := ethernet.Address(dst)
	if err != nil {
		return 0, err
	}
	frame := ethernet.BuildEthernetFrame(e.Dev.Address, *d, typ, data)
	buf, err := frame.Serialize()
	if err != nil {
		return 0, err
	}
	return e.Dev.Write(buf)
}
