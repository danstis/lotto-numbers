package main

import (
	"fmt"
	"math/rand"
	"sort"
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
	if len(numbersList) < numPerLine {
		return nil // Not enough numbers to generate a line
	}

	lotteryNumbers := make([][]int, 0, lines)
	linesMap := make(map[string]bool)
	for len(lotteryNumbers) < lines {
		rand.Shuffle(len(numbersList), func(i, j int) {
			numbersList[i], numbersList[j] = numbersList[j], numbersList[i]
		})
		uniqueLine := make(map[int]bool)
		line := make([]int, 0, numPerLine)
		for _, num := range numbersList[:numPerLine] {
			if !uniqueLine[num] {
				uniqueLine[num] = true
				line = append(line, num)
			}
		}
		if len(line) == numPerLine {
			sort.Ints(line) // Sort to normalize the line for comparison
			lineKey := fmt.Sprint(line)
			if !linesMap[lineKey] {
				linesMap[lineKey] = true
				lotteryNumbers = append(lotteryNumbers, line)
			}
		}
	}
	return lotteryNumbers
}
