package net

import (
	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/udp"
)

type UDP struct {
	RegisteredProtocol []ApplicationProtocol
	Table              udp.Table
}

func NewUDP() *UDP {

}

// Type
func (u *UDP) Type() ip.IPProtocol {
	return ip.IPUDPProtocol
}

// Handle

// Write
