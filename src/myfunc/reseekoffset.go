package myfunc

import (
	"bufio"
	"io"
	"os"
)

func reseek(infile string, inOffset int64) (offset int64, err error) {
	fitfi, err := os.Open(infile)
	if err != nil {
		return 0, err
	}
	defer fitfi.Close()
	finfo, err := fitfi.Stat()
	if err != nil {
		return 0, err
	}
	//如果已经超过最大的位置了，那么返回最大位移的位置
	maxSize := finfo.Size()
	if err != nil {
		return 0, err
	}
	if inOffset > maxSize {
		return maxSize, nil
	}
	//重新设定文件偏移量
	fitfi.Seek(inOffset, 0)
	br := bufio.NewReader(fitfi)
	for {
		line, isPrefix, err := br.ReadLine()
		if err == io.EOF {
			return 0, err
			break
		} else {
			if !isPrefix {
				//返回终止的位移
				return inOffset + int64(len(string(line))), nil
				break
			}
		}
	}
	return 0, nil
}
