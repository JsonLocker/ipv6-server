package dao

import (
	"aliyun/aliyun"
	"aliyun/helper"
	"fmt"
	"log"
	"time"
)

// 查看本地ip,并把ip更新在线数据库白名单组auto_flash
func frushDbIpOnce() {
	ip := helper.GetIp()
	err := aliyun.IpWhite(ip, "auto_flash_group_name", "pc-bp12n5p82q7z0of44")
	helper.ErrorMsg(err)
	log.Print("数据库白名单添加成功! ip: " + ip + "\n")
}

// 更新数据库白名单
func DbWhiteIP() {
	var dataDir string
	fmt.Print("需要循环更新ip请输入 yes, 否则回车跳过 ：")
	fmt.Scanln(&dataDir)

	frushDbIpOnce()
	if dataDir == "yes" {
		fmt.Println("请勿关闭本窗口, 每隔10s更新一次好慕课专属白名单, 防止ip变化造成爬虫数据入库中断.")
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			frushDbIpOnce()
		}
	}
}
