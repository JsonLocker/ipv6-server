package main

import "aliyun/dao"

func main() {
	dao.BuildServe() // 搭建ipv6穿透网
	// dao.DbWhiteIP()  // 循环更新数据库ip白名单
}
