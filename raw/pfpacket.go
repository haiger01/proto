package raw

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/spectrex02/router-shakyo-go/ioctl"
)

type PFPacket struct {
	Fd   int
	Name string
}

func NewPFPacket(name string) (*PFPacket, error) {
	fd, err := openPFPacket(name)
	if err != nil {
		return nil, err
	}
	return &PFPacket{
		Fd:   fd,
		Name: name,
	}, nil
}

func (p PFPacket) String() string {
	return p.Name
}

func (p PFPacket) Address() []byte {
	addr, _ := getAddress(p.Name)
	return addr[:6]
}

func openPFPacket(name string) (int, error) {
	if name == "" {
		return -1, fmt.Errorf("name is empty")
	}
	if len(name) >= syscall.IFNAMSIZ {
		return -1, fmt.Errorf("name is too long")
	}
	protocol := hton16(syscall.ETH_P_ALL)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(protocol))
	if err != nil {
		return -1, err
	}
	index, err := ioctl.SIOCGIFINDEX(name)
	if err != nil {
		syscall.Close(fd)
		return -1, err
	}
	addr := &syscall.SockaddrLinklayer{
		Protocol: protocol,
		Ifindex:  int(index),
	}
	if err = syscall.Bind(fd, addr); err != nil {
		syscall.Close(fd)
		return -1, err
	}
	flags, err := ioctl.SIOCGIFFLAGS(name)
	if err != nil {
		syscall.Close(fd)
		return -1, nil
	}
	flags |= syscall.IFF_PROMISC
	if err := ioctl.SIOCSIFFLAGS(name, flags); err != nil {
		syscall.Close(fd)
		return -1, nil
	}
	return fd, nil
}

func getAddress(name string) ([]byte, error) {
	addr, err := ioctl.SIOCGIFHWADDR(name)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func hton16(i uint16) uint16 {
	var ret uint16
	binary.BigEndian.PutUint16((*[2]byte)(unsafe.Pointer(&ret))[:], i)
	return ret
}
