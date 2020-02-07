package arp

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	HardwareAddress []byte
	ProtocolAddress []byte
	ProtocolType    ProtocolType
	TimeStamp       time.Time
}

type ARPTable struct {
	Entrys []*Entry
	Mutex  sync.RWMutex
}

func NewEntry(hwaddr, protoaddr []byte, typ ProtocolType) *Entry {
	return &Entry{
		HardwareAddress: hwaddr,
		ProtocolAddress: protoaddr,
		ProtocolType:    typ,
		TimeStamp:       time.Now(),
	}
}

func NewARPTable() *ARPTable {
	return &ARPTable{
		Entrys: make([]*Entry, 0, 1024),
		Mutex:  sync.RWMutex{},
	}
}

func (at *ARPTable) Search(protoaddr []byte) *Entry {
	at.Mutex.RLock()
	defer at.Mutex.RUnlock()
	for _, e := range at.Entrys {
		if bytes.Equal(e.ProtocolAddress, protoaddr) {
			// fmt.Printf("[info] found entry (%s)\n", printProtocolAddress(protoaddr))
			return e
		}
	}
	// fmt.Println("[info] not found the entry")
	return nil
}

func (at *ARPTable) Insert(hwaddr, protoaddr []byte, typ ProtocolType) error {
	at.Mutex.Lock()
	defer at.Mutex.Unlock()
	if len(at.Entrys) > 1023 {
		return fmt.Errorf("arp table is full")
	}
	for _, e := range at.Entrys {
		if bytes.Equal(e.ProtocolAddress, protoaddr) {
			return fmt.Errorf("this address pair is already entried")
		}
	}
	e := NewEntry(hwaddr, protoaddr, typ)
	at.Entrys = append(at.Entrys, e)
	// fmt.Printf(">>>>>>>>>>>>>>>>>>>insert into arp table [%v -> %v]\n", hwaddr, protoaddr)
	return nil
}

func (at *ARPTable) Update(hwaddr, protoaddr []byte) (bool, error) {
	at.Mutex.Lock()
	defer at.Mutex.Unlock()
	for _, e := range at.Entrys {
		if bytes.Equal(e.ProtocolAddress, protoaddr) {
			e.HardwareAddress = hwaddr
			e.TimeStamp = time.Now()
			return true, nil
		}
	}
	return false, nil
}

func (at *ARPTable) Show() {
	fmt.Println("---------------arp table---------------")
	for _, e := range at.Entrys {
		fmt.Printf("hwaddr= %s\n", printHadrwareAddress(e.HardwareAddress))
		fmt.Printf("protoaddr=%s\n", printProtocolAddress(e.ProtocolAddress))
		fmt.Printf("time=%v\n", e.TimeStamp)
	}
	fmt.Println("---------------------------------------")
}
