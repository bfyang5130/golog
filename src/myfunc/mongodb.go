package myfunc

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertIngoMongodb(match []string) {
	//获得一个连接session
	session, err1 := GetMgoSession()
	if err1 != nil {
		fmt.Println("Unable to get the session")
		return
	}
	defer session.Close()
	scode, _ := strconv.Atoi(match[7])
	sbody, _ := strconv.Atoi(match[8])
	sRtime, _ := strconv.ParseFloat(match[11], 32)
	newPerson := VisitDate{
		bson.NewObjectId(),
		match[1],
		match[5],
		match[6],
		scode,
		sbody,
		match[9],
		match[10],
		sRtime}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("nginx").C("visit_date")
	err1 = c.Insert(newPerson)
	if err1 != nil {
		log.Fatal(err1)
	}
}
