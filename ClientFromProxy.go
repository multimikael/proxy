package proxy

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"h12.io/socks"
)

// ClientFromProxy returns a HTTP client of type http.Client.
// The client will timeout after "timeout" seconds of dialing the proxy.
// Please note that the HTTP client has no transport or overall connection
// timeout.
//
// HTTP proxy dialers are created using net/http. SOCKS4 and SOCKS5 proxy
// dialers use h12.io/socks.
func ClientFromProxy(p *Proxy, timeout time.Duration) (*http.Client, error) {
	switch p.Protocol {
	case HTTP:
		uri := url.URL{
			Scheme: "http",
			Host:   p.Host + ":" + p.Port,
		}
		tr := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: timeout * time.Second,
			}).Dial,
			Proxy: http.ProxyURL(&uri),
		}
		return &http.Client{Transport: tr}, nil
	case SOCKS4:
		uri := url.URL{
			Scheme:   "socks4",
			Host:     p.Host + ":" + p.Port,
			RawQuery: "timeout=" + string(timeout) + "s",
		}
		dialSocksProxy := socks.Dial(uri.String())
		tr := &http.Transport{Dial: dialSocksProxy}
		return &http.Client{Transport: tr}, nil
	case SOCKS5:
		uri := url.URL{
			Scheme:   "socks5",
			Host:     p.Host + ":" + p.Port,
			RawQuery: "timeout=" + string(timeout) + "s",
		}
		dialSocksProxy := socks.Dial(uri.String())
		tr := &http.Transport{Dial: dialSocksProxy}
		return &http.Client{Transport: tr}, nil
	default:
		return nil, errors.New("unknown proxy protocol")
	}
}
