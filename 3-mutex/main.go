package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()

	msg = s
}

func main() {
	msg = "hello, world"

	wg.Add(2)
	go updateMessage("hello, universe")
	go updateMessage("hello, cosmos")
	wg.Wait()

	fmt.Println(msg)
}

/**
tread safe operation
*/
// go run --race .
