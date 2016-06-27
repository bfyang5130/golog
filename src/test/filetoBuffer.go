package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"myfunc"
	"net"
	"os"
	"strconv"
	"strings"
)

/**
 *读取一个文件到内存并操作内存的功能
 */
func main() {

	fitline := `2016-06-02 16:00:00 192.168.1.41 POST /pages/ajax/newinvest_list.ashx - 20022 - 111.40.197.7 Mozilla/5.0+(Windows+NT+6.1;+WOW64;+rv:21.0)+Gecko/20100101+Firefox/21.0 200 0 0 31`
	arr := strings.Split(fitline, ` `)
	for k, a := range arr {
		fmt.Println(a, `---------`, k)
	}
	return
	//
	datetime, errr11 := myfunc.UtcTimeToNormalDateTime(`18/Apr/2016:00:00:02 +0800`)
	if errr11 != nil {
		fmt.Println(errr11)
	}
	fmt.Println(datetime)
	//相对路径测试
	fmt.Println(os.Getwd())
	//os.Chdir("D:/")
	fmt.Println(os.Getwd())
	//读取文件，获得文件的大小
	//var ipCountry, ipProvince, ipCity string
	var ip = `183.186.40.209`
	var fileIpPath = `../logFit/source/tinyipdata.dat`
	//获得IP的标记地址
	newBinary1 := binary.BigEndian.Uint32(net.ParseIP(ip).To4())
	//将IP标记地址转化为16进制用来做下面的比较
	newBinary := fmt.Sprintf(`%x`, newBinary1)
	//把IP地址分开，后面用他们组合的16进制与IP标记地址做比较得到相应数据
	//这里主要用第一个节点IP做标记
	ipdot := strings.Split(ip, `.`)
	fi, err := os.Open(fileIpPath)
	if err != nil {
		fmt.Println(`IP文件不能打开`)
		return
	}
	info, err := fi.Stat()
	if err != nil {
		fmt.Println(`IP文件有问题`)
		return
	}
	defer fi.Close()
	fise := info.Size()
	//定义一个buffer用来放这个文件
	buf := make([]byte, fise)
	n, err := fi.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println(`IP文件大小有误`, n)
		return
	}
	a, b, c := myfunc.Convertip_CachTiny(ip, buf)
	fmt.Println(a, b, c)
	return
	offset_len := binary.BigEndian.Uint32(buf[0:4])
	fmt.Println(offset_len)
	//处理标记信息，并配置相关查询条件
	ipdot1, err1 := strconv.Atoi(ipdot[0])
	if err1 != nil {
		//处理错误的转换
	}

	start := ipdot1*4 + 4
	start_lens := binary.LittleEndian.Uint32(buf[start : start+4])
	fmt.Println(start_lens)

	//配置最大长度
	length := offset_len - 1028
	var index_offset uint32
	var index_length int
	//循环处理数据
	for newstart := start_lens*8 + 1024; newstart < length; newstart += 8 {
		snb := fmt.Sprintf("%x", string(buf[newstart+4:newstart+4+4]))
		if snb > newBinary {
			//找到相应数据后，得到相应标记的文件地址，然后读取数据
			newb := make([]byte, 4)
			newb[0] = buf[newstart+4+4 : newstart+5+4][0]
			newb[1] = buf[newstart+5+4 : newstart+6+4][0]
			newb[2] = buf[newstart+6+4 : newstart+7+4][0]
			newb[3] = 0x00
			index_offset = binary.LittleEndian.Uint32(newb)
			index_len := fmt.Sprintf(`%d`, buf[newstart+7+4:newstart+8+4])
			index_len = strings.Replace(index_len, `[`, ``, -1)
			index_len = strings.Replace(index_len, `]`, ``, -1)
			index_length, _ = strconv.Atoi(index_len)
			break
		}

	}
	//文件指针跳转
	readindex := offset_len + index_offset - 1024
	fmt.Println(readindex)
	fmt.Println(index_length)
	if index_length > 0 {
		fmt.Println(string(buf[readindex : readindex+uint32(index_length)]))
	}
	//return ipCountry, ipProvince, ipCity
}
