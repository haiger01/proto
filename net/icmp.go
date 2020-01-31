package net

import (
	"fmt"

	"github.com/spectrex02/router-shakyo-go/icmp"
	"github.com/spectrex02/router-shakyo-go/ip"
)

type ICMP struct {
	IPProtocolType ip.IPProtocol
}

func newICMP() *ICMP {
	return &ICMP{
		IPProtocolType: ip.IPICMPv4Protocol,
	}
}

func (ic *ICMP) Type() ip.IPProtocol {
	return ic.IPProtocolType
}

func (ic *ICMP) Handle(dst []byte, protocol LinkNetProtocol, data []byte) error {
	packet, err := icmp.NewICMPPacket(data)
	if err != nil {
		return fmt.Errorf("encoding error: %v", err)
	}
	// packet.Header.PrintICMPHeader()
	packet.Header.PrintICMPHeader()
	reply, err := packet.Handle()
	if err != nil {
		return err
	}
	if reply == nil {
		return nil
	}
	buf, err := reply.Serialize()
	if err != nil {
		return err
	}
	_, err = protocol.Write(dst, ip.IPICMPv4Protocol, buf)
	if err != nil {
		return err
	}
	return nil
}

func IPCMEchoRequest(dst []byte, protocol NetTransProtocol) error {
	packet, err := icmp.BuildICMPPacket(icmp.Echo, icmp.EchoRequestCode, make([]byte, 0))
	if err != nil {
		return fmt.Errorf("failed to build packet")
	}
	buf, err := packet.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize packet")
	}
	_, err = protocol.Write(dst, ip.IPICMPv4Protocol, buf)
	return err
}

func (ic *ICMP) Write(dst []byte, protocol interface{}, data []byte) (int, error) {
	return 0, fmt.Errorf("this function is not implemented")
}
