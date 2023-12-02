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

func printMessage() {
	fmt.Println(msg)
}

func main() {

	msg = "Hello World"

	wg.Add(1)
	go updateMessage("Hello, Universe")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello, Cosmos")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello, Earth")
	wg.Wait()
	printMessage()

}
