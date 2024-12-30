package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Ways to count file content: by lines, by words and by bytes
const WAYS = 3

// Channel to write values
var WriteVal = make(chan int)

// Channel to read values
var ReadVal = make(chan int)

// Function to set val value of the data variable
func set(val int) {
	WriteVal <- val
}

// Function to get the last value of the data variable
func get() int {
	return <-ReadVal
}

// Function that implements the control goroutine
func monitor() {
	var data int
	for {
		select {
		case value := <-WriteVal:
			data = value
			fmt.Print(data, " ")
		case ReadVal <- data:
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		log.Fatalf("Need a file to process!\nUsage: %s filename\n", filepath.Base(arguments[0]))
	}

	file, err := os.Open(arguments[1])

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	str, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	go monitor()

	// Signaling channels to synchronize
	// the execution sequence of goroutines
	a := make(chan struct{})
	b := make(chan struct{})
	c := make(chan struct{})
	d := make(chan struct{})
	e := make(chan struct{})
	f := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1 + WAYS)
	go func(a, b chan struct{}) {
		<-a
		set(strings.Count(string(str), "\n"))
		close(b)
		wg.Done()
	}(a, b)

	go func(a, b chan struct{}) {
		<-a
		set(len(strings.Fields(string(str))))
		close(b)
		wg.Done()
	}(c, d)

	go func(a, b chan struct{}) {
		<-a
		set(len(str))
		close(b)
		wg.Done()
	}(e, f)

	// Goroutine to get results
	go func() {
		<-b
		get()
		close(c)
		<-d
		get()
		close(e)
		<-f
		get()
		wg.Done()
	}()

	close(a)
	wg.Wait()
	fmt.Printf("%s\n", arguments[1])
}
