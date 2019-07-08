package proxy

import "fmt"

// Proxy is a struct type containing the protocol, host, port and status of the
// proxy. A Proxy should be managed by a Manager.
//
// Protocol should be of the constants HTTP, SOCKS4 or SOCKS5.
// Status should be of the constants ALIVE, BUSY or BAD.
type Proxy struct {
	Protocol int
	Host     string
	Port     string
	Status   int
}

// HTTP, SOCKS4 and SOCKS5 are possible proxy protocols as constants iota.
const (
	HTTP = iota
	SOCKS4
	SOCKS5
)

// ALIVE, BUSY and BAD are possible proxy status' as constants iota.
const (
	ALIVE = iota
	BUSY
	BAD
)

func (p *Proxy) String() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}
