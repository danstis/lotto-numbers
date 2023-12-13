package main

import (
	"sort"
	"testing"
)

func TestGenerateLotteryNumbers_CorrectNumberOfLines(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) != lines {
		t.Errorf("Expected %d lines, got %d", lines, len(generatedNumbers))
	}
}

func TestGenerateLotteryNumbers_CorrectNumbersPerLine(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) == 0 {
		t.Fatalf("No lines were generated, expected %d lines", lines)
	}

	for _, line := range generatedNumbers {
		if len(line) != numPerLine {
			t.Fatalf("Expected %d numbers per line, got %d", numPerLine, len(line))
		}
	}
}

func TestGenerateLotteryNumbers_UniqueNumbersInLine(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) == 0 {
		t.Fatalf("No lines were generated, expected %d lines", lines)
	}

	for _, line := range generatedNumbers {
		numMap := make(map[int]bool)
		for _, num := range line {
			if numMap[num] {
				t.Fatalf("Duplicate number found in a line: %d", num)
			}
			numMap[num] = true
		}
	}
}

func TestGenerateLotteryNumbers_NumbersFromList(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) == 0 {
		t.Fatalf("No lines were generated, expected %d lines", lines)
	}

	for _, line := range generatedNumbers {
		for _, num := range line {
			if !contains(numbersList, num) {
				t.Fatalf("Number %d in generated line is not in the original numbersList", num)
			}
		}
	}
}

func contains(slice []int, item int) bool {
	sort.Ints(slice)
	index := sort.SearchInts(slice, item)
	return index < len(slice) && slice[index] == item
}
