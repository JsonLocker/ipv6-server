package dao

import (
	"aliyun/aliyun"
	"aliyun/helper"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 根据输入的域名解析到当前
func ipv6Serv(domain string) {
	// 提取一级域名
	hostParts := strings.Split(domain, ".")
	var RR string
	switch len(hostParts) {
	case 2:
		RR = "@"
	case 3:
		RR = hostParts[0]
		domain = hostParts[1] + "." + hostParts[2]
	default:
		fmt.Println("域名格式错误!")
		os.Exit(0)
	}

	// 得到域名列表
	// list, err := aliyun.RecordList("xcmooc.com")
	list, err := aliyun.RecordList(domain)
	helper.ErrorMsg(err)

	// 获取dns id
	var recordId int = 0
	for _, v := range list {
		if v.RR == RR {
			recordId, err = strconv.Atoi(v.RecordId)
			helper.ErrorMsg(err)
		}
	}

	// 获取ipv6地址
	ipv6, err := helper.Ipv6()
	helper.ErrorMsg(err)

	// 解析域名到本地ipv6地址, 没有创建 有更新
	if recordId != 0 {
		aliyun.UpdateDns(RR, strconv.Itoa(recordId), ipv6, "AAAA")
		fmt.Println(recordId)
	} else {
		aliyun.AddDns(RR, domain, ipv6, "AAAA")
		log.Print("create DNS")
	}

	fmt.Printf("外网访问地址：http://%s.%s\n", RR, domain)
}

// 创建ipv6服务器
func BuildServe() {
	// 输入需要共享的路径
	var dataDir, domain string
	fmt.Print("请输入需要分享的磁盘路径：")
	fmt.Scanln(&dataDir)
	if dataDir == "" {
		fmt.Println("路径为空，无法启动")
		os.Exit(0)
	}
	fmt.Print("请输入要解析的域名如 a.abc.com 暂时只支持1级和2级域名, 没有请回车跳过 ：")
	fmt.Scanln(&domain)
	if domain != "" {
		ipv6Serv(domain)
	}

	// 提示本地局域网地址
	ipv4 := helper.Ipv4()
	fmt.Printf("请访问本机地址：http://%s\n", ipv4)
	// 搭建文件服务器
	http.ListenAndServe(":80", http.FileServer(http.Dir(dataDir)))
}
