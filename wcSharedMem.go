package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// Ways to count file content: by lines, by words and by bytes
const WAYS = 3

var sharedVar int

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

	// Silly implementation: the goroutine writes result to the sharedVar
	// then to the counts slice
	// the counts slice is used to display the results in the required order
	var counts []int
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(WAYS)
	go func() {
		defer wg.Done()
		mutex.Lock()
		sharedVar = strings.Count(string(str), "\n")
		counts = append(counts, sharedVar)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex.Lock()
		sharedVar = len(strings.Fields(string(str)))
		counts = append(counts, sharedVar)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex.Lock()
		sharedVar = len(string(str))
		counts = append(counts, sharedVar)
		mutex.Unlock()
	}()
	wg.Wait()

	slices.Sort(counts)
	for _, count := range counts {
		fmt.Print(count, " ")
	}
	fmt.Printf("%s\n", arguments[1])
}
