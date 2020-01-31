package util

import (
	"fmt"
	"os"
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
