package myfunc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

/**
 *将nginx的时间变成YYYY-MM-DD H:i:s的形式
 */
func UtcTimeToNormalDateTime(timeString string) (dateTime string, err error) {
	shortMonthNames := map[string]string{
		`Jan`: `01`,
		`Feb`: `02`,
		`Mar`: `03`,
		`Apr`: `04`,
		`May`: `05`,
		`Jun`: `06`,
		`Jul`: `07`,
		`Aug`: `08`,
		`Sep`: `09`,
		`Oct`: `10`,
		`Nov`: `11`,
		`Dec`: `12`,
	}
	//为空就返回错误
	dateTime = ``
	err = errors.New("wrong string")
	//第一部分切割
	timeStep := strings.Split(timeString, " ")
	if len(timeStep) < 2 {
		return dateTime, err
	}
	//第二部分切割
	timeStep2 := strings.Split(timeStep[0], `/`)
	if len(timeStep2) < 2 {
		return dateTime, err
	}
	day := timeStep2[0]
	month := shortMonthNames[timeStep2[1]]

	//第三部分切割
	timeStep3 := strings.Split(timeStep2[2], `:`)
	year := timeStep3[0]
	hour := timeStep3[1]
	minute := timeStep3[2]
	sec := timeStep3[3]
	dateTime = fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hour, minute, sec)
	return dateTime, nil
}

/**
获得平台类型
*/
func Platforms(agent string) (plat string) {
	platNames := map[string]string{
		`windows nt 10.0`: `Windows 10`,
		`windows nt 6.3`:  `Windows 8.1`,
		`windows nt 6.2`:  `Windows 8`,
		`windows nt 6.1`:  `Windows 7`,
		`windows nt 6.0`:  `Windows Vista`,
		`windows nt 5.2`:  `Windows 2003`,
		`windows nt 5.1`:  `Windows XP`,
		`windows nt 5.0`:  `Windows 2000`,
		`windows nt 4.0`:  `Windows NT 4.0`,
		`winnt4.0`:        `Windows NT 4.0`,
		`winnt 4.0`:       `Windows NT`,
		`winnt`:           `Windows NT`,
		`windows 98`:      `Windows 98`,
		`win98`:           `Windows 98`,
		`windows 95`:      `Windows 95`,
		`win95`:           `Windows 95`,
		`windows phone`:   `Windows Phone`,
		`windows`:         `Unknown Windows OS`,
		`android`:         `Android`,
		`blackberry`:      `BlackBerry`,
		`iphone`:          `iOS`,
		`ipad`:            `iOS`,
		`ipod`:            `iOS`,
		`os x`:            `Mac OS X`,
		`ppc mac`:         `Power PC Mac`,
		`freebsd`:         `FreeBSD`,
		`ppc`:             `Macintosh`,
		`linux`:           `Linux`,
		`debian`:          `Debian`,
		`sunos`:           `Sun Solaris`,
		`beos`:            `BeOS`,
		`apachebench`:     `ApacheBench`,
		`aix`:             `AIX`,
		`irix`:            `Irix`,
		`osf`:             `DEC OSF`,
		`hp-ux`:           `HP-UX`,
		`netbsd`:          `NetBSD`,
		`bsdi`:            `BSDi`,
		`openbsd`:         `OpenBSD`,
		`gnu`:             `GNU/Linux`,
		`unix`:            `Unknown Unix OS`,
		`symbian`:         `Symbian OS`,
	}
	//将传参数据转成小写
	newAgent := strings.ToLower(agent)
	//遍历平台信息
	// 遍历map
	for k, _ := range platNames {
		//查找字符串
		if strings.Contains(newAgent, k) {
			plat = platNames[k]
			return plat
		}

	}
	plat = `other`
	return plat
}

/**
获得浏览器
*/
func Browsers(agent string) (brow string) {
	platNames := map[string]string{
		`OPR`:               `Opera`,
		`Flock`:             `Flock`,
		`Edge`:              `Spartan`,
		`Chrome`:            `Chrome`,
		`Opera`:             `Opera`,
		`MSIE`:              `Internet Explorer`,
		`Internet Explorer`: `Internet Explorer`,
		`Trident`:           `Internet Explorer`,
		`Shiira`:            `Shiira`,
		`Firefox`:           `Firefox`,
		`Chimera`:           `Chimera`,
		`Phoenix`:           `Phoenix`,
		`Firebird`:          `Firebird`,
		`Camino`:            `Camino`,
		`Netscape`:          `Netscape`,
		`OmniWeb`:           `OmniWeb`,
		`Safari`:            `Safari`,
		`Mozilla`:           `Mozilla`,
		`Konqueror`:         `Konqueror`,
		`icab`:              `iCab`,
		`Lynx`:              `Lynx`,
		`Links`:             `Links`,
		`hotjava`:           `HotJava`,
		`amaya`:             `Amaya`,
		`IBrowse`:           `IBrowse`,
		`Maxthon`:           `Maxthon`,
		`Ubuntu`:            `Ubuntu Web Browser`,
	}
	//将传参数据转成小写
	newAgent := strings.ToLower(agent)
	//遍历平台信息
	// 遍历map
	for k, _ := range platNames {
		//查找字符串
		if strings.Contains(newAgent, strings.ToLower(k)) {
			brow = platNames[k]
			return brow
		}

	}
	brow = `other`
	return brow
}

func Mobiles(agent string) (brow string) {
	platNames := map[string]string{
		`mobileexplorer`:       `Mobile Explorer`,
		`palmsource`:           `Palm`,
		`palmscape`:            `Palmscape`,
		`motorola`:             `Motorola`,
		`nokia`:                `Nokia`,
		`palm`:                 `Palm`,
		`iphone`:               `Apple iPhone`,
		`ipad`:                 `iPad`,
		`ipod`:                 `Apple iPod Touch`,
		`sony`:                 `Sony Ericsson`,
		`ericsson`:             `Sony Ericsson`,
		`blackberry`:           `BlackBerry`,
		`cocoon`:               `O2 Cocoon`,
		`blazer`:               `Treo`,
		`lg`:                   `LG`,
		`amoi`:                 `Amoi`,
		`xda`:                  `XDA`,
		`mda`:                  `MDA`,
		`vario`:                `Vario`,
		`htc`:                  `HTC`,
		`samsung`:              `Samsung`,
		`sharp`:                `Sharp`,
		`sie-`:                 `Siemens`,
		`alcatel`:              `Alcatel`,
		`benq`:                 `BenQ`,
		`ipaq`:                 `HP iPaq`,
		`mot-`:                 `Motorola`,
		`playstation portable`: `PlayStation Portable`,
		`playstation 3`:        `PlayStation 3`,
		`playstation vita`:     `PlayStation Vita`,
		`hiptop`:               `Danger Hiptop`,
		`nec-`:                 `NEC`,
		`panasonic`:            `Panasonic`,
		`philips`:              `Philips`,
		`sagem`:                `Sagem`,
		`sanyo`:                `Sanyo`,
		`spv`:                  `SPV`,
		`zte`:                  `ZTE`,
		`sendo`:                `Sendo`,
		`nintendo dsi`:         `Nintendo DSi`,
		`nintendo ds`:          `Nintendo DS`,
		`nintendo 3ds`:         `Nintendo 3DS`,
		`wii`:                  `Nintendo Wii`,
		`open web`:             `Open Web`,
		`openweb`:              `OpenWeb`,
		`android`:              `Android`,
		`symbian`:              `Symbian`,
		`SymbianOS`:            `SymbianOS`,
		`elaine`:               `Palm`,
		`series60`:             `Symbian S60`,
		`windows ce`:           `Windows CE`,
		`obigo`:                `Obigo`,
		`netfront`:             `Netfront Browser`,
		`openwave`:             `Openwave Browser`,
		`mobilexplorer`:        `Mobile Explorer`,
		`operamini`:            `Opera Mini`,
		`opera mini`:           `Opera Mini`,
		`opera mobi`:           `Opera Mobile`,
		`fennec`:               `Firefox Mobile`,
		`digital paths`:        `Digital Paths`,
		`avantgo`:              `AvantGo`,
		`xiino`:                `Xiino`,
		`novarra`:              `Novarra Transcoder`,
		`vodafone`:             `Vodafone`,
		`docomo`:               `NTT DoCoMo`,
		`o2`:                   `O2`,
		`mobile`:               `Generic Mobile`,
		`wireless`:             `Generic Mobile`,
		`j2me`:                 `Generic Mobile`,
		`midp`:                 `Generic Mobile`,
		`cldc`:                 `Generic Mobile`,
		`up.link`:              `Generic Mobile`,
		`up.browser`:           `Generic Mobile`,
		`smartphone`:           `Generic Mobile`,
		`cellphone`:            `Generic Mobile`,
	}
	//将传参数据转成小写
	newAgent := strings.ToLower(agent)
	//遍历平台信息
	// 遍历map
	for k, _ := range platNames {
		//查找字符串
		if strings.Contains(newAgent, strings.ToLower(k)) {
			brow = platNames[k]
			return brow
		}

	}
	brow = `other`
	return brow
}
func Convertip_tiny(ip string) (ipCountry, ipProvince, ipCity string) {
	//获得IP的标记地址
	newBinary1 := binary.BigEndian.Uint32(net.ParseIP(ip).To4())
	//将IP标记地址转化为16进制用来做下面的比较
	newBinary := fmt.Sprintf(`%x`, newBinary1)
	//把IP地址分开，后面用他们组合的16进制与IP标记地址做比较得到相应数据
	//这里主要用第一个节点IP做标记
	ipdot := strings.Split(ip, `.`)
	//fmt.Println(ipdot[0])
	//读取IP二进制文件
	fi, err := os.Open("/home/golangpath/src/logFit/source/tinyipdata.dat")
	//fi, err := os.Open("./source/tinyipdata.dat")

	if err != nil {
		//任务中不想出错，读不到再说吧
		fmt.Println(err)
		ipCountry, ipProvince, ipCity = ``, ``, ``
		return ipCountry, ipProvince, ipCity
	}
	//读取文件头部信息得到文件的基本信息用来做标记
	mode := make([]byte, 4)
	fi.Read(mode)
	offset_len := binary.BigEndian.Uint32(mode)
	//读取内容
	index := make([]byte, offset_len-4)
	fi.Read(index)
	//处理标记信息，并配置相关查询条件
	ipdot1, err1 := strconv.Atoi(ipdot[0])
	if err1 != nil {
		return ipCountry, ipProvince, ipCity
	}
	start := ipdot1 * 4
	start_lens := binary.LittleEndian.Uint32(index[start : start+4])
	//配置最大长度
	length := offset_len - 1028
	var index_offset uint32
	var index_length int
	//循环处理数据
	for newstart := start_lens*8 + 1024; newstart < length; newstart += 8 {
		snb := fmt.Sprintf("%x", string(index[newstart:newstart+4]))
		if snb > newBinary {
			//找到相应数据后，得到相应标记的文件地址，然后读取数据
			newb := make([]byte, 4)
			newb[0] = index[newstart+4 : newstart+5][0]
			newb[1] = index[newstart+5 : newstart+6][0]
			newb[2] = index[newstart+6 : newstart+7][0]
			newb[3] = 0x00
			index_offset = binary.LittleEndian.Uint32(newb)
			index_len := fmt.Sprintf(`%d`, index[newstart+7:newstart+8])
			index_len = strings.Replace(index_len, `[`, ``, -1)
			index_len = strings.Replace(index_len, `]`, ``, -1)
			index_length, _ = strconv.Atoi(index_len)
			break
		}

	}
	//文件指针跳转
	readindex := offset_len + index_offset - 1024
	fi.Seek(int64(readindex), 0)
	if index_length > 0 {
		newmode := make([]byte, index_length)
		fi.Read(newmode)
		ipCountry, ipProvince, ipCity = SplitToCPC(string(newmode))
	}
	return ipCountry, ipProvince, ipCity
}

//将一个地址分解成国家省份城市
func SplitToCPC(ctcString string) (ipCountry, ipProvince, ipCity string) {
	//如果有中国两个字,那么就是中国否则返回所有的数据记为一个国家
	if strings.Contains(ctcString, `中国`) {
		ipCountry = `中国`
	} else {
		ipCountry = ctcString
		ipProvince = ``
		ipCity = ``
		return ipCountry, ipProvince, ipCity
	}
	//把中国替换掉
	ctcString = strings.Replace(ctcString, `中国`, ``, -1)
	//截取前两个字符
	var newByte = []byte(ctcString)
	//判断中国的32
	ipProvince = string(newByte[0:6])
	if strings.Contains(ipProvince, `黑龙`) {
		ipProvince = `黑龙江`
	} else if strings.Contains(ipProvince, `内蒙`) {
		ipProvince = `内蒙古`
	}
	//最后一个余下的就是城市了
	ipCity = strings.Replace(ctcString, ipProvince, ``, -1)
	return ipCountry, ipProvince, ipCity
}

//检查并创建文件
func CheckFileSize(configfile string) (fileSize int64) {
	//如果不存在文件那么就创建一个文件，并且返回文件大小为0
	if !FileExist(configfile) {
		//创建文件
		err := os.MkdirAll(path.Dir(configfile), 0777)

		if err != nil {
			fmt.Println(err)
			return -1
		}
		f, err1 := os.Create(configfile)
		//如果创建成功那么就写入文件
		if err1 != nil {
			//返回一个错误的状态
			fmt.Println(err1)
			return -1
		}
		io.WriteString(f, "[filesize]")
	}

	fileSize = 1
	return fileSize
}

//判断文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//返回一个集合时间,我不会告诉你什么叫集合时间
func ToCollectiontime(date string) (coldate string) {
	//时间格式为2016-06-14 10:06:57
	if date == `` {
		coldate = `0000-00-00 00:00:00`
		return coldate
	}
	//将时间用:切分
	minuteTime := strings.Split(date, `:`)
	if len(minuteTime) < 1 {
		coldate = `0000-00-00 00:00:00`
		return coldate
	}
	marktime, err := strconv.Atoi(minuteTime[1])
	if err != nil {
		coldate = `0000-00-00 00:00:00`
		return coldate
	}
	switch {
	case marktime >= 50:
		minuteTime[1] = `50`
	case marktime >= 40:
		minuteTime[1] = `40`
	case marktime >= 30:
		minuteTime[1] = `30`
	case marktime >= 20:
		minuteTime[1] = `20`
	case marktime >= 10:
		minuteTime[1] = `10`
	default:
		minuteTime[1] = `00`
	}
	//组合coldate
	coldate = fmt.Sprintf(`%s:%s:00`, minuteTime[0], minuteTime[1])
	return coldate
}

func FitWebsite(confilefile string) (visitwebsite string) {
	website := strings.Split(confilefile, `/`)
	logname := website[len(website)-1]
	splitlog := strings.Split(logname, `.`)
	visitwebsite = strings.Replace(logname, splitlog[len(splitlog)-1], ``, -1)
	visitwebsite = strings.Replace(visitwebsite, splitlog[len(splitlog)-2], ``, -1)
	visitwebsite = strings.Replace(visitwebsite, `..`, ``, -1)
	return visitwebsite
}
