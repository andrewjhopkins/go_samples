package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strings"
)

const (
	NegativeColor = "\033[91m%q\033[39m"
	PositiveColor = "\033[92m%q\033[39m"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Correct usage: go run main.go <ticker symbol>")
		return
	}

	var price string
	var change string

	c := colly.NewCollector()

	c.OnHTML("tr td[class=wsod_last] span", func(e *colly.HTMLElement) {
		price = e.Text
	})

	c.OnHTML("tr td[class=wsod_change] span span", func(e *colly.HTMLElement) {
		if len(change) == 0 {
			change = e.Text
		} else {
			change += " / " + e.Text
		}
	})

	url := "https://money.cnn.com/quote/forecast/forecast.html?symb=" + os.Args[1]

	c.Visit(url)

	if len(price) == 0 && len(change) == 0 {
		fmt.Printf("Ticker symbol not found")
		return
	}

	var gain bool
	if strings.HasPrefix(change, "+") {
		gain = true
	} else {
		gain = false
	}

	fmt.Printf("Price: %s\n", addColor(price, gain))
	fmt.Printf("Change: %s\n", addColor(change, gain))
}

func addColor(input string, positive bool) string {
	if positive {
		return fmt.Sprintf(PositiveColor, input)
	} else {
		return fmt.Sprintf(NegativeColor, input)
	}
}
