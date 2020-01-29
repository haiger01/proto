package arp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ARPHeader struct {
	HardwareType HardwareType
	ProtocolType ProtocolType
	HardwareSize uint8
	ProtocolSize uint8
	OpCode       OperationCode
}

type ARP struct {
	Header                ARPHeader
	SourceHardwareAddress []byte
	SourceProtocolAddress []byte
	TargetHardwareAddress []byte
	TargetProtocolAddress []byte
}

type HardwareType uint16
type ProtocolType uint16
type OperationCode uint16

func (arp *ARP) String() {
	fmt.Println("---------------arp---------------")
	fmt.Printf("hardware type = %02x\n", arp.Header.HardwareType)
	fmt.Printf("protocol type = %02x\n", arp.Header.ProtocolType)
	fmt.Printf("hardware address size = %02x\n", arp.Header.HardwareSize)
	fmt.Printf("protocol address size = %02x\n", arp.Header.ProtocolSize)
	fmt.Printf("operation code = %02x%s\n", arp.Header.OpCode, arp.Header.OpCode.String())
	fmt.Printf("src hwaddr = %s\n", printHadrwareAddress(arp.SourceHardwareAddress))
	fmt.Printf("src protoaddr = %s\n", arp.Header.printProtocolAddress(arp.SourceProtocolAddress))
	fmt.Printf("target hwaddr = %s\n", printHadrwareAddress(arp.TargetHardwareAddress))
	fmt.Printf("target protoaddr = %s\n", arp.Header.printProtocolAddress(arp.TargetProtocolAddress))
}

func (op OperationCode) String() string {
	switch op {
	case ARP_REQUEST:
		return "(REQUEST)"
	case ARP_REPLY:
		return "(REPLY)"
	default:
		return "(UNKNOWN)"
	}
}

func NewARP(data []byte) (*ARP, error) {
	arpHeader := &ARPHeader{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, arpHeader); err != nil {
		return nil, err
	}
	arpPacket := &ARP{
		Header:                *arpHeader,
		SourceHardwareAddress: make([]byte, arpHeader.HardwareSize),
		SourceProtocolAddress: make([]byte, arpHeader.ProtocolSize),
		TargetHardwareAddress: make([]byte, arpHeader.HardwareSize),
		TargetProtocolAddress: make([]byte, arpHeader.ProtocolSize),
	}
	if err := binary.Read(buf, binary.BigEndian, arpPacket.SourceHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, arpPacket.SourceProtocolAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, arpPacket.TargetHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, arpPacket.TargetProtocolAddress); err != nil {
		return nil, err
	}
	return arpPacket, nil
}

func (arp *ARP) Serialize() ([]byte, error) {
	packet := bytes.NewBuffer(make([]byte, 28))
	if err := binary.Write(packet, binary.BigEndian, arp.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(packet, binary.BigEndian, arp.SourceHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(packet, binary.BigEndian, arp.SourceProtocolAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(packet, binary.BigEndian, arp.TargetHardwareAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(packet, binary.BigEndian, arp.TargetProtocolAddress); err != nil {
		return nil, err
	}
	return packet.Bytes(), nil
}

func (arp *ARP) Handle() {
	arp.String()
}

func printHadrwareAddress(hwaddr []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", hwaddr[0], hwaddr[1], hwaddr[2], hwaddr[3], hwaddr[4], hwaddr[5])
}

func (arphdr ARPHeader) printProtocolAddress(addr []byte) string {
	switch arphdr.ProtocolType {
	case PROTOCOL_IPv4:
		if len(addr) == 4 {
			return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
		} else {
			return "unknown address"
		}
	case PROTOCOL_IPv6:
		if len(addr) == 16 {
			return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x", addr[0], addr[1], addr[2], addr[3], addr[4], addr[5], addr[6], addr[7], addr[8], addr[9], addr[10], addr[11], addr[12], addr[13], addr[14], addr[15])
		} else {
			return "unknown address"
		}
	default:
		return "unknown address"
	}
}
