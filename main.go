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
	var proxy51 bool
	var proxyfofa bool
	var fofaupdate bool
	var useTime int
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
			cli.BoolFlag{
				Name:        "51",
				Usage:       "使用51代理进行代理地址获取",
				Destination: &proxy51,
			},
			cli.BoolFlag{
				Name:        "fofa",
				Usage:       "使用fofa代理池",
				Destination: &proxyfofa,
			},
			cli.BoolFlag{
				Name:        "fofaup",
				Usage:       "更新fofa代理池",
				Destination: &fofaupdate,
			},
			cli.IntFlag{
				Name:        "time,t",
				Usage:       "一个ip默认的使用时间",
				Destination: &useTime,
				Value:       20,
			},
		},
		Action: func(c *cli.Context) error {
			// 创建管道
			Addrch := make(chan string)
			MyConfig := config.ReadConfig()
			proxyAddrs, err := proxy.ReadFofaProxyFile()
			if err != nil {
				log.Fatal(err)
			}
			go func(ch chan<- string) {
				var proxyAddr string
				for {
					// 获取代理地址
					if proxy51 == true {
						proxyAddr = proxy.Get51ProxyAddr(MyConfig.Url51)
					} else if proxyfofa == true {
						var targetBool bool
						proxyAddr, targetBool = proxy.GetFofaProxyAddr(&proxyAddrs)
						if !targetBool {
							log.Println("[-] 代理池用完了！重头再来一遍！")
							proxyAddrs, err = proxy.ReadFofaProxyFile()
							if err != nil {
								log.Fatal(err)
							}
						}
					} else if fofaupdate == true {
						proxy.UpdateFofaProxyAddr(MyConfig.FofaEmail, MyConfig.FofaApiKey)
						log.Println("[+] socks5服务器更新成功！")
						os.Exit(1)
					} else {
						log.Fatal("请选择使用的代理池！")
					}
					// 将代理地址发送到管道getProxyAddr
					ch <- proxyAddr
					// 打印提示信息
					log.Printf("[*] 本地代理监听在:%s，将流量转发到 SOCKS5 代理服务器 %s\n", localPort, proxyAddr)
					time.Sleep(time.Duration(useTime) * time.Second)
				}
			}(Addrch)
			server.ProxySocks(localPort, Addrch)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
