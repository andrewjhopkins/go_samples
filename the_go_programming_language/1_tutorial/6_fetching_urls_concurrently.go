package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	fetchAll()
}

func fetchAll() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetchUrl(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetchUrl(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf(err.Error())
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

// fetchall fetches URLs concurrently
// reports the size and elapsed time for each one

// goroutine is a concurrent function execution
// A channel is a communication mechanism that allows one goroutine to pass values of a specified type into another goroutine

// main function creates channel of strings using make
// for each command line argument the go statement in the frist range loop starts a new goroutine that calls fetch asynchronously

// when a goroutine attempts to send or recieve on a channel it blocks until another goroutine attempts the corresponding send or recieve
