package main

import (
	"fmt"
	"regexp"
)

func main() {
	var fitline string
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
	//fitline = `60.2.123.235 - 111.202.74.142 - 60.2.123.235, 111.202.74.142 - - [18/Apr/2016:00:01:05 +0800] "GET /ajaxCross/Reg.ashx?Cmd=GetVoiceVerifyCode&mobile=18330323015 HTTP/1.1" 200 5 "https://www.tuandai.com/ajaxCross/Reg.ashx?Cmd=GetVoiceVerifyCode&mobile=18330323015" "Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1; 125LA; .NET CLR 2.0.50727; .NET CLR 3.0.04506.648; .NET CLR 3.5.21022)" "0.000"`
	fitline = `60.2.123.235 - 111.202.74.142 - 60.2.123.235, 111.202.74.142 - - [18/Apr/2016:00:01:05 +0800] "GET /ajaxCross/Reg.ashx?Cmd=GetVoiceVerifyCode&mobile=18330323015 HTTP/1.1" 200 5 "https://www.tuandai.com/ajaxCross/Reg.ashx?Cmd=GetVoiceVerifyCode&mobile=18330323015" "Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1; 125LA; .NET CLR 2.0.50727; .NET CLR 3.0.04506.648; .NET CLR 3.5.21022)" "0.000"`
	replaceString := fmt.Sprintf(`%s - %s - %s, %s - - \[%s\] \"%s\" %s %s \"%s\" \"%s\" \"%s\"`, ip1Reg, ip2Reg, ip3Reg, ip4Reg, dateReg, requestUrl, statuscode, bodysize, fromUrl, agent, requesttime)
	//replaceString := `\[(\d{1,2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})\]`
	reg := regexp.MustCompile(replaceString)
	match := reg.FindStringSubmatch(fitline)
	fmt.Println(len(match))
	for index := 0; index < len(match); index++ {
		fmt.Println(match[index])
	}
	/**
	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		// result[name] = match[i]
		if i != 0 {
			result[name] = match[i]
			fmt.Println(result[name])
		}
	}
	*/
}
