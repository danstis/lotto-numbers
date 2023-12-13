package main

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
	lotteryNumbers := make([][]int, lines)
	for i := 0; i < lines; i++ {
		line := make([]int, numPerLine)
		for j := 0; j < numPerLine; j++ {
			line[j] = numbersList[j]
		}
		lotteryNumbers[i] = line
	}
	return lotteryNumbers
}
