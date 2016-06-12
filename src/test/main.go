package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"myfunc"

	"log"

	"github.com/widuu/goini"
	"gopkg.in/mgo.v2"
)

const URL = "192.168.8.188:27017" //mongodb连接字符串

var (
	mgoSession  *mgo.Session
	dataBase    = "nginx"
	w           sync.WaitGroup
	logFileName = flag.String("log", "cServer.log", "Log file name")
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
		fitstatus := toFitFile(fileSize, filepath+info.Name())
		fmt.Println(fitstatus)
		w.Wait()
	}

}

//判断文件是否存在
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func toFitFile(alreadyfitSize int64, file string) (fitstatus int) {
	//已经处理的大小
	fmt.Println(alreadyfitSize)
	//定义每个文件分割的大小
	var splitFileSize = int64(13107200)
	//分割成每份splitFileSize大小的文件
	splitStatus, splitNums := myfunc.SplitFile(file, splitFileSize, alreadyfitSize)
	if !splitStatus {
		fmt.Println("分割文件失败")
		return 0
	}
	//多通道处理文件
	var i int64
	fmt.Println("splintNums=" + strconv.FormatInt(splitNums, 10))

	for i = 1; i <= splitNums; i++ {
		w.Add(1)
		fmt.Println("iamhere")
		go sptoF(file, i)
	}
	fitstatus = 1
	return fitstatus
}

func sptoF(infile string, i int64) {
	fmt.Println("what is going on")
	//拼接成一个正常的文件
	newFileString := infile + strconv.FormatInt(i, 10)
	//打开文件
	fmt.Println("open:", newFileString)
	f, err := os.Open(newFileString)
	if err != nil {
		fmt.Println("A Wrong File")
		fmt.Println("finish", i)
		defer w.Done()
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	count := 0
	//用list来装这500个数据，然后一次性插入，因为一条条地插入会很耗又耗IO
	lineArr := list.New()
	for {
		fmt.Println(fmt.Printf(`第%d行`, count))
		line, isPrefix, err := br.ReadLine()
		if err == io.EOF {
			if count != 0 {
				count = 0
				toFitOneLine(lineArr)
			}
			fmt.Println(err)
			//删除文件
			err = os.Remove(newFileString) //删除文件test.txt
			if err != nil {
				//如果删除失败则输出 file remove Error!
				fmt.Println("file remove Error!")
				//输出错误详细信息
				fmt.Printf("%s", err)
			} else {
				//如果删除成功则输出 file remove OK!
				fmt.Print("file remove OK!")
			}
			break
		} else {
			if !isPrefix {
				if count == 500 {
					toFitOneLine(lineArr)
					//传入后，那么就清0
					count = 0
					lineArr.Init()
				}
				//将一个元素推入list
				lineArr.PushBack(string(line))
				count++
			}
		}

	}
	fmt.Println("finish", i)
	defer w.Done()
}

func toFitOneLine(fitlines *list.List) {
	//原来的表
	//mysqlArrList := []myfunc.MysqlVisitDate{}
	//新表
	mysqlArrList := []myfunc.Accesslogss{}
	for onefitline := fitlines.Front(); onefitline != nil; onefitline = onefitline.Next() {
		fitline := fmt.Sprintf("%s", onefitline.Value)
		ip1Reg := `(?P<ip1>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
		ip2Reg := `(?P<ip2>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
		ip3Reg := `(?P<io3>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
		ip4Reg := `(?P<ip4>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
		dateReg := `(?P<local_time>\d{1,2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})`
		requestUrl := `(?P<requestUrl>.*)`
		statuscode := `(?P<statuscode>\d+)`
		bodysize := `(?P<bodysize>\d+)`
		fromUrl := `(?P<fromUrl>.*)`
		agent := `(?P<agent>.*)`
		requesttime := `(?P<requesttime>.*)`
		replaceString := fmt.Sprintf(`%s - %s - %s, %s - - \[%s\] \"%s\" %s %s \"%s\" \"%s\" \"%s\"`, ip1Reg, ip2Reg, ip3Reg, ip4Reg, dateReg, requestUrl, statuscode, bodysize, fromUrl, agent, requesttime)
		reg := regexp.MustCompile(replaceString)
		match := reg.FindStringSubmatch(fitline)
		//如果不匹配表达式,那么记录这一行的数据
		if len(match) == 0 {
			wf, err := os.OpenFile("F:/nglog/nopasstheurl.txt", os.O_APPEND, 0775)
			if err != nil {
				return
			}
			_, err1 := io.WriteString(wf, fitline+"\n")
			if err1 != nil {
				fmt.Println("can not write file")
			}
		} else {
			//fmt.Println(len(match))
			//调用mongodb
			//myfunc.InsertIngoMongodb(match)
			//调用mysql
			//分割match[6] request_url
			request_url := strings.Split(match[6], " ")
			scode, _ := strconv.Atoi(match[7])
			sbody, _ := strconv.Atoi(match[8])
			sRtime, _ := strconv.ParseFloat(match[11], 32)
			//原来的表
			//mysqlArrList = append(mysqlArrList, myfunc.MysqlVisitDate{match[1], match[5], request_url[0], request_url[1], request_url[2], scode, sbody, match[9], match[10], sRtime})
			//新表-------------
			//新表增加的字段
			ipCountry, ipProvince, ipCity := myfunc.Convertip_tiny(match[1])
			plat := myfunc.Platforms(match[10])
			browser := myfunc.Browsers(match[10])
			mobile := myfunc.Mobiles(match[10])
			date_Reg, _ := myfunc.UtcTimeToNormalDateTime(match[5])
			mysqlArrList = append(mysqlArrList, myfunc.Accesslogss{match[1], ipCountry, ipProvince, ipCity, date_Reg, request_url[0], request_url[1], request_url[2], scode, sbody, match[9], match[10], plat, browser, mobile, sRtime})
		}
	}
	myfunc.InsertIndb(mysqlArrList)
}
