package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data []struct {
		Ip string `json:"ip"`
	} `json:"data"`
	Success string `json:"success"`
}

// 从url获取代理服务器
func Get51ProxyAddr(url string) string {
	var targetAddr string

	response, err := http.Get(url)
	if err != nil {
		log.Printf("[-] 发送 HTTP GET 请求时发生错误：%v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[-] 读取响应体时发生错误：%v\n", err)
		os.Exit(1)
	}
	var jsonresponse Response
	err = json.Unmarshal(body, &jsonresponse)
	if err != nil {
		log.Printf("[-] 解析 JSON 时发生错误：%v\n", err)
	}
	targetAddr = jsonresponse.Data[0].Ip
	return targetAddr
}
