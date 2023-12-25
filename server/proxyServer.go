package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"socks5proxy/proxy"
	"time"
)

func ProxySocks(localPort string, Addrch chan proxy.Data) {
	// 在本地监听指定端口
	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		log.Fatalf("启动本地监听器时发生错误: %v", err)
	}
	defer listener.Close()

	var proxyAddr string = ""
	loc, err := time.LoadLocation("Local")
	mm, _ := time.ParseDuration("5s")
	for {
		// 从管道中接收值
		select {
		case proxyAddrFromChan := <-Addrch:
			now := time.Now()
			now = now.Add(mm)
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", proxyAddrFromChan.ExpireTime, loc)
			//log.Println(now, t)
			//log.Println(now.Before(t))
			if now.Before(t) {
				proxyAddr = proxyAddrFromChan.Ip
			} else {
				//log.Println(proxyAddrFromChan.Ip + " 已过期!")
				continue
			}
		}

		if proxyAddr == "" {
			log.Println("等待代理服务器中...")
			time.Sleep(2 * time.Second)
			continue
		}

		// 接受本地客户端连接
		clientConn, err := listener.Accept()
		err = clientConn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			log.Printf("接受客户端连接时发生错误: %v", err)
			continue
		}

		// 在新的 goroutine 中处理连接
		go handleConnection(clientConn, proxyAddr)
	}
}

// 处理连接请求
func handleConnection(clientConn net.Conn, proxyAddr string) {
	defer clientConn.Close()

	// 连接到 SOCKS5 代理服务器
	proxyConn, err := net.Dial("tcp", proxyAddr)
	err = clientConn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		log.Printf("连接到 SOCKS5 代理服务器时发生错误: %v", err)
		return
	}
	defer proxyConn.Close()

	// 用于传递错误信息的通道
	var errCh = make(chan error, 1)

	// 从本地客户端到 SOCKS5 代理服务器的复制
	go func() {
		_, err := io.Copy(proxyConn, clientConn)
		if err != nil {
			errCh <- fmt.Errorf("复制数据到 SOCKS5 代理服务器时发生错误: %v", err)
		}
		errCh <- nil
	}()

	// 从 SOCKS5 代理服务器到本地客户端的复制
	go func() {
		_, err := io.Copy(clientConn, proxyConn)
		if err != nil {
			errCh <- fmt.Errorf("从 SOCKS5 代理服务器复制数据到本地客户端时发生错误: %v", err)
		}
		errCh <- nil
	}()

	// 等待复制完成或发生错误
	if err := <-errCh; err != nil {
		log.Println(err)
	}
}
