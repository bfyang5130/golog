package main

import (
	"fmt"
	"math/rand"
)

// 函数 rand_generator_2，返回 通道(Channel)
func rand_generator_2() chan int {
	// 创建通道
	out := make(chan int)
	// 创建协程
	go func() {
		for {
			//向通道内写入数据，如果无人读取会等待
			out <- rand.Int()
		}
	}()
	return out
}

func main() {
	// 生成随机数作为一个服务
	rand_service_handler := rand_generator_2()
	// 从服务中读取随机数并打印
	fmt.Println(fmt.Sprintf("%d", <-rand_service_handler))
}
