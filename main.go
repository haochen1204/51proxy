package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"socks5proxy/config"
	"socks5proxy/proxy"
	"socks5proxy/server"
	"time"
)

func main() {
	var localPort string
	app := cli.App{
		Name:  "代理池",
		Usage: "socks5代理池",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "port,p",
				Usage:       "设置本地代理端口，默认为1080",
				Destination: &localPort,
				Value:       "1080",
			},
		},
		Action: func(c *cli.Context) error {
			// 创建管道 设置最大容量为100w
			Addrch := make(chan proxy.Data, 10000000)
			CheckAddrch := make(chan proxy.Data, 10000000)
			MyConfig := config.ReadConfig()
			go func() {
				for {
					// 获取代理地址
					proxy.Get51ProxyAddr(MyConfig.Url51, CheckAddrch)
					time.Sleep(1 * time.Second)
				}
			}()
			go server.CheckProxy(Addrch, CheckAddrch)
			log.Printf("[*] 本地代理监听在:%s", localPort)
			server.ProxySocks(localPort, Addrch)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
