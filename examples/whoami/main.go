/*
whoami example of using the proxy manager. Please note that this program
connects to CloudFlare's server (https://www.cloudflare.com/cdn-cgi/trace).
This example cannot be ran without a text file "proxies.txt" with HTTP proxy.
*/

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
