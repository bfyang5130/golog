package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	//"time"
)

func main() {
	stringsarray := strings.Split(`1 2 3 4`, ` `)

	//t := time.Now()
	//nt := t.Format(`20060102`)
	fmt.Println(stringsarray[1])
	//ReadLine(`F:/nglog/test2.log`)
}
func ReadLine(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 4*1024)
	line, isPrefix, err := r.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		fmt.Println(s)
		line, isPrefix, err = r.ReadLine()
	}
	if isPrefix {
		fmt.Println("buffer size to small")
		return
	}
	if err != io.EOF {
		fmt.Println(err)
		return
	}
}
