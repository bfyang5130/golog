package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 100)
	go produce(ch)
	for i := 0; i < 100; i++ {
		go customer(ch)
	}
	time.Sleep(1 * time.Second)
}

func produce(p chan<- string) {
	for i := 0; i <= 100; i++ {
		p <- fmt.Sprintf(`%d`, i)
		fmt.Println(`send:`, i)
	}

}
func customer(ch <-chan string) {
	for i := 0; i < 200; i++ {
		v := <-ch
		fmt.Println(`receive`, v)
	}

}
