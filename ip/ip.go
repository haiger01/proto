package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  | Type of Service|          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Identification        |Flags|      Fragment Offset    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Time to Live |    Protocol   |         Header Checksum       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Source Address                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Destination Address                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

type IPHeader struct {
	VHL      VerIHL              // 8bits
	TOS      uint8               // 8bits
	Length   uint16              // 16bits
	Ident    uint16              // 16bits
	FlOffset FlagsFragmentOffset // 16bits
	TTL      uint8               // 8bits
	Protocol IPProtocol          // 8bits
	Checksum uint16              // 16bits
	Src      IPAddress           // 32bits
	Dst      IPAddress           // 32bits
	// Options  []byte
	// Padding  []byte
}
type VerIHL uint8

func (vi VerIHL) Version() uint8 {
	return uint8(vi) >> 4
}

func (vi VerIHL) IHL() uint8 {
	return uint8(vi) & 0x0F
}

type FlagsFragmentOffset uint16

func (fo FlagsFragmentOffset) Flags() uint8 {
	return uint8(fo) >> 5
}

func (fo FlagsFragmentOffset) FragmentOffset() uint16 {
	return uint16(fo) & 0x1FF
}

type IPAddress [4]byte

func (ipaddr IPAddress) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ipaddr[0], ipaddr[1], ipaddr[2], ipaddr[3])
}

type IPProtocol uint8

type IP struct {
	Header IPHeader
	// Options   []byte
	// Padding []byte
	OptionPadding []byte
	Data          []byte
}

func (iphdr *IPHeader) PrintIPHeader() {
	fmt.Println("----------ip header----------")
	fmt.Println("version = ", iphdr.VHL.Version())
	fmt.Println("ihl = ", iphdr.VHL.IHL())
	fmt.Printf("tos = %02x\n", iphdr.TOS)
	fmt.Printf("length = %02x\n", iphdr.Length)
	fmt.Printf("identifier = %02x\n", iphdr.Ident)
	fmt.Printf("flags = %02x\n", iphdr.FlOffset.Flags())
	fmt.Printf("fragment offset = %02x\n", iphdr.FlOffset.FragmentOffset())
	fmt.Printf("ttl = %02x\n", iphdr.TTL)
	fmt.Printf("protocol = %s\n", iphdr.Protocol.String())
	fmt.Printf("checksum = %02x\n", iphdr.Checksum)
	fmt.Printf("src = %s\n", iphdr.Src.String())
	fmt.Printf("dst = %s\n", iphdr.Dst.String())
}

func (ipp IPProtocol) String() string {
	switch ipp {
	case IPICMPv4Protocol:
		return "(ICMP)"
	case IPTCPProtocol:
		return "(TCP)"
	case IPUDPProtocol:
		return "(UDP)"
	default:
		return "(UNKNOWN)"
	}
}

func NewIP(data []byte) (*IP, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("ip header is too short (%d)", len(data))
	}
	header := &IPHeader{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}
	if header.VHL.Version() != uint8(4) {
		return nil, fmt.Errorf("ip version is not ipv4")
	}
	if int(header.VHL.IHL())*4 > len(data) {
		return nil, fmt.Errorf("header length is too short")
	}
	if int(header.Length) > len(data) {
		return nil, fmt.Errorf("header length is too short")
	}
	if int(header.TTL) == 0 {
		return nil, fmt.Errorf("ttl is zero")
	}

	packet := &IP{
		Header:        *header,
		OptionPadding: make([]byte, header.VHL.IHL()-20),
	}
	if err := binary.Read(buf, binary.BigEndian, packet.OptionPadding); err != nil {
		return nil, fmt.Errorf("error making option and padding: %v", err)
	}
	packet.Data = buf.Bytes()
	return packet, nil
}

func (ip *IP) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 20))
	if err := binary.Write(buf, binary.BigEndian, ip.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ip.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (ip *IP) Handle() {
	ip.Header.PrintIPHeader()
}
