package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	w sync.WaitGroup
)

func main() {
	/**
	configName := `dfdcdc/dfdfd/cdcdsadf/dfdfdd/www.go.com.access.log`
	website := strings.Split(configName, `/`)
	logname := website[len(website)-1]
	splitlog := strings.Split(logname, `.`)
	visitwebsite := strings.Replace(logname, splitlog[len(splitlog)-1], ``, -1)
	visitwebsite = strings.Replace(visitwebsite, splitlog[len(splitlog)-2], ``, -1)
	visitwebsite = strings.Replace(visitwebsite, `..`, ``, -1)
	*/
	if t1, err := time.Parse("02/Jan/2006:15:04:05", strings.TrimSpace(`04/Jan/2006:15:04:05`)); err == nil {
		fmt.Println(t1.Format(`2006-01-02 15:04:05`))
	} else {
		fmt.Println(`nono`)
	}
	return

	//生成一个线程用来往通道里面放东西
	intChannel := make(chan int)
	w.Add(1)
	go Product(intChannel)
	//生成10个线程取东西
	for j := 0; j < 10; j++ {
		w.Add(1)
		go Consumer(intChannel)
	}
	w.Wait()
	/**
	for i := 0; i < 1000; i++ {
		intChannel <- i
		fmt.Println("send:", i)
		v := <-intChannel
		fmt.Println("receive:", v)
	}
	*/
	//退出堵塞
	//<-isOut
}

func Product(queue chan<- int) {

	defer func() {
		fmt.Println(`close channl000`)
		close(queue)
		w.Done()
	}()
	for i := 0; i < 1000000; i++ {
		queue <- i
		//fmt.Println("send:", i)
	}
}

func Consumer(queue <-chan int) {
	//用来做一下记录，看看线程是不是一直在运行的
	//打开一个文件向里面写东西
	//wf, err := os.OpenFile(`./testchan.txt`, os.O_APPEND, 0775)
	//if err != nil {
	//	return
	//}

	//var v int
	ok := true

	for ok {
		if _, ok = <-queue; ok {
			//fmt.Println(`receive:`, v)
		} else {
			fmt.Println(`close channl`)
			w.Done()
		}
	}
}
