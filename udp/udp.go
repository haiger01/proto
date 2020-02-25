package udp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/util"
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
	fmt.Printf("source port = %v\n", udphdr.SourcePort)
	fmt.Printf("destination port = %v\n", udphdr.DestinationPort)
	fmt.Printf("length = %x\n", udphdr.Length)
	fmt.Printf("checksum = %x\n", udphdr.Checksum)
	fmt.Println("-----------------------------")
}

func NewUDPDatagram(data []byte) (*UDPDatagram, error) {
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

func (udpd *UDPDatagram) CalculateChecksum() error {
	data, err := udpd.Serialize()
	if err != nil {
		return err
	}
	sum := util.Checksum2(data, int(udpd.Header.Length), 0)
	udpd.Header.Checksum = sum
	return nil
}

func BuildUDPDatagram(src, dst uint16, data []byte) (*UDPDatagram, error) {
	header := UDPHeader{
		SourcePort:      src,
		DestinationPort: dst,
		Length:          uint16(8 + len(data)),
		Checksum:        0,
	}
	datagram := &UDPDatagram{
		Header: header,
		Data:   data,
	}
	err := datagram.CalculateChecksum()
	if err != nil {
		return nil, err
	}
	return datagram, nil
}

func (udpd *UDPDatagram) PrintUDPDatagram() {
	udpd.Header.PrintUDPHeader()
	fmt.Println(string(udpd.Data))
	fmt.Println("-----------------------------")
}
