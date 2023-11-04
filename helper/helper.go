package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// 报错并登出
func ErrorMsg(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 获得本地ip
func GetIp() string {
	type response struct {
		Ip string `json:"origin"`
	}
	url := "https://httpbin.org/get"
	data := Fetch(url)
	res := response{}
	json.Unmarshal([]byte(data), &res)
	return res.Ip
}

// 返回本机ipv6地址
func Ipv6() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To16() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("IPv6 address not found")
}

// 获取局域网ipv4地址
func Ipv4() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

// 获取数据
func Fetch(url string) string {
	// 获取网页数据
	resp, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}
