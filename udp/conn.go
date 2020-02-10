package udp

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/spectrex02/proto/ip"
)

// port lange
// 40000 ~ 65535

type Conn struct {
	peer  *Address // remote address
	entry *Entry   // Address field in entry fileld has my address and port info
}

func NewConn(peer *Address, entry *Entry) *Conn {
	return &Conn{
		peer:  peer,
		entry: entry,
	}
}

type Entry struct {
	Queue   chan Buffer
	Address *Address
}

type Buffer struct {
	address ip.IPAddress
	port    uint16
	data    []byte
}

func NewBuffer(addr ip.IPAddress, port uint16, data []byte) Buffer {
	return Buffer{
		address: addr,
		port:    port,
		data:    data,
	}
}

type Table struct {
	List  *list.List
	Mutex sync.RWMutex
}

func NewTable() *Table {
	return &Table{
		List:  list.New(),
		Mutex: sync.RWMutex{},
	}
}

func (t *Table) Add(addr *Address) (*Entry, error) {
	// if port is specified, i have to look up eather this port is used or not.
	// if port is not specified, look up available port and allocate its port.
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	if addr.Port == 0 {
		port := t.getAvailablePort(addr.IPAddress)
		if port == 0 {
			return nil, fmt.Errorf("there is no available port")
		}
		addr.Port = port
	} else {
		entry := t.search(addr)
		if entry != nil {
			return nil, fmt.Errorf("entry is already exists")
		}
	}
	newEntry := &Entry{
		Queue:   make(chan Buffer),
		Address: addr,
	}
	t.List.PushBack(newEntry)
	return newEntry, nil
}

func (t *Table) Delete(entry *Entry) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	for element := t.List.Front(); element != nil; element = element.Next() {
		if element.Value.(*Entry) == entry {
			t.List.Remove(element)
			return nil
		}
	}
	return fmt.Errorf("no such entry")
}

func (t *Table) Search(addr *Address) *Entry {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	return t.search(addr)
}

func (t *Table) search(addr *Address) *Entry {
	for element := t.List.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*Entry)
		if entry.Address.Port == addr.Port && entry.Address.IPAddress == addr.IPAddress {
			return entry
		}
	}
	return nil
}

func (t *Table) getAvailablePort(addr ip.IPAddress) uint16 {
	var port uint16
	for port = MIN_PORT_RANGE; port < MAX_PORT_RANGE; port++ {
		var element *list.Element
		for element = t.List.Front(); element != nil; element = element.Next() {
			entry := element.Value.(*Entry)
			if entry.Address.Port == port {
				break
			}
		}
		if element == nil {
			return port
		}
	}
	return 0
}

func (conn *Conn) Port() uint16 {
	return conn.peer.Port
}

func (conn *Conn) Address() ip.IPAddress {
	return conn.peer.IPAddress
}

func (conn *Conn) Peer() *Address {
	return conn.peer
}

func (conn *Conn) Entry() *Entry {
	return conn.entry
}
func (conn *Conn) Close() error {
	return nil
}

func (conn *Conn) Read(buf []byte) (int, error) {
	if conn.peer == nil {
		return -1, fmt.Errorf("invalid connection")
	}
	len, _, err := conn.ReadFrom(buf)
	if err != nil {
		return -1, fmt.Errorf("failed to read")
	}
	return len, nil
}

func (conn *Conn) ReadFrom(buf []byte) (int, *Address, error) {
	select {
	case b := <-conn.entry.Queue:
		len := copy(buf, b.data)
		peer := &Address{
			IPAddress: b.address,
			Port:      b.port,
		}
		return len, peer, nil
	}
}

// func (conn *Conn) Write(buf []byte) (int, error) {
//
// }

// func (conn *Conn) WriteTo(buf []byte, addr *Address) (int, error) {
//
// }
