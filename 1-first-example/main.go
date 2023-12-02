package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, group *sync.WaitGroup) {
	defer group.Done()
	fmt.Println(s)
}
func main() {
	var wg sync.WaitGroup

	words := []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}

	wg.Add(len(words))
	for i, word := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, word), &wg)
	}

	wg.Wait()
	wg.Add(1)

	printSomething("this is second", &wg)
}
