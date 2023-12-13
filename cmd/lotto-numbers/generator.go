package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateLotteryNumbers generates a list of lottery numbers based on the given parameters.
//
// Parameters:
// - numbersList: a list of integers representing the available numbers for the lottery.
// - lines: an integer indicating the number of lines of lottery numbers to generate.
// - numPerLine: an integer indicating the number of numbers per line in the lottery.
//
// Returns:
// A 2D slice of integers representing the generated lottery numbers.
func generateLotteryNumbers(numbersList []int, lines, numPerLine int) [][]int {
func generateLotteryNumbers(numbersList []int, lines, numPerLine int) [][]int {
	if len(numbersList) < numPerLine {
		return nil // Not enough numbers to generate a line
	}

	lotteryNumbers := make([][]int, 0, lines)
	for i := 0; i < lines; i++ {
		rand.Shuffle(len(numbersList), func(i, j int) {
			numbersList[i], numbersList[j] = numbersList[j], numbersList[i]
		})
		line := make([]int, numPerLine)
		copy(line, numbersList[:numPerLine])
		lotteryNumbers = append(lotteryNumbers, line)
	}
	return lotteryNumbers
}
