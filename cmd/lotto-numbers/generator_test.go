package main

import (
	"sort"
	"testing"
)

func TestGenerateLotteryNumbers(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	// Test for correct number of lines and numbers per line
	if len(generatedNumbers) != lines {
		t.Errorf("Expected %d lines, got %d", lines, len(generatedNumbers))
	}
	for _, line := range generatedNumbers {
		if len(line) != numPerLine {
			t.Errorf("Expected %d numbers per line, got %d", numPerLine, len(line))
		}

		// Test for uniqueness of numbers in each line
		numMap := make(map[int]bool)
		for _, num := range line {
			if numMap[num] {
				t.Errorf("Duplicate number found in a line: %d", num)
			}
			numMap[num] = true

			// Test that number is in numbersList
			if !contains(numbersList, num) {
				t.Errorf("Number %d in generated line is not in the original numbersList", num)
			}
		}
	}
}

func contains(slice []int, item int) bool {
	sort.Ints(slice)
	index := sort.SearchInts(slice, item)
	return index < len(slice) && slice[index] == item
}
