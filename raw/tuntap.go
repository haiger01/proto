package raw

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/spectrex02/proto/ioctl"
)

const device = "/dev/net/tun"

type TunDevice struct {
	// File *os.File
	File io.ReadWriteCloser
	Name string
}

func NewTunDevice(name string) (*TunDevice, error) {
	devName, file, err := openTun(name)
	if err != nil {
		return nil, err
	}
	return &TunDevice{
		File: file,
		Name: devName,
	}, nil
}

func openTun(name string) (string, *os.File, error) {
	if len(name) >= syscall.IFNAMSIZ {
		return "", nil, fmt.Errorf("name is too long")
	}
	file, err := os.OpenFile(device, os.O_RDWR, 0600)
	if err != nil {
		return "", nil, err
	}
	name, err = ioctl.TUNSETIFF(file.Fd(), name, syscall.IFF_TAP|syscall.IFF_NO_PI)
	if err != nil {
		return "", nil, err
	}
	flags, err := ioctl.SIOCGIFFLAGS(name)
	if err != nil {
		file.Close()
		return "", nil, err
	}
	flags |= (syscall.IFF_UP | syscall.IFF_RUNNING)
	if err := ioctl.SIOCSIFFLAGS(name, flags); err != nil {
		file.Close()
		return "", nil, err
	}
	return name, file, nil
}

func (tun *TunDevice) Address() []byte {
	addr, _ := getAddress(tun.Name)
	return addr[:6]
}
