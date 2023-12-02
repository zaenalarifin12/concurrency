package main

import "testing"

func Test_updateMessage(t *testing.T) {
	msg = "Hello, world"

	wg.Add(1)
	go updateMessage("Goodbye, cruel world!")
	wg.Wait()

	if msg != "Goodbye, cruel world" {
		t.Error("Incorrect value in message")
	}
}
