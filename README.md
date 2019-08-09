# proxy
[![GoDoc](https://godoc.org/github.com/multimikael/proxy?status.png)](http://godoc.org/github.com/multimikael/proxy)

proxy is a simple proxy manager library. It supports HTTP(s)/SOCK4(a)/SOCKS5.


# Installation
proxy can be installed using `go get`:
```sh
go get github.com/multimikael/proxy
```

# Example
This is a simple example of using the proxy manager and ClientFromProxy. This example reads HTTP proxies from a text file "proxies.txt". It gets a proxy from the manager and passes it to ClientFromProxy to get an HTTP client. Finally it makes a GET request through the proxy. This a cut from the whoami example in the [examples folder](https://github.com/multimikael/proxy/tree/master/examples/).
```go
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/multimikael/proxy"
)

func main() {
	// Create proxy manager and read from proxies from file "proxies.txt"
	pm := proxy.NewManager()
	pm.AppendProxiesFromFile(proxy.HTTP, "proxies.txt")

	// Get a random alive proxy from the proxy manager
	p, err := pm.Get()
	if err != nil {
		panic(err.Error())
	}

	// Get an HTTP client with the proxy, and timeout after 10 seconds
	client, err := proxy.ClientFromProxy(p, 10)
	if err != nil {
		panic(err.Error())
	}

	// Make a GET request through the proxy
	resp, err := client.Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```
