package myfunc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func fitOffsetFile(fi string, start int64, end int64) (fitStatus bool, fitSize int64, err error) {
	fitfi, err := os.Open(fi)
	if err != nil {
		return false, start, err
	}
	//重新设定文件偏移量
	fitfi.Seek(start, 0)
	//创建一个读入器读入文件
	br := bufio.NewReader(fitfi)
	//记录处理到的位置
	var makeSize int64
	makeSize = start
	for {
		line, isPrefix, err := br.ReadLine()
		if err == io.EOF {
			return true, end, nil
			break
		} else {
			if !isPrefix {
				//返回终止的位移
				//return inOffset + len(string(line)), nil
				fmt.Println(string(line))
				makeSize = makeSize + int64(len(string(line)))
				//break
			}
		}
	}
	if makeSize == end {
		return true, end, nil
	} else {
		return false, makeSize, errors.New("not to the end")
	}

}
