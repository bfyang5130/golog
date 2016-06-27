// read_large_file project main.go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type AC struct {
	Id            int
	UserIP1       string
	UserIP2       string
	UserIP3       string
	UserIP4       string
	RequestTime   time.Time
	RequestType   string
	Protocol      string
	AccessAddress string
	Status        int
	ContentSize   int64
	HttpReferer   string
	ClientType    string
	System        string
	Browser       string
	TakeTime      float64
	access_type   string
	source        string
}

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

var (
	wg   sync.WaitGroup
	reg  *regexp.Regexp
	regX *regexp.Regexp
	regD *regexp.Regexp
	mydb *sqlx.DB
)

const (
	workers = 15
	oneLot  = 10
)

func init() {
	if a1, err := regexp.Compile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`); err != nil {
		panic(err)
	} else {
		reg = a1
	}
	if a1, err := regexp.Compile(`"(.*?)"`); err != nil {
		panic(err)
	} else {
		regX = a1
	}
	if a1, err := regexp.Compile(`"\s+\d{3}\s+\d+\s+"`); err != nil {
		panic(err)
	} else {
		regD = a1
	}
}

//117.68.247.121 - 60.174.241.91 - 117.68.247.121, 60.174.241.91 - - [06/Jun/2016:00:00:03 +0800] "GET /view/ip.aspx HTTP/1.1" 200 14 "-" "Apache-HttpClient/UNAVAILABLE (java 1.4)" "0.002"
//分析行
func parseLine(data []byte) (*AC, error) {
	ac := &AC{}
	line := string(data)
	ix := strings.Index(line, "[")
	ixEnd := strings.Index(line, "]")
	ips := reg.FindAllString(line, ix)
	date := line[ix : ixEnd+1]
	otherS := line[ixEnd+1:]
	for k, v := range ips {
		if k == 0 {
			ac.UserIP1 = v
		}
		if k == 1 {
			ac.UserIP2 = v
		}
		if k == 2 {
			ac.UserIP3 = v
		}
		if k == 3 {
			ac.UserIP4 = v
		}
	}
	fs := regX.FindAllString(otherS, -1)
	if len(fs) == 4 {
		ac.AccessAddress = strings.Replace(fs[0], `"`, "", -1)
		ac.HttpReferer = strings.Replace(fs[1], `"`, "", -1)
		brs := strings.Replace(fs[2], `"`, "", -1)
		if len(brs) > 100 {
			ac.Browser = brs[0:100]
		} else {
			ac.Browser = brs
		}
		vvv := strings.Replace(fs[3], `"`, "", -1)
		if tt, err := strconv.ParseFloat(vvv, 10); err == nil {
			ac.TakeTime = tt
		} else {
			fmt.Println(err.Error(), vvv)
		}
	}
	//	RFC822      = "02 Jan 06 15:04 MST"
	//06/Jun/2016:00:00:26
	//2006-01-02 15:04:05.999999999
	dstr := date[1 : len(date)-6]
	if t1, err := time.Parse("02/Jan/2006:15:04:05", strings.TrimSpace(dstr)); err == nil {
		ac.RequestTime = t1
	} else {
		fmt.Println(dstr, err.Error())
	}
	//	fmt.Println(ips, ">>>")
	//	fmt.Println("len=", len(fs), " >>", fs)

	sts1 := regD.Find([]byte(otherS))
	if len(sts1) > 0 {
		ss := strings.TrimSpace(strings.Replace(string(sts1), `"`, "", -1))
		slist := strings.Split(ss, " ")
		if len(slist) >= 2 {
			if iv, err := strconv.Atoi(slist[0]); err != nil {
				fmt.Println(slist[0], err.Error())
			} else {
				ac.Status = iv
			}

			if iv, err := strconv.ParseInt(slist[1], 10, 64); err != nil {
				fmt.Println(slist[1], err.Error())
			} else {
				ac.ContentSize = iv
			}
		}

	}

	return ac, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//	flag.Parse()
	filename := "F:/nglog/Log_source/nginx_proxy_log/17_nginx_proxy/20160614/www.tuandai.com.access.log" // flag.Arg(0)
	//	filename := "a.log"
	chLine := make(chan []byte, 1000)

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

	wg.Add(1)
	go readfile(filename, chLine)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(chLine)
	}

	wg.Wait()
	fmt.Println("all ok")
}

func worker(ch <-chan []byte) {
	icount := 0

	acs := make([]*AC, 0, oneLot)
	for {
		select {
		case oneline, ok := <-ch:
			if ok {
				ac, err := parseLine(oneline)
				if err != nil {
					fmt.Println(err)
				} else {
					icount++
					acs = append(acs, ac)
					if icount >= oneLot {
						if err := savedb(acs[:]); err != nil {
							fmt.Println(err)
						} else {
							icount = 0
							acs = acs[0:0]
							//copy(acs, acs[0:0])
							//							fmt.Println("clear -->", len(acs), cap(acs))
						}
					}
				}
			} else {
				//保存剩余的
				if err := savedb(acs[:]); err != nil {
					fmt.Println(err)
				}

				fmt.Println("chan close")
				wg.Done()
				return
			}
		default:
			//time.After(time.Second * 1)
		}
	}
}

//save to db
func savedb(acs []*AC) error {
	il := len(acs)
	if il > 0 {
		bf := bytes.NewBufferString("insert into AccessLog(UserIP1,UserIP2,UserIP3,UserIP4,RequestTime,AccessAddress,HttpReferer,Browser,TakeTime,Status,ContentSize)values")
		vs := make([]interface{}, 0)
		for k, v := range acs {
			if k == (il - 1) {
				bf.WriteString("(?,?,?,?,?,?,?,?,?,?,?);")
			} else {
				bf.WriteString("(?,?,?,?,?,?,?,?,?,?,?),")
			}
			vs = append(vs, v.UserIP1, v.UserIP2, v.UserIP3, v.UserIP4, v.RequestTime, v.AccessAddress, v.HttpReferer, v.Browser, v.TakeTime, v.Status, v.ContentSize)
		}

		if _, err := mydb.Exec(bf.String(), vs...); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

//读取日志文件
func readfile(filename string, ch chan []byte) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
		close(ch)
		wg.Done()
	}()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		ch <- line
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
