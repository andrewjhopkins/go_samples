package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	echo()
}

func echo() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Iterate over a range of values
// _, arg = index and value
func echo2() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

//Join from strings package
func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

// Print the name of the command that invoked it
func exercise1() {
	fmt.Println(os.Args[0])
}

// print the index and value of each of its arguments one per line
func exercise2() {
	for index, arg := range os.Args[1:] {
		fmt.Println(strconv.Itoa(index) + " " + arg)
	}
}

// Find difference in runtime between inefficient version and strings.Join
func exercise3() {
	start := time.Now()
	echo()
	fmt.Printf("%dms elapsed\n", time.Since(start).Microseconds())
	start = time.Now()
	echo3()
	fmt.Printf("%dms elapsed\n", time.Since(start).Microseconds())
}

// command line arguments available to a program in a variable called os.Args
// os.Args is a slice of strings
// a slice is a dynamically sized sequence s of array elements
// indexing in Go uses half open interavls that includes the first index but not the last
// os.Args[0] is the command itself
// + for string concatenation
// := short variable declaration that declares and gives appropriate types based on initializer value
// for is the only loop in go
// post can be omitted to simulate a while loop
// initialization, post and conditional can be omitted for an infinite loop
// _ blank identifiers used when syntax requires a variable name but program logic does not
