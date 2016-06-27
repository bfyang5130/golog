package myfunc

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/DannyBen/filecache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/widuu/goini"
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

func NewInsertIndb(matchs []Accesslogss, mydb *sqlx.DB, logtype string) {
	var InsertDBString string
	switch logtype {
	case `nginx`:
		InsertDBString = "INSERT INTO AccessLogss (Ip1,country,province,city, date_reg, request_method,request_url,request_protocol,status_code, body_size, from_url, agent,plat,bower,mobile_plat, request_time,visitwebsite,proxy,collectimtime) VALUES"
	case `iis`:
		InsertDBString = "INSERT INTO IisAccessLog (Ip1,country,province,city, date_reg, request_method,request_url,request_protocol,status_code, body_size, from_url, agent,plat,bower,mobile_plat, request_time,visitwebsite,proxy,collectimtime) VALUES"

	default:
		InsertDBString = "INSERT INTO AccessLogss (Ip1,country,province,city, date_reg, request_method,request_url,request_protocol,status_code, body_size, from_url, agent,plat,bower,mobile_plat, request_time,visitwebsite,proxy,collectimtime) VALUES"
	}
	bf := bytes.NewBufferString(InsertDBString)
	il := len(matchs)
	vs := make([]interface{}, 0)
	//fmt.Println(matchs)
	for k, oneMatch := range matchs {
		if k == (il - 1) {
			bf.WriteString("(?,?,?,?,?,?, ?,?, ?, ?,?, ?, ?,?,?,?,?,?,?);")
		} else {
			bf.WriteString("(?,?,?,?,?,?, ?,?, ?, ?,?, ?, ?,?,?,?,?,?,?),")
		}
		vs = append(vs, oneMatch.Ip1, oneMatch.Country, oneMatch.Province, oneMatch.City, oneMatch.DateReg, oneMatch.Request_method, oneMatch.Request_url, oneMatch.Request_protocol, oneMatch.StatusCode, oneMatch.BodySize, oneMatch.FromUrl, oneMatch.Agent, oneMatch.Plat, oneMatch.Bower, oneMatch.Mobile_plat, oneMatch.RequestTime, oneMatch.Visitwebsite, oneMatch.Proxy, oneMatch.Collectiontime)
	}
	if _, err := mydb.Exec(bf.String(), vs...); err != nil {
		fmt.Println(bf.String())
		panic(err)
	}
}

//原表是MysqlVisitDate新表是
func InsertIndb(matchs []Accesslogss, fcaches filecache.Handler, fixconf *goini.Config, configName string) {
	//从配置文件拿到链接
	var mysqllinks string

	mysqllinks = fixconf.GetValue("dblink", `mysqllinks`)
	if len(mysqllinks) == 0 || mysqllinks == "no value" {
	} else {

		//获得一个连接session
		db, err1 := sql.Open("mysql", mysqllinks)
		if err1 != nil {
			checkErr(err1, "db open fail", matchs[0].Marksize, fcaches, configName)
			return
		}
		defer db.Close()
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		//新表
		stmt, err := db.Prepare(`INSERT INTO AccessLogss (Ip1,country,province,city, date_reg, request_method,request_url,request_protocol,status_code, body_size, from_url, agent,plat,bower,mobile_plat, request_time,visitwebsite,proxy,collectimtime) VALUES (?,?,?,?,?,?, ?,?, ?, ?,?, ?, ?,?,?,?,?,?,?)`)
		checkErr(err, "str prepare", matchs[0].Marksize, fcaches, configName)
		for _, oneMatch := range matchs {
			fmt.Println("-------------------------------")
			//新表
			res, err := stmt.Exec(oneMatch.Ip1, oneMatch.Country, oneMatch.Province, oneMatch.City, oneMatch.DateReg, oneMatch.Request_method, oneMatch.Request_url, oneMatch.Request_protocol, oneMatch.StatusCode, oneMatch.BodySize, oneMatch.FromUrl, oneMatch.Agent, oneMatch.Plat, oneMatch.Bower, oneMatch.Mobile_plat, oneMatch.RequestTime, oneMatch.Visitwebsite, oneMatch.Proxy, oneMatch.Collectiontime)
			checkErr(err, "str Exec fail", oneMatch.Marksize, fcaches, configName)
			num, err := res.RowsAffected()
			checkErr(err, "get row fail", oneMatch.Marksize, fcaches, configName)
			fmt.Println(num)
		}
	}
}

//原表是MysqlVisitDate新表是
func InsertIisIndb(matchs []Accesslogss, fcaches filecache.Handler, fixconf *goini.Config, configName string) {
	//从配置文件拿到链接
	var mysqllinks string

	mysqllinks = fixconf.GetValue("dblink", `mysqllinks`)
	if len(mysqllinks) == 0 || mysqllinks == "no value" {
	} else {

		//获得一个连接session
		db, err1 := sql.Open("mysql", mysqllinks)
		if err1 != nil {
			checkErr(err1, "db open fail", matchs[0].Marksize, fcaches, configName)
			return
		}
		defer db.Close()
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10) //新表
		stmt, err := db.Prepare(`INSERT INTO IisAccessLog (Ip1,country,province,city, date_reg, request_method,request_url,request_protocol,status_code, body_size, from_url, agent,plat,bower,mobile_plat, request_time,visitwebsite,proxy,collectimtime) VALUES (?,?,?,?,?,?, ?,?, ?, ?,?, ?, ?,?,?,?,?,?,?)`)
		checkErr(err, "str prepare", matchs[0].Marksize, fcaches, configName)
		for _, oneMatch := range matchs {
			fmt.Println("-------------------------------")
			//新表
			res, err := stmt.Exec(oneMatch.Ip1, oneMatch.Country, oneMatch.Province, oneMatch.City, oneMatch.DateReg, oneMatch.Request_method, oneMatch.Request_url, oneMatch.Request_protocol, oneMatch.StatusCode, oneMatch.BodySize, oneMatch.FromUrl, oneMatch.Agent, oneMatch.Plat, oneMatch.Bower, oneMatch.Mobile_plat, oneMatch.RequestTime, oneMatch.Visitwebsite, oneMatch.Proxy, oneMatch.Collectiontime)
			checkErr(err, "str Exec fail", oneMatch.Marksize, fcaches, configName)
			num, err := res.RowsAffected()
			checkErr(err, "get row fail", oneMatch.Marksize, fcaches, configName)
			fmt.Println(num)
		}
	}
}
func checkErr(err error, errstr string, marksize int64, fcaches filecache.Handler, configName string) {
	if err != nil {
		fmt.Println(marksize)
		//记录新的marksize地址
		fcaches.Set(configName, []byte(fmt.Sprintf(`%d`, marksize)))
		fmt.Println(marksize)
		panic(err)
	}
}

/**
Accesslogs
*/
type MysqlVisitDate struct {
	Ip1              string
	DateReg          string
	Request_method   string
	Request_url      string
	request_protocol string
	StatusCode       int
	BodySize         int
	FromUrl          string
	Agent            string
	RequestTime      float64
}

/**
Accesslogss
*/
type Accesslogss struct {
	Ip1              string
	Country          string
	Province         string
	City             string
	DateReg          string
	Request_method   string
	Request_url      string
	Request_protocol string
	StatusCode       int
	BodySize         int
	FromUrl          string
	Agent            string
	Plat             string
	Bower            string
	Mobile_plat      string
	RequestTime      float64
	Visitwebsite     string
	Proxy            int
	Collectiontime   string
	Marksize         int64
}

/**
*一个用来存储500的行数据
 */
type LineMessage struct {
	Line             string
	FileSizePosition int64
}
