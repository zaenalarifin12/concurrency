package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	wg.Add(1)

	go updateMessage("nine")

	wg.Wait()

	if msg != "nine" {
		t.Error("Expected to find nine, but it is not there")
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "nine"
	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "nine") {
		t.Error("Expected to find nine, but it is not there")
	}
}

func Test_main(t *testing.T) {

	stdOut := os.Stdout

	w, r, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if strings.Contains(output, "Hello, Universe") {
		t.Error("Expected to find Hello, Universe, but it is not there")
	}

	if strings.Contains(output, "Hello, Cosmos") {
		t.Error("Expected to find Hello, Cosmos, but it is not there")
	}

	if strings.Contains(output, "Hello, Earth") {
		t.Error("Expected to find Hello, Earth, but it is not there")
	}
}
