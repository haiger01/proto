package net

import (
	"context"
	"fmt"

	"github.com/spectrex02/proto/ip"
)

type Dialer struct {
	LocalAddr ip.IPAddress
}

func (d *Dialer) Dial(ctx context.Context, network, address string) (Conn, error) {
	switch network {
	case "udp":
		return DialUDP(ctx, address)
	// case "tcp":
	default:
		return nil, fmt.Errorf("unknown network protocol")
	}
}
