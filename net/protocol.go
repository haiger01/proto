package net

import (
	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/ip"
)

type LinkNetProtocol interface {
	Type() ethernet.EtherType
	Handle(data []byte) error
	Write(dst []byte, protocol interface{}, data []byte) (int, error)
}

type NetTransProtocol interface {
	Type() ip.IPProtocol
	Handle(dst []byte, protocol LinkNetProtocol, data []byte) error
	Write(dstAddress []byte, protocol interface{}, data []byte) (int, error)
}

// type LinkProtocol interface {

// }
