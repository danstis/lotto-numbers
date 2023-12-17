package generator

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// GetNumbers generates a list of lottery numbers based on the given parameters.
//
// Parameters:
// - numbersList: a list of integers representing the available numbers for the lottery.
// - lines: an integer indicating the number of lines of lottery numbers to generate.
// - numPerLine: an integer indicating the number of numbers per line in the lottery.
//
// Returns:
// A 2D slice of integers representing the generated lottery numbers.
func GetNumbers(numbersList []int, lines, numPerLine int) [][]int {
	if len(numbersList) < numPerLine || lines <= 0 || numPerLine <= 0 {
		return nil // Not enough numbers to generate a line, zero lines requested, or non-positive numPerLine
	}

	lotteryNumbers := make([][]int, 0)
	linesMap := make(map[string]bool)
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lines; i++ {
		localRand.Shuffle(len(numbersList), func(i, j int) {
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
