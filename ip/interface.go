package ip

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type IPSubnetMask struct {
	Address IPAddress
	Subnet  IPAddress
	Netmask IPAddress
}

func NewIPSubnetMask(info string) (*IPSubnetMask, error) {
	str := strings.Split(info, "/")
	s := strings.Split(str[0], ".")
	var addr []byte
	for _, v := range s {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		addr = append(addr, byte(n))
	}
	addrSubnetMask := &IPSubnetMask{
		Address: NewIPAddress(addr),
	}
	mask, err := strconv.Atoi(str[1])
	if err != nil {
		return nil, err
	}
	netmask := []byte{255, 255, 255, 255}
	a := mask / 8
	b := 8 - (mask % 8)
	masking := uint8(0)
	for i := 7; i >= b; i-- {
		fmt.Println("shift")
		masking += 1 << i
	}
	if a != 4 {
		netmask[a] = masking & netmask[a]
	}
	for i := a + 1; i < 4; i++ {
		netmask[i] = 0
	}
	addrSubnetMask.Netmask = NewIPAddress(netmask)
	subnet := []byte{byte(addr[0] & netmask[0]), byte(addr[1] & netmask[1]), byte(addr[2] & netmask[2]), byte(addr[3] & netmask[3])}
	addrSubnetMask.Subnet = NewIPAddress(subnet)
	return addrSubnetMask, nil
}

func (sm IPSubnetMask) Show() {
	fmt.Println("----------subnetmask----------")
	fmt.Printf("address = %s\n", sm.Address.String())
	fmt.Printf("subnet = %s\n", sm.Subnet.String())
	fmt.Printf("netmask = %s\n", sm.Netmask.String())
}

func (sm IPSubnetMask) IsInSegment(addr IPAddress) bool {
	// addr & netmask == subnet -> true
	b := addr.Bytes()
	maskedAddr := make([]byte, 4)
	for i, v := range sm.Netmask[:] {
		maskedAddr[i] = v & b[i]
	}
	// fmt.Printf("interface address-> %s ::: searching address-> %s ::: netmask-> %v\n", sm.Address.String(), addr.String(), maskedAddr)
	return bytes.Equal(sm.Subnet.Bytes(), maskedAddr)
}
