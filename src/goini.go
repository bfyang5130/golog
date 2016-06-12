package main

import (
	"fmt"

	"github.com/widuu/goini"
)

func main() {
	conf := goini.SetConfig("./config/config.ini")    //goini.SetConfig(filepath) filepath = directory+file
	username := conf.GetValue("database", "username") //username is your key you want get the value
	fmt.Println(username)                             //root
	conf.DeleteValue("database", "username")          //username is your delete the key
	conf.SetValue("database", "username", "widuu")
	username = conf.GetValue("database", "username")
	fmt.Println(username) //widuu Adding/section configuration information if there is to add or modify the value of the corresponding, if there is no add section
}
