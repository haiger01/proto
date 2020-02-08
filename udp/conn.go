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
	entry *Entry
}

type Entry struct {
	Queue   chan buffer
	Address *Address
}

type buffer struct {
	address ip.IPAddress
	port    uint16
	data    []byte
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
		Queue:   make(chan buffer),
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
