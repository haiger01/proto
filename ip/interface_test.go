package ip

import (
	"fmt"
	"testing"
)

func TestNewIPSubnetMask(t *testing.T) {
	a, err := NewIPSubnetMask("192.168.2.251/24")
	if err != nil {
		t.Fatal(err)
	}
	a.Show()
}

func TestIsInSegment(t *testing.T) {
	a, err := NewIPSubnetMask("192.168.2.251/24")
	if err != nil {
		t.Fatal(err)
	}
	a.Show()
	fmt.Println(a.IsInSegment(IPAddress{192, 168, 3, 4}))
}
