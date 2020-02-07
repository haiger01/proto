package router

import (
	"fmt"

	"github.com/spectrex02/router-shakyo-go/arp"
	"github.com/spectrex02/router-shakyo-go/ethernet"
	"github.com/spectrex02/router-shakyo-go/icmp"
	"github.com/spectrex02/router-shakyo-go/ip"
	"github.com/spectrex02/router-shakyo-go/net"
)

type Router struct {
	Devices    map[bool]net.Device
	ArpTables  map[bool]*arp.ARPTable
	Nexthops   ip.IPAddress
	ReadyQueue chan buffer
	EndFlag    bool
}

type buffer struct {
	deviceNo bool
	address  ip.IPAddress
	frame    *ethernet.EthernetFrame
}

func NewRouter(name1, name2, addr1, addr2, next string) (*Router, error) {
	dev1, table1, err := setUpDevice(name1, addr1)
	if err != nil {
		return nil, fmt.Errorf("failed to set up %v: %v\n", name1, err)
	}
	dev2, table2, err := setUpDevice(name2, addr2)
	if err != nil {
		return nil, fmt.Errorf("failed to set up %v: %v\n", name1, err)
	}
	nextAddr, err := ip.StringToIPAddress(next)
	if err != nil {
		return nil, err
	}

	return &Router{
		Devices:    map[bool]net.Device{true: dev1, false: dev2},
		ArpTables:  map[bool]*arp.ARPTable{true: table1, false: table2},
		Nexthops:   *nextAddr,
		ReadyQueue: make(chan buffer, 1),
	}, nil
}

func setUpDevice(name string, addr string) (net.Device, *arp.ARPTable, error) {
	dev, err := net.NewDevicePFPacket(name, 1500)
	if err != nil {
		return nil, nil, err
	}
	// link := net.NewEthernet(dev)
	arp := net.NewARP(dev)
	// ip := net.NewIP(*address, link)
	// icmp := net.NewICMP()
	err = dev.RegisterProtocol(arp)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%v register arp protocol\n", dev.Name())
	// err = ip.RegisterProtocol(icmp)
	// if err != nil {
	// return nil, err
	// }
	err = dev.RegisterNetInfo(addr)
	if err != nil {
		return nil, nil, err
	}
	return dev, arp.Table, nil
}

func (r *Router) Handle() {
	for {
		select {
		case frame := <-r.Devices[true].Buffer():
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			frame.Header.PrintEthernetHeader()
			err := r.handle(frame, true)
			if err != nil {
				fmt.Printf("[info]%v error: %v", r.Devices[true].Name(), err)
			}
		case frame := <-r.Devices[false].Buffer():
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			frame.Header.PrintEthernetHeader()
			err := r.handle(frame, false)
			if err != nil {
				fmt.Printf("[info]%v error: %v", r.Devices[false].Name(), err)
			}
		}
	}
}

func (r *Router) arpHandle(frame *ethernet.EthernetFrame, devNo bool) {
	r.Devices[devNo].Next()
}

func (r *Router) handle(frame *ethernet.EthernetFrame, devNo bool) error {
	// this ethernet frame should be ip datagram
	// fmt.Printf("--------------------%v--------------------\n", r.Devices[devNo].Name())
	// check is data correct ip packet
	packet, err := ip.NewIPPacket(frame.Payload())
	if err != nil {
		return fmt.Errorf("not ip datagram: %v\n", err)
	}
	// packet.Header.PrintIPHeader()
	if packet.Header.Protocol == ip.IPICMPv4Protocol {
		icmpPacket, err := icmp.NewICMPPacket(packet.Data)
		if err != nil {
			return err
		}
		if icmpPacket.Header.Type == icmp.EchoReply {
			icmpPacket.Header.PrintICMPHeader()
		}
	}
	// decriment ttl. if ttl is zero, throw away
	packet.Header.DeclTTL()
	if packet.Header.TTL <= 0 {
		err := icmpTimeExceeded(r.Devices[devNo])
		if err != nil {
			return fmt.Errorf("ttl is zero and error while sending icmp message: %v", err)
		}
		return fmt.Errorf("ttl is zero")
	}
	// check destination segment
	var hwaddr *ethernet.HardwareAddress
	if r.Devices[devNo].NetInfo().IsInSegment(packet.Header.Dst) {
		// to my address
		if packet.Header.Dst == r.Devices[devNo].IPAddress() { // この比較大丈夫？
			fmt.Printf("%v: received to my address\n", r.Devices[devNo].Name())
			return nil
		}
		// check mac address talble
		entry := r.ArpTables[devNo].Search(packet.Header.Dst.Bytes())
		if entry == nil {
			// there is no entry
			fmt.Printf("[info]%v: there is no entry[%s]: send arp request\n", r.Devices[devNo].Name(), packet.Header.Dst.String())
			r.ArpTables[devNo].Show()
			// send arp request
			protocols := r.Devices[devNo].RegisteredProtocol()
			a := protocols[0].(*net.ARP)
			err = a.ARPRequest(packet.Header.Dst.Bytes(), arp.PROTOCOL_IPv4)
			if err != nil {
				return err
			}
			// add send buffer
			err = packet.ReCalculateChecksum()
			if err != nil {
				return err
			}
			buf, err := packet.Serialize()
			if err != nil {
				return err
			}
			frame.Data = buf
			frame.Header.Src = r.Devices[devNo].Address()
			frame.Header.Dst = ethernet.InvalidAddress
			r.ReadyQueue <- buffer{
				deviceNo: devNo,
				address:  packet.Header.Dst,
				frame:    frame,
			}
			return nil
		}
		// there is
		hwaddr, err = ethernet.Address(entry.HardwareAddress)
		if err != nil {
			return fmt.Errorf("failed to encode hwaddr:%v", err)
		}
	} else if r.Devices[!devNo].NetInfo().IsInSegment(packet.Header.Dst) {
		// ether segment of router interfaces
		fmt.Println(">>>>>>>>>>>>> other interface segment <<<<<<<<<<<<<")
		entry := r.ArpTables[!devNo].Search(packet.Header.Dst.Bytes())
		if entry == nil {
			fmt.Printf("[info]%v: there is no entry[%s]: send arp request\n", r.Devices[!devNo].Name(), packet.Header.Dst.String())
			r.ArpTables[!devNo].Show()
			// send arp request
			protocols := r.Devices[!devNo].RegisteredProtocol()
			a := protocols[0].(*net.ARP)
			err = a.ARPRequest(packet.Header.Dst.Bytes(), arp.PROTOCOL_IPv4)
			if err != nil {
				return err
			}
			// add send buffer
			err = packet.ReCalculateChecksum()
			if err != nil {
				return err
			}
			buf, err := packet.Serialize()
			if err != nil {
				return err
			}
			frame.Data = buf
			frame.Header.Src = r.Devices[!devNo].Address()
			frame.Header.Dst = ethernet.InvalidAddress
			r.ReadyQueue <- buffer{
				deviceNo: !devNo,
				address:  packet.Header.Dst,
				frame:    frame,
			}
			return nil
		}
		//if there is in arp table
		hwaddr, err = ethernet.Address(entry.HardwareAddress)
		if err != nil {
			return fmt.Errorf("failed to encode hwaddr: %v", err)
		}
	} else {
		// different segment
		// check mac address talble
		fmt.Printf("[info]%v: to next router(%s): %s\n", r.Devices[!devNo].Name(), r.Nexthops.String(), packet.Header.Dst.String())
		entry := r.ArpTables[!devNo].Search(r.Nexthops.Bytes())
		if entry == nil {
			// there is no entry
			r.ArpTables[!devNo].Show()
			fmt.Printf("[info]%v: there is no entry[%s]: send arp request\n", r.Devices[!devNo].Name(), r.Nexthops.String())
			// send arp request
			protocols := r.Devices[!devNo].RegisteredProtocol()
			a := protocols[0].(*net.ARP)
			err = a.ARPRequest(r.Nexthops.Bytes(), arp.PROTOCOL_IPv4)
			if err != nil {
				return err
			}
			// add send buffer
			err = packet.ReCalculateChecksum()
			if err != nil {
				return err
			}
			buf, err := packet.Serialize()
			if err != nil {
				return err
			}
			frame.Data = buf
			frame.Header.Src = r.Devices[!devNo].Address()
			frame.Header.Dst = ethernet.InvalidAddress
			r.ReadyQueue <- buffer{
				deviceNo: !devNo,
				address:  r.Nexthops,
				frame:    frame,
			}
			return nil
		}
		// there is
		hwaddr, err = ethernet.Address(entry.HardwareAddress)
		if err != nil {
			return fmt.Errorf("[error]failed to encode hwaddr: %v", err)
		}
	}
	// rewrite destination hardware address
	frame.Header.Dst = *hwaddr
	// rewrite source hardware address to router's hardware address
	frame.Header.Src = r.Devices[devNo].Address()
	// write to another(destination segment's) interface device
	// calculate checksum again because ttl is changed
	err = packet.ReCalculateChecksum()
	if err != nil {
		return fmt.Errorf("[error]failed to recalc checksum: %v", err)
	}
	buf, err := packet.Serialize()
	if err != nil {
		return fmt.Errorf("[error]failed to serialize ip packet: %v", err)
	}
	frame.Data = buf
	data, err := frame.Serialize()
	if err != nil {
		return fmt.Errorf("[error]failed to serialize ethernet frame: %v", err)
	}
	size, err := r.Devices[!devNo].Write(data)
	if err != nil {
		return fmt.Errorf("[error]%v: failed to write ethernet frame: %v", r.Devices[!devNo].Name(), err)
	}
	fmt.Printf("[info]%v: write to %s: %dbyte\n", r.Devices[!devNo].Name(), packet.Header.Dst.String(), size)
	return nil
}

func (r *Router) sendBuffer() {
	buf := make([]buffer, 128)
	for {
		b := <-r.ReadyQueue
		buf = append(buf, b)
		fmt.Printf("[info] got buffer from %v (to %s)\n", r.Devices[b.deviceNo].Name(), b.address.String())
		for i, b := range buf {
			entry := r.ArpTables[b.deviceNo].Search(b.address.Bytes())
			fmt.Println("send buffer routine------------------")
			r.ArpTables[b.deviceNo].Show()
			if entry != nil {
				// found

				hwaddr, err := ethernet.Address(entry.HardwareAddress)
				if err != nil {
					fmt.Printf("[error]%v: failed to encode hardware address: %v\n", r.Devices[b.deviceNo].Name(), err)
				}
				b.frame.Header.Dst = *hwaddr
				data, err := b.frame.Serialize()
				if err != nil {
					fmt.Printf("[error]%v: failed to serialize: %v\n", r.Devices[b.deviceNo].Name(), err)
				}
				size, err := r.Devices[b.deviceNo].Write(data)
				if err != nil {
					fmt.Printf("[error]%v: failed to write: %v\n", r.Devices[b.deviceNo].Name(), err)
				}
				fmt.Printf("[info](sendBuffer routine):%v:write to %s %dbytes\n", r.Devices[b.deviceNo].Name(), b.address.String(), size)
				buf = removeBuf(buf, i)
				continue
			}
		}
	}
}

func removeBuf(buf []buffer, num int) []buffer {
	if num >= len(buf) {
		return buf
	}
	return append(buf[:num], buf[num+1:]...)
}

func (r *Router) Run() {
	r.Devices[true].DeviceInfo()
	r.Devices[false].DeviceInfo()
	defer r.Devices[true].Close()
	defer r.Devices[false].Close()
	// handling arp packet
	go r.Devices[true].Handle()
	go r.Devices[false].Handle()
	go r.Devices[true].Next()
	go r.Devices[false].Next()
	// handling ip packet and so on
	go r.sendBuffer()
	fmt.Println("------------------router start--------------------")
	r.Handle()
}

func icmpTimeExceeded(dev net.Device) error {

	return nil
}
