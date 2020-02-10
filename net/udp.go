package net

import (
	"context"
	"fmt"

	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/udp"
	"github.com/spectrex02/proto/util"
)

type UDP struct {
	// RegisteredProtocol []ApplicationProtocol
	Table *udp.Table
	Iface IP
}

func NewUDP(iface IP) *UDP {
	return &UDP{
		Table: udp.NewTable(),
		Iface: iface,
	}
}

// Type
func (u *UDP) Type() ip.IPProtocol {
	return ip.IPUDPProtocol
}

// Handle(dst []byte, protocol LinkNetProtocol data []byte) error
// dst is destination address
func (u *UDP) Handle(src, dst []byte, protocol LinkNetProtocol, data []byte) error {
	fmt.Println("------------------------------udp handling--------------------")
	datagram, err := udp.NewUDPDatagram(data)
	if err != nil {
		return err
	}
	// datagram.Header.PrintUDPHeader()
	datagram.PrintUDPDatagram()
	a, err := ip.Address(dst)
	if err != nil {
		return err
	}
	addr := &udp.Address{
		IPAddress: *a,
		Port:      datagram.Header.DestinationPort,
	}
	entry := u.Table.Search(addr)
	if entry == nil {
		return fmt.Errorf("port is unreachable(%s:%d)\n", addr.String(), addr.Port)
	}
	buf := udp.NewBuffer(addr.IPAddress, datagram.Header.DestinationPort, datagram.Data)
	select {
	case entry.Queue <- buf:
		fmt.Println(" <- in udp queue")
		return nil
	default:
		return fmt.Errorf("failed to handle")
	}
}

func (u *UDP) Write(dstAddress []byte, protocol interface{}, data []byte) (int, error) {
	return 0, fmt.Errorf("this function is dummy")
}

// Write

// udpConn fills Conn interface
type UDPConn struct {
	conn  *udp.Conn
	iface IP
}

func NewUDPConn(u *UDP, local, remote *udp.Address) (*UDPConn, error) {
	if local == nil {
		local = &udp.Address{
			IPAddress: u.Iface.Address(),
			Port:      0,
		}
	}
	entry, err := u.Table.Add(local)
	if err != nil {
		return nil, err
	}
	return &UDPConn{
		conn:  udp.NewConn(remote, entry),
		iface: u.Iface,
	}, nil
}

func (uc *UDPConn) Read(b []byte) (int, error) {
	return uc.conn.Read(b)
}

func (uc *UDPConn) ReadFrom(b []byte) (int, *udp.Address, error) {
	return uc.conn.ReadFrom(b)
}

func (uc *UDPConn) Write(b []byte) (int, error) {
	return uc.WriteTo(b, *uc.conn.Peer())
}

func (uc *UDPConn) WriteTo(b []byte, dstAddr udp.Address) (int, error) {
	datagram, err := udp.BuildUDPDatagram(uc.conn.Entry().Address.Port, dstAddr.Port, b)
	if err != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	// datagram.Header.PrintUDPHeader()
	datagram.PrintUDPDatagram()
	data, err := datagram.Serialize()
	if err != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	return uc.iface.Write(dstAddr.IPAddress.Bytes(), ip.IPUDPProtocol, data)
}

func (uc *UDPConn) Close() error {
	return nil
}

func DialUDP(ctx context.Context, addr string) (*UDPConn, error) {
	a, port, err := util.ParseAddressAndPort(addr)
	if err != nil {
		return nil, err
	}
	ipaddr, err := ip.Address(a)
	if err != nil {
		return nil, err
	}
	peer := &udp.Address{IPAddress: *ipaddr, Port: port}
	// var udp *UDP
	// switch i := ctx.Value("udp").(type) {
	// case UDP:
	// udp = &i
	// default:
	// return nil, fmt.Errorf("udp is not registered")
	// }
	u := ctx.Value("udp").(*UDP)
	return NewUDPConn(u, nil, peer)
}

func ListenUDP(ctx context.Context, addr string) (*UDPConn, error) {
	// addr is listen port
	var u *UDP
	switch i := ctx.Value("udp").(type) {
	case *UDP:
		u = i
	default:
		return nil, fmt.Errorf("udp is not registered")
	}
	a, port, err := util.ParseAddressAndPort(addr)
	if err != nil {
		return nil, err
	}
	ipaddr, err := ip.Address(a)
	peer := &udp.Address{IPAddress: *ipaddr, Port: port}

	return NewUDPConn(u, peer, nil)
}
