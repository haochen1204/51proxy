package proxy

import (
	"bufio"
	"fmt"
	"github.com/haochen1204/HCGOTools/fofa"
	"log"
	"net"
	"os"
	"time"
)

func GetFofaProxyAddr(targetAddrs *[]string) (string, bool) {
	var targetAddr string
	for index, value := range *targetAddrs {
		targetAddr = value
		*targetAddrs = append((*targetAddrs)[:index], (*targetAddrs)[index+1:]...)
		if len(*targetAddrs) == 0 {
			return targetAddr, false
		}
		if testTargetServer(targetAddr) {
			break
		}
	}
	return targetAddr, true
}

func UpdateFofaProxyAddr(email, apiKey string) {
	var targetAddrs []string
	fofaclient := fofa.New_FoFa_Client(email, apiKey)
	fofaQ := "after=\"2023-10-01\" && protocol=\"socks5\" && country=\"CN\" && \"Method:No Authentication(0x00)\""
	fofainfo := fofa.New_FoFa_InfoSearch(fofaQ)
	fofainfo.Size = 10000
	fofaresult, err := fofaclient.HostSearch(fofainfo)
	if err != nil {
		log.Fatal(err)
	}
	for _, value := range fofaresult.Results {
		targetAddrs = append(targetAddrs, fmt.Sprintf("%s:%s", value[1], value[2]))
	}
	writeFile(&targetAddrs)
}

func writeFile(targetAddrs *[]string) {
	file, err := os.Create("fofaproxy.txt")
	if err != nil {
		log.Println("[-] 创建文件失败:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range *targetAddrs {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Println("[-] 写入文件错误:", err)
			return
		}
	}
	// 刷新缓冲区，确保所有数据都写入文件
	err = writer.Flush()
	if err != nil {
		log.Println("[-] 缓冲区刷新错误:", err)
		return
	}
}

func ReadFofaProxyFile() ([]string, error) {
	file, err := os.Open("fofaproxy.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func testTargetServer(targetProxyAddr string) bool {
	targetServer := "www.baidu.com:80"
	timeout := 5 * time.Second // 设置超时时间为5秒

	// 连接到 SOCKS5 服务器
	conn, err := net.DialTimeout("tcp", targetProxyAddr, timeout)
	if err != nil {
		log.Println("[-] 连接到 SOCKS5 服务器时发生错误:", err)
		return false
	}
	defer conn.Close()

	// 发送 SOCKS5 握手请求
	_, err = conn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		log.Println("[-] 发送 SOCKS5 握手请求时发生错误:", err)
		return false
	}

	// 读取 SOCKS5 握手响应
	response := make([]byte, 2)
	_, err = conn.Read(response)
	if err != nil {
		log.Println("[-] 读取 SOCKS5 握手响应时发生错误:", err)
		return false
	}

	// 检查握手响应是否为 0x05 0x00，表示成功
	if response[0] != 0x05 || response[1] != 0x00 {
		log.Printf("[-] %s SOCKS5 握手失败.", targetProxyAddr)
		return false
	}

	// 解析目标服务器地址
	targetAddr, err := net.ResolveTCPAddr("tcp", targetServer)
	if err != nil {
		log.Println("[-] 解析目标服务器地址时发生错误:", err)
		return false
	}

	targetIP := targetAddr.IP
	targetPort := uint16(targetAddr.Port)

	// 发送 SOCKS5 连接请求
	request := []byte{0x05, 0x01, 0x00, 0x01}
	request = append(request, targetIP.To4()...)
	request = append(request, byte(targetPort>>8), byte(targetPort))

	_, err = conn.Write(request)
	if err != nil {
		log.Println("[-] 发送 SOCKS5 连接请求时发生错误:", err)
		return false
	}

	// 读取 SOCKS5 连接响应
	response = make([]byte, 10)
	_, err = conn.Read(response)
	if err != nil {
		log.Println("[-] 读取 SOCKS5 连接响应时发生错误:", err)
		return false
	}

	// 检查连接响应是否为 0x05 0x00，表示成功
	if response[0] != 0x05 || response[1] != 0x00 {
		log.Printf("[-] %s SOCKS5 连接请求失败.", targetProxyAddr)
		return false
	}

	log.Printf("[+] %s SOCKS5 服务器成功运行.", targetProxyAddr)
	return true
}
