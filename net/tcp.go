package net

import (
	"context"
	"fmt"

	"github.com/spectrex02/proto/ip"
	"github.com/spectrex02/proto/tcp"
	"github.com/spectrex02/proto/util"
)

type TCP struct {
	table *tcbTable
	iface IP
}

type TCPConn struct {
	peer  *Address
	entry *entry
	tcp   *TCP
	cb    *tcp.ControlBlock
}

// this struct satisfies Listener interface
type TCPListener struct {
}

func NewTCP(table *tcbTable, iface IP) *TCP {
	return &TCP{
		table: table,
		iface: iface,
	}
}

func (t *TCP) Type() ip.IPProtocol {
	return ip.IPTCPProtocol
}

func (t *TCP) Write(dstAddr []byte, protocol interface{}, data []byte) (int, error) {
	return t.iface.Write(dstAddr, protocol, data)
}

func (t *TCP) Handle(src, dst []byte, protocol LinkNetProtocol, data []byte) error {
	fmt.Println("------------------------------tcp handling--------------------")
	packet, err := tcp.NewTCPPacket(data)
	if err != nil {
		return err
	}
	packet.PrintTCPPacket()
	for element := t.table.list.Front(); element != nil; element = element.Next() {
		e := element.Value.(*entry)
		if e.cb == nil {
			continue
		}
		resp, err := e.cb.HandleEvent(packet)
		if err != nil {
			return err
		}
		if resp != nil {
			data, err := resp.Serialize()
			if err != nil {
				return err
			}
			t.iface.Write(dst, ip.IPTCPProtocol, data)
		}
	}
	// a, err := ip.Address(dst)
	// if err != nil {
	// 	return err
	// }
	// addr := NewAddress(*a, packet.Header.DestinationPort)

	return nil
}

func NewTCPConn(t *TCP, local, remote *Address) (*TCPConn, error) {
	if local == nil {
		local = NewAddress(t.iface.Address(), 0)
	}
	e, err := t.table.add(local)
	if err != nil {
		return nil, err
	}
	return &TCPConn{
		peer:  remote,
		entry: e,
		tcp:   t,
	}, nil
}

func (tc *TCPConn) ReadFrom(b []byte) (int, *Address, error) {
	select {
	case buf := <-tc.entry.queue:
		len := copy(b, buf.data)
		peer := NewAddress(buf.address, buf.port)
		return len, peer, nil
	}
}

func (tc *TCPConn) Read(data []byte) (int, error) {
	l, _, err := tc.ReadFrom(data)
	return l, err
}

func (tc *TCPConn) Write(b []byte) (int, error) {
	return tc.WriteTo(b, *tc.peer)
}

func (tc *TCPConn) WriteTo(b []byte, dst Address) (int, error) {
	return -1, nil
}

// user interface
// func open(host, peer *Address) (*tcp.ControlBlock, error) {
// 	cb := tcp.NewControlBlock(&host.IPAddress, &peer.IPAddress, host.Port, peer.Port)
// 	return cb, nil
// }

func (t *TCP) connect(peer *Address) (*entry, error) {
	e, err := t.table.add(NewAddress(t.iface.Address(), 0))
	if err != nil {
		return nil, fmt.Errorf("invalid tcp entry")
	}
	cb := tcp.NewControlBlock(&e.address.IPAddress, &peer.IPAddress, e.address.Port, peer.Port)
	cb.Mutex.Lock()
	defer cb.Mutex.Unlock()
	if cb.HostPort == 0 {
		cb.HostPort = e.address.Port
	}
	cb.PeerAddr = &peer.IPAddress
	cb.PeerPort = peer.Port
	cb.Rcv.WND = uint32(cap(cb.Window))
	cb.Snd.ISS = tcp.Random()
	p, err := tcp.BuildTCPPacket(cb.HostPort, cb.PeerPort, cb.Snd.ISS, 0, tcp.SYN, uint16(cb.Rcv.WND), 0, nil)
	if err != nil {
		return nil, err
	}
	data, err := p.Serialize()
	if err != nil {
		return nil, err
	}
	t.Write(t.iface.Address().Bytes(), ip.IPTCPProtocol, data)
	cb.Snd.NXT = cb.Snd.ISS + 1
	cb.State = tcp.SYN_SENT
	e.cb = cb
	return e, nil
}

func (t *TCP) bind(port uint16) (*entry, error) {
	addr := NewAddress(t.iface.Address(), port)
	return t.table.add(addr)
}

func (e *entry) listen() error {
	e.cb.Mutex.Lock()
	defer e.cb.Mutex.Unlock()
	if e.cb.State != tcp.CLOSED {
		return fmt.Errorf("invalid tcb state:%s", e.cb.State.String())
	}
	if e.cb.HostPort == 0 {
		return fmt.Errorf("host port is no specified")
	}
	e.cb.State = tcp.LISTEN
	return nil
}

func (e *entry) accept() error {
	e.cb.Mutex.Lock()
	defer e.cb.Mutex.Unlock()
	if e.cb.State != tcp.LISTEN {
		return fmt.Errorf("not listen state: current state is %s", e.cb.State.String())
	}
	return nil
}

func (e *entry) send(buf []byte) (int, error) {

}

func (e *entry) recv(buf []byte) (int, error) {
	// var total int
	e.cb.Mutex.Lock()
	defer e.cb.Mutex.Unlock()
	for {
		if !e.cb.IsReadyRecv() {
			return -1, fmt.Errorf("not ready to receive")
		}
		if len(e.cb.Window)-int(e.cb.Rcv.WND) > 0 {
			break
		}
	}
	copy(buf, e.cb.Window)
}

func ListenTCP(ctx context.Context, addr *Address) (*TCPListener, error) {
	var t *TCP
	switch i := ctx.Value("tcp").(type) {
	case *TCP:
		t = i
	default:
		return nil, fmt.Errorf("tcp is not registered")
	}

}

// func (t *TCPListener) AcceptTCP() (*TCPConn, error) {

// }

func DialTCP(ctx context.Context, addr string) (*TCPConn, error) {
	a, port, err := util.ParseAddressAndPort(addr)
	if err != nil {
		return nil, err
	}
	ipaddr, err := ip.Address(a)
	if err != nil {
		return nil, err
	}
	peer := NewAddress(*ipaddr, port)
	t := ctx.Value("tcp").(*TCP)

	t.connect()
}
