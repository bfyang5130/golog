package main

import (
	"fmt"
	"myfunc"
)

func main() {
	timestring := `18/Apr/2016:08:20:04 +0800`
	for i := 0; i < 100; i++ {
		the_time, err := myfunc.UtcTimeToNormalDateTime(timestring)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(the_time)
		}
	}
}
