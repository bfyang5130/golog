package main

import (
	"fmt"
	"time"
)

func main() {

	//定义一个信道
	intChannel := make(chan int, 10)
	//生成一个线程用来往通道里面放东西
	go Product(intChannel)
	//生成10个线程取东西
	for j := 0; j < 10; j++ {
		go Consumer(intChannel, j)
	}
	/**
	for i := 0; i < 1000; i++ {
		intChannel <- i
		fmt.Println("send:", i)
		v := <-intChannel
		fmt.Println("receive:", v)
	}
	*/
	time.Sleep(30 * time.Second)
}

func Product(queue chan<- int) {

	for i := 0; i < 10000000; i++ {
		queue <- i
		fmt.Println("send:", i)
	}
}

func Consumer(queue <-chan int, j int) {
	for {
		v := <-queue
		fmt.Println(fmt.Sprintf("receive%d:%d", j, v))
	}
}
