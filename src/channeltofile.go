package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"

	"myfunc"

	"log"

	"github.com/widuu/goini"
)

var (
	w           sync.WaitGroup
	logFileName = flag.String("log", "cServer.log", "Log file name")
	flag_status bool
)

func main() {
	//定义一个多线程

	//设置文件目录
	var filepath string
	var configfile string
	var fitfilelist string
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//write log
	log.Printf("Server abort! Cause:%v \n", "test log file")

	fitfilelist = "F:/nglog/config/fitfile.ini"
	fixconf := goini.SetConfig(fitfilelist)
	fmt.Println(fixconf)
	//设置读取文件的目录
	filepath = "F:/nglog/"
	//设置配置文件
	configfile = "F:/nglog/config/config.ini"
	//如果不存在配置文件就生成配置文件
	fileSize := myfunc.CheckFileSize(configfile)
	fmt.Println(fileSize)
	if fileSize != 1 {
		return
	}
	//read configfile
	conf := goini.SetConfig(configfile)
	//read file
	list, err := ioutil.ReadDir(filepath)
	if err != nil {
		fmt.Println("Wrong Dir")
		return
	}
	//循环目录历表获得文件信息
	for _, info := range list {
		//得到文件名称
		configName := info.Name()
		if info.IsDir() == true {
			fmt.Println(configName + " is Dir")
			continue
		}
		//是否存在配置文件中要处理的文件
		isfitfile := fixconf.GetValue("fitfile", configName)
		if len(isfitfile) == 0 || isfitfile == "no value" {
			fmt.Println(configName + "不在配置文件中,不需要处理")
			continue
		}
		//get lasttime filesize
		rfileSize := conf.GetValue("filesize", configName)

		var fileSize int64
		if len(rfileSize) != 0 {
			fileSize, err = strconv.ParseInt(rfileSize, 10, 64)
		} else {
			fileSize = 0
		}
		fmt.Println(info.Name())
		fmt.Println(fileSize)
		fmt.Println(info.Size())
		if fileSize == info.Size() {
			fmt.Println("the same file size,not need to fix")
			return
		}
		//创建一个有100个缓冲的通道
		ch := make(chan string, 100)
		//生产者读取数据并放入通道
		w.Add(1)
		go produce(ch, filepath+info.Name())
		//循环读取文件,然后分成100个线程去处理数据
		for chi := 1; chi <= 100; chi++ {
			w.Add(1)
			go consumer(ch)
		}
		w.Wait()

	}
}

//生产者
func produce(p chan<- string, file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("A Wrong File")
		defer w.Done()
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		line, isPrefix, err := br.ReadLine()
		if err == io.EOF {
			fmt.Println(err)
			w.Done()
			break
		} else {
			if !isPrefix {
				//将一个元素推入通道
				p <- string(line)
				fmt.Println(`send:`, string(line))
			}
		}

	}
}

//消费者
func consumer(c <-chan string) {
	for i := 100; i < 100; i++ {
		v := <-c
		fmt.Println(`receive:`, v)
	}
	w.Done()
}
