package myfunc

import (
	"strconv"
	//  "bufio"
	"fmt"
	"io"
	"math"
	"os"
)

/**
传入文件file 传入分割文件数splintNums 返回是否分割成功splitStatus
*/
func SplitFile(infile string, splitFileSize int64) (splitStatus bool, splintNums int64, alreadyfitSize int64) {

	//读取文件
	file, err := os.Open(infile)
	//
	if err != nil {
		fmt.Println("failed to open")
		return false, 0, 0
	}

	defer file.Close()
	//得到文件信息
	finfo, err := file.Stat()
	if err != nil {
		fmt.Println("异常的文件")
		return false, 0, 0
	}
	//将文件分成30M一份的数据
	var splitN float64
	splitN = float64(finfo.Size()-alreadyfitSize) / float64(splitFileSize)
	//splitN = float64(thisFileSize) / float64(15)
	//成为一份整数
	splitN = math.Ceil(splitN)
	//循环处理分割每一份文件
	//设置第一份文件的开始读取点
	var fitFileStartMarkOffset int64
	fitFileStartMarkOffset = alreadyfitSize
	var eSize int64
	var j int64
	for j = 1; j <= int64(splitN); j++ {
		fmt.Println(fitFileStartMarkOffset)
		//获得对应分割文件的起点和终点
		eSize, err = reseek(infile, int64(j*splitFileSize))
		fmt.Println(eSize)
		if err != nil {
			fmt.Println("处理结束分割位置时出错")
			return false, 0, 0
		}
		//得到每份文件分割的大小
		size := eSize - fitFileStartMarkOffset
		if size < 1 {
			continue
		}
		//每次最多拷贝1m
		var bufsize int64
		bufsize = 1024 * 1024
		if size < bufsize {
			bufsize = size
		}
		buf := make([]byte, bufsize)
		//生成新文件的名字
		newfilename := infile + strconv.FormatInt(j, 10)
		//创建新文件
		newfile, err1 := os.Create(newfilename)
		if err1 != nil {

			fmt.Println("failed to create file", newfilename)
			return false, 0, 0
		} else {
			fmt.Println("create file:", newfilename)
		}
		//设置文件偏移位移
		file.Seek(fitFileStartMarkOffset, 0)

		var copylen int64
		copylen = 0
		for copylen < size {
			//如果最后一次读取超过设置的最后位置就设置为正确的位置
			if (copylen + bufsize) > size {
				bufsize = size - copylen
				buf = make([]byte, bufsize)
			}

			n, err2 := file.Read(buf)
			if err2 != nil && err2 != io.EOF {
				fmt.Println(err2, "failed to read from:", file)
				return false, 0, 0
				break
			}

			if n <= 0 {
				break
			}

			//写文件
			w_buf := buf[:n]
			newfile.Write(w_buf)
			copylen += int64(n)
		}
		//设置下一个文件的起点
		fitFileStartMarkOffset = eSize + 1
	}
	return true, int64(splitN), 0

}
