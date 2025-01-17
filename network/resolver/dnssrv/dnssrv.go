// Package dns srv resolves names to dns srv records
package dnssrv

import (
	"fmt"
	"net"

	"github.com/arun-spire/go-micro/network/resolver"
)

// Resolver is a DNS network resolve
type Resolver struct{}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	_, addrs, err := net.LookupSRV("network", "udp", name)
	if err != nil {
		return nil, err
	}
	var records []*resolver.Record
	for _, addr := range addrs {
		address := addr.Target
		if addr.Port > 0 {
			address = fmt.Sprintf("%s:%d", addr.Target, addr.Port)
		}
		records = append(records, &resolver.Record{
			Address: address,
		})
	}
	return records, nil
}
