package server

import (
	"golang.org/x/net/proxy"
	"net/http"
	myproxy "socks5proxy/proxy"
)

func CheckProxy(Addrch chan myproxy.Data, CheckAddrch chan myproxy.Data) {
	for {
		ipData := <-CheckAddrch
		go Check(ipData, Addrch)
	}
}

func Check(proxy_Data myproxy.Data, Addrch chan myproxy.Data) {
	proxy_addr := proxy_Data.Ip
	dialer, err := proxy.SOCKS5("tcp", proxy_addr, nil, proxy.Direct)
	if err != nil {
		return
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	if resp, err := httpClient.Get("https://www.baidu.com"); err == nil {
		if resp.StatusCode == 200 {
			Addrch <- proxy_Data
		}
		defer resp.Body.Close()
	}
}
