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

var data = make(chan int, WAYS)

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		log.Fatalf("Need a file to process!\nUsage: %s filename\n", filepath.Base(arguments[0]))
	}

	f, err := os.Open(arguments[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	str, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// Signaling channels to synchronize
	// the execution sequence of goroutines
	x := make(chan struct{})
	y := make(chan struct{})
	z := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(WAYS)
	go func(a, b chan struct{}) {
		defer wg.Done()
		<-a
		data <- strings.Count(string(str), "\n")
		close(b)
	}(x, y)

	go func(a, b chan struct{}) {
		defer wg.Done()
		<-a
		data <- len(strings.Fields(string(str)))
		close(b)
	}(y, z)

	go func(a chan struct{}) {
		defer wg.Done()
		<-a
		data <- len(str)
		close(data)
	}(z)

	close(x)
	wg.Wait()

	for d := range data {
		fmt.Printf("%d ", d)
	}
	fmt.Printf("%s\n", arguments[1])
}
