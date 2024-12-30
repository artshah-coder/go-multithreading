package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// The struct to store job info for worker:
// file content as string value
// split function
type Job struct {
	Content   string
	SplitFunc bufio.SplitFunc
}

// The struct to store the result of the J job
type Result struct {
	J     Job
	count int
}

var (
	sz   = 10
	jobs = make(chan Job, sz)    // buffered channel to store jobs
	res  = make(chan Result, sz) // buffered channel to store results
	// Split functions:
	sFs = []bufio.SplitFunc{bufio.ScanLines, bufio.ScanWords, bufio.ScanBytes}
	// Ways to count file content: by lines, by words and by bytes:
	ways = len(sFs)
)

// Worker implementation
func worker(w *sync.WaitGroup) {
	for j := range jobs {
		nR := strings.NewReader(j.Content)
		scanner := bufio.NewScanner(nR)
		scanner.Split(j.SplitFunc)

		count := 0
		for scanner.Scan() {
			count++
		}
		res <- Result{j, count}
	}
	w.Done()
}

// Function to create workers pool
func makeWP(n int) {
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		w.Add(1)
		go worker(&w)
	}
	w.Wait()
	close(res)
}

// Function to create jobs queue
func create(s string) {
	for i := 0; i < ways; i++ {
		jobs <- Job{s, sFs[i]}
	}
	close(jobs)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Need filename and #workers!")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	str, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	nWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	go create(string(str))

	finished := make(chan interface{})
	var counts []int
	go func() {
		for r := range res {
			counts = append(counts, r.count)
		}
		slices.Sort(counts)
		finished <- true
	}()
	makeWP(nWorkers)
	<-finished
	for _, c := range counts {
		fmt.Print(c, " ")
	}
	fmt.Printf("%s\n", os.Args[1])
}
