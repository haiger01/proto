package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func DisableIPForward() error {
	// file, err := os.Open("/proc/sys/net/ipv4/ip_forward")
	file, err := os.OpenFile("/proc/sys/net/ipv4/ip_forward", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open ip_forward file")
	}
	defer file.Close()
	_, err = file.WriteString("0")
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}
	return nil
}

func ParseAddressAndPort(str string) ([]byte, uint16, error) {
	// required -> address:port
	parsed := strings.Split(str, ":")
	port, err := strconv.Atoi(parsed[1])
	if err != nil {
		return nil, 0, err
	}
	var addr []byte
	addrStr := strings.Split(parsed[0], ".")
	for _, v := range addrStr {
		a, err := strconv.Atoi(v)
		if err != nil {
			return nil, 0, err
		}
		addr = append(addr, byte(a))
	}
	return addr, uint16(port), nil
}
