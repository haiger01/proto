package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/spectrex02/router-shakyo-go/util"
)

type ICMPHeader struct {
	Type ICMPType
	// Code messageCode
	Code     uint8
	Checksum uint16
}

type ICMPPacket struct {
	Header ICMPHeader
	Data   []byte
}

type ICMPType uint8
type messageCode uint8

func newICMPHeader(typ ICMPType, code uint8) *ICMPHeader {
	return &ICMPHeader{
		Type:     typ,
		Code:     code,
		Checksum: uint16(0),
	}
}

func NewICMPPacket(data []byte) (*ICMPPacket, error) {
	header := &ICMPHeader{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	return &ICMPPacket{
		Header: *header,
		Data:   buf.Bytes(),
	}, nil
}

func (icmp *ICMPPacket) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 8))
	if err := binary.Write(buf, binary.BigEndian, icmp.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, icmp.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (icmphdr *ICMPHeader) PrintICMPHeader() {
	fmt.Println("----------icmp header----------")
	fmt.Printf("type = %v\n", icmphdr.Type)
	fmt.Printf("code = %02x\n", icmphdr.Code)
	fmt.Printf("checksum = %02x\n", icmphdr.Checksum)
}

func (icmp ICMPPacket) PrintICMP() {
	icmp.Header.PrintICMPHeader()
	fmt.Printf("data = %v\n", icmp.Data)
}

func (typ ICMPType) String() string {
	switch typ {
	case Echo:
		return fmt.Sprintf("Echo Request")
	case EchoReply:
		return fmt.Sprintf("Echo Reply")
	case DestinationUnreachable:
		return fmt.Sprintf("Destination Unreachable")
	case TimeExceeded:
		return fmt.Sprintf("Time Exceeded")
	case Redirect:
		return fmt.Sprintf("Redirect")
	default:
		return fmt.Sprintf("unknown")
	}
}

func (icmp *ICMPPacket) Handle() (*ICMPPacket, error) {
	switch icmp.Header.Type {
	case Echo:
		return BuildICMPPacket(EchoReply, EchoReplyCode, nil)
	case EchoReply:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupportted ICMP type: %v", icmp.Header.Type.String())
	}
}

func BuildICMPPacket(typ ICMPType, code uint8, data []byte) (*ICMPPacket, error) {
	header := newICMPHeader(typ, code)
	packet := &ICMPPacket{
		Header: *header,
		Data:   data,
	}
	buf, err := packet.Serialize()
	if err != nil {
		return nil, err
	}
	sum := util.Checksum(buf, len(buf), 0)
	packet.Header.Checksum = sum
	return packet, nil
}
