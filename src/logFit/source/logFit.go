package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"myfunc"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DannyBen/filecache"
	"github.com/jmoiron/sqlx"
	"github.com/widuu/goini"
)

var (
	w sync.WaitGroup
	//fcaches     = filecache.Handler{"/var/golang/src/logFit/cache", 43200}
	//logFileName = flag.String("log", "/var/golang/src/logFit/cServer.log", "Log file name")
	//fixconf     = goini.SetConfig(`/var/golang/src/logFit/config/config.ini`)
	mydb        *sqlx.DB
	fcaches     filecache.Handler
	logFileName = flag.String("log", "./cServer.log", "Log file name")
	fixconf     = goini.SetConfig(`./config/config.ini`)
)

//root:tuandai1921688190@tcp(192.168.8.190:3036)/Tuandai_Log
//root:jia123@tcp(127.0.0.1:3306)/godb
type DBConfig struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Database     string `json:"dbname"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

func main() {
	//获得当前时间，如果当前时间不是不大于20160617那么直接退出
	ifitmarktime := time.Now().Format("20060102")
	ifitmarktime1, ooooooerr := strconv.Atoi(ifitmarktime)
	if ooooooerr != nil {
		return
	}
	if ifitmarktime1 < 20160618 {
		fmt.Println(`未到时间`)
		return
	}
	//初始化配置文件，主要因为相对路径的问题
	fcaches = filecache.Handler{"./cache", 43200}
	fcaches.Set(`abce`)
	testAAA := fcaches.Get(`abce`)
	fmt.Println(testAAA)
	return
	dbc := DBConfig{
		UserName:     "root",
		Password:     "tuandai1921688190",
		IP:           "192.168.8.190",
		Port:         3306,
		Database:     "Tuandai_Log",
		MaxIdleConns: 6,
		MaxOpenConns: 6,
	}
	if db, err := dbInit(dbc); err != nil {
		panic(err)
	} else {
		mydb = db
	}
	//设置文件目录
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//配置IP文件
	ipBuf, err := myfunc.ReadIpDatToBuf(`../logFit/source/tinyipdata.dat`)
	if err != nil {
		fmt.Println(`IP文件配置不正确`)
		fmt.Println(err)
		return
	}

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
	//获得不符合条件时，保存的文件路径
	nomatchfilepath := fixconf.GetValue(`nomatchfile`, `nomatchfilepath`)
	if len(nomatchfilepath) == 0 || nomatchfilepath == "no value" {
		fmt.Println("未配置不符合规则记录的路径")
		return
	}
	//设置读取文件的目录,这里使用,把文件配置的目录分割起来。因为有多个目录的话，就可以配置啦。要不坑死我!
	filepathts := fixconf.GetValue(`filepath`, `logfilepath`)
	if len(filepathts) == 0 || filepathts == "no value" {
		fmt.Println("没有配置要处理的日志路径")
		return
	}
	//那么现在切割filepath
	filepaths := strings.Split(filepathts, `,`)
	//循环文件夹
	for filei := 0; filei < len(filepaths); filei++ {

		if strings.Contains(filepaths[filei], `nginx`) {
			fmt.Println(`lalalalalalala`)
			//获得当天正确的地址
			t := time.Now()
			nt := t.Format(`20060102`)
			fmt.Println(`get the read file dir`)
			newfilepath := fmt.Sprintf(`%s/%s`, filepaths[filei], nt)
			fmt.Println(newfilepath)
			list, err := ioutil.ReadDir(newfilepath)
			if err != nil {
				fmt.Println("Wrong Dir:", filepaths[filei])
				continue
			}
			fmt.Println(`start to reading the dir:`, newfilepath)
			//判断目录的类型,日志目录分有iis和nginx
			//判断是否存在nginx
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
				//配置最后读取位置的参数
				fitLastReadPosition := fmt.Sprintf(`%s/%s`, newfilepath, configName)
				//获得最后处理到的位置
				rfileSizes := fcaches.Get(fitLastReadPosition)
				var rfileSize string
				if rfileSizes == nil {
					//没有记录位置设置初始位置
					fcaches.Set(fitLastReadPosition, []byte(`0`))
					rfileSize = `0`
				} else {
					rfileSize = string(rfileSizes)
				}
				var fileSize int64
				if len(rfileSize) != 0 {
					fileSize, err = strconv.ParseInt(rfileSize, 10, 64)
				} else {
					fileSize = 0
				}
				fmt.Println(`fitfilemessage:`, fitLastReadPosition)
				fmt.Println(`-------------old position-------------`)
				fmt.Println(fileSize)
				fmt.Println(`-------------new position-------------`)
				fmt.Println(info.Size())
				fmt.Println(`-------------old-------------`)
				if fileSize >= info.Size() {
					fmt.Println("the same file size,not need to fix")
					continue
				}
				//定义一个通道用来装数据
				fitChannel := make(chan string)
				//分一个线程去读取文件
				//因为传进来的configName是一个完整的路径，所以通过完整的路径来判断代理和站点
				//如果字符串包括17_nginx_proxy
				var proxy int
				if strings.Contains(configName, `17_nginx_proxy`) {
					proxy = 17
				} else {
					proxy = 21
				}
				//处理站点visitwebsite := `www.tuandai.com`

				visitwebsite := myfunc.FitWebsite(configName)
				w.Add(1)
				go goReadOneFile(fitChannel, fileSize, fitLastReadPosition, fitLastReadPosition)
				//go toFitOneFile(fileSize, fitLastReadPosition, fitLastReadPosition, nomatchfilepath, `nginx`)
				//分15个子线程去处理文件
				for xC := 0; xC < 15; xC++ {
					w.Add(1)
					go goFitOneFile(fitChannel, nomatchfilepath, `nginx`, proxy, visitwebsite, mydb, ipBuf)
				}
				//等待多线程处理一个文件，在处理下一个文件之前进行堵塞
				w.Wait()

			}
		} else {
			//这部分处理的是iis的日志，这里只需要处理昨天的文件就行了
			//获得当天正确的地址
			t := time.Now()
			nt := t.Format(`060102`)

			//当前的文件名就是
			configName := fmt.Sprintf(`u_ex%s.log`, nt)
			newfilepath := fmt.Sprintf(`%s/%s`, filepaths[filei], configName)
			//configname其实并不能用来作为cache的区别，因为他们喜欢起相同的名字，所以我就加上路径区别吧
			//获得最后处理到的位置
			rfileSizes := fcaches.Get(newfilepath)
			var rfileSize string
			if rfileSizes == nil {
				//没有记录位置设置初始位置
				fcaches.Set(newfilepath, []byte(`0`))
				rfileSize = `0`
			} else {
				rfileSize = string(rfileSizes)
			}
			var fileSize int64

			if len(rfileSize) != 0 {
				fileSize, _ = strconv.ParseInt(rfileSize, 10, 64)
			} else {
				fileSize = 0
			}
			//获得文件的信息
			info, err22 := os.Stat(newfilepath)
			if err22 != nil {
				fmt.Println(`a wrong readfile:`, newfilepath)
				continue
			}
			fmt.Println(info.Name())
			fmt.Println(`-------------new-------------`)
			fmt.Println(fileSize)
			fmt.Println(`-------------new-------------`)
			fmt.Println(info.Size())
			fmt.Println(`-------------old-------------`)
			if fileSize >= info.Size() {
				fmt.Println("the same file size,not need to fix")
				continue
			}
			//分为多个线程去处理每一个文件
			w.Add(1)
			go toFitOneFile(fileSize, newfilepath, newfilepath, nomatchfilepath, `iis`)
			w.Wait()
		}

	}
}

/**
 *读取指定的文件到一个通道里，让其它线程从此通道分开处理数据
 */
func goReadOneFile(fitChannel chan<- string, fileSize int64, filepath string, configName string) {
	//定义一个用来记录处理到什么位置的标识，用来记录下次从什么位置开始处理文件
	var fitMarkPosition int64
	fitMarkPosition = fileSize
	//读取文件
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer func() {
		//如果出错或者处理完成，记录处理到的位置
		fcaches.Set(configName, []byte(fmt.Sprintf(`%d`, fitMarkPosition)))
		f.Close()
		close(fitChannel)
		w.Done()
	}()
	//把文件指针指向传入的文件fileSize,其实就是上一次处理到的位置
	//创建一个reader
	if fileSize > 0 {
		f.Seek(fileSize+1, 0)
	}
	//如果说配置参数出现异常
	r := bufio.NewReader(f)
	line, err1 := r.ReadString('\n')
	for err1 == nil {
		//将行line放入通道
		fitChannel <- line
		//fmt.Println(1)
		line, err1 = r.ReadString('\n')

	}
	//读完文件后，记录处理的到的位置
	//如果中途退出了，那么不会走到这一步的记录标识。
}

/**
 *多线程从通道获得数据，并进行处理后进行保存
 */
func goFitOneFile(fitChannel <-chan string, nomatchfilepath string, logtype string, proxy int, visitwebsite string, mydb *sqlx.DB, ipBuf *[]byte) {

	//原来的表
	//mysqlArrList := []myfunc.MysqlVisitDate{}
	//新表
	mysqlArrList := []myfunc.Accesslogss{}
	var line string
	ok := true
	icount := 0
	for ok {
		if line, ok = <-fitChannel; ok {
			//如果mysqlArrList有500个,那么就做入库处理
			//fmt.Println(line)
			if icount == 10 {
				icount = 0
				if len(mysqlArrList) > 0 {
					myfunc.NewInsertIndb(mysqlArrList, mydb)
				}
			}
			fitline := line
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
			if len(match) < 12 {

				nomatchfilepath = fixconf.GetValue("fitfile", nomatchfilepath)
				if len(nomatchfilepath) == 0 || nomatchfilepath == "no value" {
				} else {
					wf, err := os.OpenFile(nomatchfilepath, os.O_APPEND, 0775)
					if err != nil {
						return
					}
					_, err1 := io.WriteString(wf, fitline+"\n")
					if err1 != nil {
						fmt.Println("can not write file")
					}
				}
			} else {
				request_url := strings.Split(match[6], " ")
				scode, _ := strconv.Atoi(match[7])
				sbody, _ := strconv.Atoi(match[8])
				sRtime, _ := strconv.ParseFloat(match[11], 32)
				ipCountry, ipProvince, ipCity := myfunc.Convertip_CachTiny(match[1], ipBuf)
				plat := myfunc.Platforms(match[10])
				browser := myfunc.Browsers(match[10])
				mobile := myfunc.Mobiles(match[10])
				date_Reg, _ := myfunc.UtcTimeToNormalDateTime(match[5])
				//处理时间为collectiontime,就是一个集合时间
				collectiontime := myfunc.ToCollectiontime(date_Reg)
				mysqlArrList = append(mysqlArrList, myfunc.Accesslogss{match[1], ipCountry, ipProvince, ipCity, date_Reg, request_url[0], request_url[1], request_url[2], scode, sbody, match[9], match[10], plat, browser, mobile, sRtime, visitwebsite, proxy, collectiontime, 1})
			}
			icount++
		} else {
			//如果读完了，只要mysqlArrList>0也要入库处理
			if len(mysqlArrList) > 0 {
				myfunc.NewInsertIndb(mysqlArrList, mydb)
			}
			fmt.Println(`close channl`)
			w.Done()
		}
	}

}

//每一个线程来处理一个文件,fileSize已经处理到的地方，filepath文件路径configName记录处理到位置的文件名，nomatchfilepath,不匹配正则的记录，logtype日志分为iis,nginx要分开处理
func toFitOneFile(fileSize int64, filepath string, configName string, nomatchfilepath string, logtype string) {
	//定义一个用来记录处理到什么位置的标识，用来记录下次从什么位置开始处理文件
	var fitMarkPosition int64
	fitMarkPosition = fileSize
	//读取文件
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("A Wrong File")
		fmt.Println("finish", filepath)
		defer w.Done()
		return
	}
	defer f.Close()
	//把文件指针指向传入的文件fileSize,其实就是上一次处理到的位置
	//创建一个reader
	if fileSize > 0 {
		f.Seek(fileSize+1, 0)
	}
	count := 1
	//用array来装这500个数据，定义两个是因为不知道让第一个重置为空
	var one, two []myfunc.LineMessage
	//-----------------------------------------------------------------------------
	r := bufio.NewReader(f)
	line, err1 := r.ReadString('\n')
	for err1 == nil {
		if count == 500 {
			if logtype == `nginx` {
				toFitLines(one, configName, nomatchfilepath)
			} else {
				toFitIisLines(one, configName, nomatchfilepath)
			}
			//传入后，那么就清0
			count = 0
			one = two
		}
		//将一个元素推入list
		fitMarkPosition = fitMarkPosition + int64(len(line))
		one = append(one, myfunc.LineMessage{string(line), fitMarkPosition})
		fmt.Println(fitMarkPosition)
		count++
		line, err1 = r.ReadString('\n')

	}
	//-----------------------------------------------------------------------------
	//如果说one里还有数据,要再处理一次
	if len(one) > 0 {
		if logtype == `nginx` {
			toFitLines(one, configName, nomatchfilepath)
		} else {
			toFitIisLines(one, configName, nomatchfilepath)
		}
	}
	//如果中途退出了，那么不会走到这一步的记录标识。
	fcaches.Set(configName, []byte(fmt.Sprintf(`%d`, fitMarkPosition)))
	w.Done()
}

/**
处理一堆的行，上面的count定义为500，那么这里就是处理500行的数据
*/
func toFitLines(fitlines []myfunc.LineMessage, configName string, nomatchfilepath string) {
	//因为传进来的configName是一个完整的路径，所以通过完整的路径来判断代理和站点
	//如果字符串包括17_nginx_proxy
	var proxy int
	if strings.Contains(configName, `17_nginx_proxy`) {
		proxy = 17
	} else {
		proxy = 21
	}
	//处理站点visitwebsite := `www.tuandai.com`
	visitwebsite := myfunc.FitWebsite(configName)
	//原来的表
	//mysqlArrList := []myfunc.MysqlVisitDate{}
	//新表
	mysqlArrList := []myfunc.Accesslogss{}

	for _, onefitline := range fitlines {
		fitline := onefitline.Line
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
		if len(match) < 12 {

			nomatchfilepath = fixconf.GetValue("fitfile", `nomatchfilepath`)
			if len(nomatchfilepath) == 0 || nomatchfilepath == "no value" {
			} else {
				wf, err := os.OpenFile(nomatchfilepath, os.O_APPEND, 0775)
				if err != nil {
					return
				}
				_, err1 := io.WriteString(wf, fitline+"\n")
				if err1 != nil {
					fmt.Println("can not write file")
				}
			}
		} else {
			request_url := strings.Split(match[6], " ")
			scode, _ := strconv.Atoi(match[7])
			sbody, _ := strconv.Atoi(match[8])
			sRtime, _ := strconv.ParseFloat(match[11], 32)
			ipCountry, ipProvince, ipCity := myfunc.Convertip_tiny(match[1])
			plat := myfunc.Platforms(match[10])
			browser := myfunc.Browsers(match[10])
			mobile := myfunc.Mobiles(match[10])
			date_Reg, _ := myfunc.UtcTimeToNormalDateTime(match[5])
			//处理时间为collectiontime,就是一个集合时间
			collectiontime := myfunc.ToCollectiontime(date_Reg)
			mysqlArrList = append(mysqlArrList, myfunc.Accesslogss{match[1], ipCountry, ipProvince, ipCity, date_Reg, request_url[0], request_url[1], request_url[2], scode, sbody, match[9], match[10], plat, browser, mobile, sRtime, visitwebsite, proxy, collectiontime, onefitline.FileSizePosition})
		}
	}
	if len(mysqlArrList) > 0 {
		myfunc.InsertIndb(mysqlArrList, fcaches, fixconf, configName)
	}
}

/**
处理一堆的行，上面的count定义为500，那么这里就是处理500行的数据//iis的分析功能
*/
func toFitIisLines(fitlines []myfunc.LineMessage, configName string, nomatchfilepath string) {
	//因为传进来的configName是一个完整的路径，所以通过完整的路径来判断代理和站点
	//如果字符串包括17_nginx_proxy
	var proxy int
	if strings.Contains(configName, `41_iis`) {
		proxy = 41
	} else {
		proxy = 140
	}
	visitwebsite := `www.tuandai.com`
	//原来的表
	//mysqlArrList := []myfunc.MysqlVisitDate{}
	//新表
	mysqlArrList := []myfunc.Accesslogss{}

	for _, onefitline := range fitlines {
		fitline := onefitline.Line
		//判断第一个字符,如果第一个字符是#就是注释，直接跳过
		if strings.Contains(fitline, `#`) {
			continue
		}
		//iis的分析直接使用空格分开处理
		iisarrays := strings.Split(fitline, ` `)
		//开始配置各个分割出来对应的参数
		//ip
		ip1 := iisarrays[8]
		//ip解释
		ipCountry, ipProvince, ipCity := myfunc.Convertip_tiny(ip1)
		//时间
		date_Reg := fmt.Sprintf(`%s %s`, iisarrays[0], iisarrays[1])
		//处理时间为collectiontime,就是一个集合时间
		collectiontime := myfunc.ToCollectiontime(date_Reg)
		request_method := iisarrays[3]
		request_url := fmt.Sprintf(`%s?%s`, iisarrays[4], iisarrays[5])
		//IIS里没有这个协议的直接给空
		request_protocol := ``
		//状态值
		scode, _ := strconv.Atoi(iisarrays[10])
		sbody := 0
		sRtime, _ := strconv.ParseFloat(iisarrays[13], 32)
		//将agent中的+转为空格，然后分析
		newAgent := strings.Replace(iisarrays[9], `+`, ` `, -1)
		plat := myfunc.Platforms(newAgent)
		browser := myfunc.Browsers(newAgent)
		mobile := myfunc.Mobiles(newAgent)
		mysqlArrList = append(mysqlArrList, myfunc.Accesslogss{ip1, ipCountry, ipProvince, ipCity, date_Reg, request_method, request_url, request_protocol, scode, sbody, iisarrays[6], newAgent, plat, browser, mobile, sRtime, visitwebsite, proxy, collectiontime, onefitline.FileSizePosition})

	}
	if len(mysqlArrList) > 0 {
		myfunc.InsertIisIndb(mysqlArrList, fcaches, fixconf, configName)
	}
}

//创建db的连接
//db连接整个app中只需要一个，当网络断开时，再次请求会自动恢复的
//root:tuandai1921688190@tcp(192.168.8.190:3036)/Tuandai_Log
//root:jia123@tcp(127.0.0.1:3306)/godb
//jdbc:mysql://192.168.8.190/Tuandai_Log
func dbInit(dbconfig DBConfig) (*sqlx.DB, error) {
	myUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8", dbconfig.UserName, dbconfig.Password, dbconfig.IP, dbconfig.Port, dbconfig.Database)
	db, err := sqlx.Connect("mysql", myUrl)
	if err != nil {
		return nil, err
	}
	//设置缓存参数
	db.SetMaxOpenConns(dbconfig.MaxOpenConns)
	db.SetMaxIdleConns(dbconfig.MaxIdleConns)
	//不要设置下面这个值，连接的端口会定时变的
	///db.SetConnMaxLifetime(60 * time.Second)
	return db, nil
}