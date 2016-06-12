package main

import (
	"fmt"
	"myfunc"
)

func main() {
	country, province, city := myfunc.Convertip_tiny(`218.21.128.51`)
	fmt.Println(fmt.Sprintf(`%s%s%s`, country, province, city))
}
