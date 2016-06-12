package main

import (
	"fmt"
	"myfunc"

	//"github.com/yinheli/qqwry"
)

func main() {
	ipCountry, ipProvince, ipCity := myfunc.Convertip_tiny(`192.168.1.142`)
	fmt.Println(fmt.Sprintf(`国家：%s 省份：%s 城市：%s`, ipCountry, ipProvince, ipCity))
}
