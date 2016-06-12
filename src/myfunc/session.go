package myfunc

import (
	"time"

	"database/sql"

	_ "database/sql"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	_ "github.com/go-sql-driver/mysql"
)

const URL = "192.168.8.188:27017"                                                       //mongodb连接字符串
const myURL = "root:tuandai1921688190@tcp(192.168.8.190:3306)/Tuandai_Log?charset=utf8" //mysql连接字符串
var (
	mgoSession *mgo.Session
	dataBase   = "nginx"
	mysqldb    *sql.DB
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GetMgoSession() (mgoSession *mgo.Session, err error) {
	if mgoSession == nil {
		maxWait := time.Duration(5 * time.Second)
		mgoSession, err = mgo.DialWithTimeout(URL, maxWait)
		if err != nil {
			return mgoSession, err
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone(), nil
}

type VisitDate struct {
	Id          bson.ObjectId `bson:"_id"`
	Ip1         string        `bson:"Ip1"`
	DateReg     string        `bson:"date_reg"`
	RequestUrl  string        `bson:"request_url"`
	StatusCode  int           `bson:"status_code"`
	BodySize    int           `bson:"body_size"`
	FromUrl     string        `bson:"from_url"`
	Agent       string        `bson:"agent"`
	RequestTime float64       `bson:"request_time"`
}

func GetMysqlSession() (mysqldb *sql.DB, err error) {
	if mysqldb == nil {
		mysqldb, err = sql.Open("mysql", myURL)
		if err != nil {
			return mysqldb, err
		}
	}
	//最大连接池默认为4096
	return mysqldb, nil
}
