package proxy

import (
	"bufio"
	"errors"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// Manager is struct type containing proxies
type Manager struct {
	Proxies []*Proxy
}

// NewManager returns a new proxy manager. Seeds the PRNG with UnixNano time.
func NewManager() Manager {
	rand.Seed(time.Now().UnixNano()) //Seed PRNG
	return Manager{}
}

// Get returns a new Proxy. Can be used with ClientFromProxy to get HTTP client.
func (m *Manager) Get() (*Proxy, error) {
	proxies, err := m.Alive()
	if err != nil {
		return &Proxy{}, err
	}
	return proxies[rand.Intn(len(proxies)-1)], nil
}

// Alive returns a list of alive proxies. These proxies are neither bad or busy.
func (m *Manager) Alive() ([]*Proxy, error) {
	if len(m.Proxies) == 0 {
		return nil, errors.New("Proxy pool empty")
	}
	var alive []*Proxy
	for _, p := range m.Proxies {
		if p.Status == ALIVE { // Alive is neither busy or bad
			alive = append(alive, p)
		}
	}
	return alive, nil
}

// AppendProxiesFromFile appends proxies from a file to the proxy manager.
// First argument specifies the protocol of the proxies.
// These should be proxy.HTTP, proxy.SOCKS4 or proxy.SOCKS5.
// HTTPS should be specified as HTTP. SOCKS4a should be specified as SOCKS4.
// Each line of the file should be a new proxy in the form: host:port
func (m *Manager) AppendProxiesFromFile(protocol int, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			host, port, _ := net.SplitHostPort(line)
			p := &Proxy{
				Protocol: protocol,
				Host:     host,
				Port:     port,
				Status:   ALIVE,
			}
			m.Proxies = append(m.Proxies, p)
		}
	}
	return nil
}
