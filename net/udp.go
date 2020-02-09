package net

import (
	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/udp"
)

type UDP struct {
	// RegisteredProtocol []ApplicationProtocol
	Table              udp.Table
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
func (u *UDP) Handle(dst []byte, protocol LinkNetProtocol, data []byte) error {
	datagram, err := udp.NewUDPDatagram(data)
	if err != nil {
		return err
	}
	a, err := ip.Address(dst)
	if err != nil {
		return err
	}
	addr := &upd.Address{
		IPAddress: a,
		Port: datagram.Header.DestinationPort,
	}
	entry := u.Table.Search(addr)
	if entry == nil {
		return fmt.Errorf("port(%v) is unreachable")
	}
	buf := udp.Buffer{
		address: addr.IPAddress,
		port: datagram.Header.DestinationPort,
		data: datagram.Data,
	}
	select {
		case entry.Queue <- buf:
			return nil
		default:
			return fmt.Errorf("failed to handle")
	}
}

// Write

// udpConn fills Conn interface
type UDPConn struct {
	conn udp.Conn
	iface IP
}

func (uc *UDPConn) Read(b []byte) (int, error) {
	return uc.conn.Read(b)
}

func (uc *UDPConn) Write(b []byte) (int, error) {
	return uc.WriteTo(b, uc.conn.Port())
}

func (uc *UDPConn) WriteTo(b []byte, dstAddr udp.Address) (int, error) {
	datagram, err := udp.BuildUDPDatagram(uc.conn.Entry().Address.Port, dstAddr.Port, b)
	if err != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	data, err := frame.Serialize()
	if er != nil {
		return -1, fmt.Errorf("failed to write: %v\n", err)
	}
	return uc.IP.Write(dstAddr.IPAddress, ip.IPUDPProtocol, data)
}

func (uc *UDPConn) Close() error {
	return nil
}