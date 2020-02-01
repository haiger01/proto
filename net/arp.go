package net

import (
	"bytes"
	"fmt"

	"github.com/spectrex02/router-shakyo-go/arp"
	"github.com/spectrex02/router-shakyo-go/ethernet"
)

type ARP struct {
	HardwareType ethernet.EtherType
	Table        *arp.ARPTable
	Dev          Device
}

func newARP(dev Device) *ARP {
	return &ARP{
		HardwareType: ethernet.ETHER_TYPE_ARP,
		Table:        arp.NewARPTable(),
		Dev:          dev,
	}
}

func (a *ARP) Handle(data []byte) error {
	packet, err := arp.NewARPPacket(data)
	if err != nil {
		return fmt.Errorf("failed to create ARP packet")
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	packet.String()
	if packet.Header.HardwareType != arp.HARDWARE_ETHERNET {
		return fmt.Errorf("invalid hardware type")
	}
	if packet.Header.ProtocolType != arp.PROTOCOL_IPv4 && packet.Header.ProtocolType != arp.PROTOCOL_IPv6 {
		return fmt.Errorf("invalid protocol type")
	}
	mergeFlag, err := a.Table.Update(packet.SourceHardwareAddress, packet.SourceProtocolAddress)
	if err != nil {
		return err
	}
	if bytes.Equal(packet.TargetProtocolAddress, a.Dev.ProtocolAddressIP().Bytes()) {
		if !mergeFlag {
			err := a.Table.Insert(packet.SourceHardwareAddress, packet.SourceProtocolAddress, packet.Header.ProtocolType)
			if err != nil {
				return fmt.Errorf("Failed to insert: %v", err)
			}
		}
		if packet.Header.OpCode == arp.ARP_REQUEST {
			err := a.ARPReply(packet.SourceHardwareAddress, packet.SourceProtocolAddress, packet.Header.ProtocolType)
			if err != nil {
				return err
			}
		}
	}
	a.Table.Show()
	return nil
}

func (a *ARP) Type() ethernet.EtherType {
	return a.HardwareType
}

func (a *ARP) ARPRequest(targetProtocolAddress []byte, protocolType arp.ProtocolType) error {
	hwaddr := a.Dev.Address()
	protoaddr := a.Dev.ProtocolAddressIP()
	request, err := arp.Request(hwaddr[:], protoaddr[:], targetProtocolAddress, protocolType)
	if err != nil {
		return fmt.Errorf("failed to create ARP request")
	}
	buf, err := request.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize")
	}
	_, err = a.Dev.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to send ARP request")
	}
	return nil
}

func (a *ARP) ARPReply(targetHardwareAddress, targetProtocolAddress []byte, protocolType arp.ProtocolType) error {
	hwaddr := a.Dev.Address()
	protoaddr := a.Dev.ProtocolAddressIP()
	reply, err := arp.Reply(hwaddr[:], protoaddr[:], targetHardwareAddress, targetProtocolAddress, protocolType)
	if err != nil {
		return fmt.Errorf("failed to create arp reply")
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>")
	reply.String()
	buf, err := reply.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize")
	}
	_, err = a.Dev.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to send arp reply")
	}
	return nil
}

func (a *ARP) Write(dst []byte, protocol interface{}, data []byte) (int, error) {
	// d, err := ethernet.Address(dst)
	// if err != nil {
	// return 0, err
	// }
	// return a.Dev.Write(data)
	return 0, fmt.Errorf("this function is dummy function")
}
