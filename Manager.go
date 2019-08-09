package proxy

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Manager is struct type containing proxies
type Manager struct {
	Proxies []*Proxy
	Status  int
}

// OK and LOCKED are possible status for the proxy manager. Locked will prevent
// clients from getting any proxies.
const (
	OK = iota
	LOCKED
)

// NewManager returns a new proxy manager. Seeds the PRNG with UnixNano time.
func NewManager() Manager {
	rand.Seed(time.Now().UnixNano()) //Seed PRNG
	return Manager{}
}

// Get returns a new Proxy. Can be used with ClientFromProxy to get HTTP client.
func (m *Manager) Get() (*Proxy, error) {
	// If the manager is locked, pause goroutine for 200ms and try again.
	if m.Status == LOCKED {
		time.Sleep(200 * time.Millisecond)
		return m.Get()
	}

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

// RemoveBadProxies removes the proxies in the Manager with status "BAD".
func (m *Manager) RemoveBadProxies() {
	// Make a new slice with not bad proxies. Using this method instead of index
	// swapping, so we don't end up changing the index' while for range.
	var ps []*Proxy
	for _, p := range m.Proxies {
		if p.Status != BAD {
			ps = append(ps, p)
		}
	}
	m.Proxies = ps
}

// AppendProxiesFromReader appends proxies from an io.Reader using a scanner.
// First argument specifies the protocol of the proxies.
// These should be proxy.HTTP, proxy.SOCKS4 or proxy.SOCKS5.
// HTTPS should be specified as HTTP. SOCKS4a should be specified as SOCKS4.
// Each token of the scanner should be a new proxy in the form: host:port
func (m *Manager) AppendProxiesFromReader(protocol int, r io.Reader) error {
	scanner := bufio.NewScanner(r)

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
	return scanner.Err()
}

// AppendProxiesFromFile appends proxies from a file to the proxy manager.
// First argument specifies the protocol of the proxies.
// The content of the file is processed by AppendProxiesFromReader.
// Each line of the file should be a new proxy in the form: host:port
func (m *Manager) AppendProxiesFromFile(protocol int, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	err = m.AppendProxiesFromReader(protocol, file)
	return err
}

// AppendProxiesFromURL appends proxies from an URL. The ressource should be
// available through a blank HTTP GET request.
// First argument specifies the protocol of the proxies.
// The body of the response is processed by AppendProxiesFromReader.
func (m *Manager) AppendProxiesFromURL(protocol int, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = m.AppendProxiesFromReader(protocol, resp.Body)
	return err
}
