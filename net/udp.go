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
	table *tcbTable
	iface IP
}

func NewUDP(table *tcbTable, iface IP) *UDP {
	return &UDP{
		table: table,
		iface: iface,
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
	addr := NewAddress(*a, datagram.Header.DestinationPort)
	entry := u.table.search(addr)
	if entry == nil {
		return fmt.Errorf("port is unreachable(%s:%d)\n", addr.String(), addr.Port)
	}
	buf := newBuffer(addr.IPAddress, datagram.Header.DestinationPort, datagram.Data)
	select {
	case entry.queue <- buf:
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
	peer  *Address
	entry *entry
	iface *UDP
}

func NewUDPConn(u *UDP, local, remote *Address) (*UDPConn, error) {
	if local == nil {
		local = NewAddress(u.iface.Address(), 0)
	}
	e, err := u.table.add(local)
	if err != nil {
		return nil, err
	}
	return &UDPConn{
		peer:  remote,
		entry: e,
		iface: u,
	}, nil
}

func (uc *UDPConn) Read(b []byte) (int, error) {
	l, _, err := uc.ReadFrom(b)
	return l, err
}

func (uc *UDPConn) ReadFrom(b []byte) (int, *Address, error) {
	select {
	case buf := <-uc.entry.queue:
		len := copy(b, buf.data)
		peer := &Address{
			IPAddress: buf.address,
			Port:      buf.port,
		}
		return len, peer, nil
	}
}

func (uc *UDPConn) Write(b []byte) (int, error) {
	return uc.WriteTo(b, *uc.peer)
}

func (uc *UDPConn) WriteTo(b []byte, dstAddr Address) (int, error) {
	datagram, err := udp.BuildUDPDatagram(uc.entry.address.Port, dstAddr.Port, b)
	if err != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	// datagram.Header.PrintUDPHeader()
	datagram.PrintUDPDatagram()
	data, err := datagram.Serialize()
	if err != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	return uc.iface.iface.Write(dstAddr.IPAddress.Bytes(), ip.IPUDPProtocol, data)
}

func (uc *UDPConn) Close() error {
	return uc.iface.table.delete(uc.entry)
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
	peer := NewAddress(*ipaddr, port)
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
	peer := NewAddress(*ipaddr, port)
	return NewUDPConn(u, peer, nil)
}
