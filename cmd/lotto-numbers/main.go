package main

import (
	"log"

	"github.com/danstis/lotto-numbers/internal/version"
)

// Main entry point for the app.
func main() {
	log.Printf("Version %q", version.Version)

	num := generateLotteryNumbers([]int{1, 1, 1, 1, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29}, 5, 5)
	log.Printf("Numbers: %v", num)
}
