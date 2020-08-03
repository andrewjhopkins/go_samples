package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}

// go run helloworld.go

// go natively handles Unicode

// go build helloworld.go
// create binary executable

// go oragnized into packets similar to libraries
// consists of one or more .go source files in a single directory that defines what the package does
// package main defines the standalone executable
// where new lines are placed matters in go as they are converted into semicolons
// opening braces of functions must be on the same line as declaration
// use gofmt to rewrite code into standard format
