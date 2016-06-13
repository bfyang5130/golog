package main

import (
	"fmt"
)

var (
	intcount int
	isOut    chan bool
)

func main() {
	//定义一个控制主线程退出的处理
	intcount = 1

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
	//退出堵塞
	<-isOut
}

func Product(queue chan<- int) {

	for i := 0; i < 1000; i++ {
		queue <- i
		fmt.Println("send:", i)
	}
}

func Consumer(queue <-chan int, j int) {
	//用来做一下记录，看看线程是不是一直在运行的
	//打开一个文件向里面写东西
	//wf, err := os.OpenFile(`./testchan.txt`, os.O_APPEND, 0775)
	//if err != nil {
	//	return
	//}

	for {
		v := <-queue
		fmt.Println(`receive:`, v)
		//_, err1 := io.WriteString(wf, fmt.Sprintf(`%d`, v)+"\n")
		//if err1 != nil {
		//	fmt.Println("can not write file")
		//}
	}
	addCount()
}

func addCount() {
	intcount++
	if intcount == 10 {
		isOut <- true
	}
}
