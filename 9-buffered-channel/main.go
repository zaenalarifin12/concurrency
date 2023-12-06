package main

import (
	"fmt"
	"time"
)

func listenToChan(ch chan int) {
	for {
		// print a got data message
		i := <-ch
		fmt.Println("Got ", i, "from channel")

		// simulate doing of a lot work
		time.Sleep(1 * time.Second)
	}
}
func main() {
	ch := make(chan int)

	go listenToChan(ch)

	for i := 0; i < 101; i++ {
		fmt.Println("Sending ", i, "to channel...")
		ch <- i
		fmt.Println("Sent ", i, "to channel!")
	}

	fmt.Println("Done!")
	close(ch)
}
