package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/spectrex02/proto/ip"
)

type UDPDatagram struct {
	Header UDPHeader
	Data   []byte
}

type UDPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}

type Address struct {
	ip.IPAddress
	Port uint16
}

func (udphdr UDPHeader) PrintUDPHeader() {
	fmt.Println("---------udp header----------")
	fmt.Printf("source header = %v\n", udphdr.SourcePort)
	fmt.Printf("destination header = %v\n", udphdr.DestinationPort)
	fmt.Printf("length = %x\n", udphdr.Length)
	fmt.Printf("checksum = %x\n", udphdr.Checksum)
	fmt.Println("-----------------------------")
}

func NewUDDatagram(data []byte) (*UDPDatagram, error) {
	header := &UDPHeader{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("header encoding error: %v\n", err)
	}
	datagram := &UDPDatagram{
		Header: *header,
	}
	datagram.Data = data[8:]
	return datagram, nil
}

func (udpd *UDPDatagram) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, int(udpd.Header.Length)))
	if err := binary.Write(buf, binary.BigEndian, udpd.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, udpd.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BuildUDPDatagram(src, dst uint16, data []byte) (*UDPDatagram, error) {
	header := UDPHeader{
		SourcePort:      src,
		DestinationPort: dst,
		Length:          uint16(8 + len(data)),
		Checksum:        0,
	}
	return &UDPDatagram{
		Header: header,
		Data:   data,
	}, nil
}
