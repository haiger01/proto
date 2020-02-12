package net

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/spectrex02/proto/ip"
)

// tcb (Transmission Control Block)

type Address struct {
	ip.IPAddress
	Port uint16
}

type tcbTable struct {
	list  *list.List
	mutex sync.RWMutex
}

type entry struct {
	queue   chan buffer
	address *Address
}

type buffer struct {
	address ip.IPAddress
	port    uint16
	data    []byte
}

const (
	MIN_PORT_RANGE uint16 = 40000
	MAX_PORT_RANGE uint16 = 65535
)

func newBuffer(addr ip.IPAddress, port uint16, data []byte) buffer {
	return buffer{
		address: addr,
		port:    port,
		data:    data,
	}
}

func newTcbTable() *tcbTable {
	return &tcbTable{
		list:  list.New(),
		mutex: sync.RWMutex{},
	}
}

func newEntry(addr *Address) *entry {
	return &entry{
		queue:   make(chan buffer),
		address: addr,
	}
}

func NewAddress(addr ip.IPAddress, port uint16) *Address {
	return &Address{
		IPAddress: addr,
		Port:      port,
	}
}

func (t *tcbTable) add(addr *Address) (*entry, error) {
	// if port is specified, i have to look up eather this port is used or not.
	// if port is not specified, look up available port and allocate its port.
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if addr.Port == 0 {
		port := t.getAvailablePort(addr.IPAddress)
		if port == 0 {
			return nil, fmt.Errorf("there is no available port")
		}
		addr.Port = port
	} else {
		entry := t.searchUnLock(addr)
		if entry != nil {
			return nil, fmt.Errorf("entry is already exists")
		}
	}
	newEntry := &entry{
		queue:   make(chan buffer),
		address: addr,
	}
	t.list.PushBack(newEntry)
	return newEntry, nil
}

func (t *tcbTable) delete(e *entry) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for element := t.list.Front(); element != nil; element = element.Next() {
		if element.Value.(*entry) == e {
			t.list.Remove(element)
			return nil
		}
	}
	return fmt.Errorf("no such entry")
}

func (t *tcbTable) search(addr *Address) *entry {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	t.show()
	return t.searchUnLock(addr)
}

func (t *tcbTable) searchUnLock(addr *Address) *entry {
	for element := t.list.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*entry)
		if entry.address.Port == addr.Port && entry.address.IPAddress == addr.IPAddress {
			return entry
		}
	}
	return nil
}

func (t *tcbTable) getAvailablePort(addr ip.IPAddress) uint16 {
	var port uint16
	for port = MIN_PORT_RANGE; port < MAX_PORT_RANGE; port++ {
		var element *list.Element
		for element = t.list.Front(); element != nil; element = element.Next() {
			entry := element.Value.(*entry)
			if entry.address.Port == port {
				break
			}
		}
		if element == nil {
			return port
		}
	}
	return 0
}

func (t *tcbTable) show() {
	fmt.Println("----port entry table----")
	for element := t.list.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*entry)
		fmt.Printf("%s:%v\n", entry.address.String(), entry.address.Port)
	}
	fmt.Println("------------------------")
}
