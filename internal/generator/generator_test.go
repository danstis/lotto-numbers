package generator

import (
	"fmt"
	"sort"
	"testing"
)

func TestGenerateLotteryNumbers_CorrectNumberOfLines(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) != lines {
		t.Errorf("Expected %d lines, got %d", lines, len(generatedNumbers))
	}
}

func TestGenerateLotteryNumbers_CorrectNumbersPerLine(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

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

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

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

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

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
func TestGenerateLotteryNumbers_UniqueLines(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	lines := 5
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

	if len(generatedNumbers) == 0 {
		t.Fatalf("No lines were generated, expected %d lines", lines)
	}

	linesMap := make(map[string]bool)
	for _, line := range generatedNumbers {
		sort.Ints(line) // Sort to normalize the line for comparison
		lineKey := fmt.Sprint(line)
		if linesMap[lineKey] {
			t.Fatalf("Duplicate line found: %v", line)
		}
		linesMap[lineKey] = true
	}
}
func TestGenerateLotteryNumbers_NoDuplicatesInLine(t *testing.T) {
	numbersList := []int{1, 1, 1, 1, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29}
	lines := 10
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

	for _, line := range generatedNumbers {
		seen := make(map[int]int)
		for _, num := range line {
			if seen[num] > 0 {
				t.Fatalf("Duplicate number %d found in line: %v", num, line)
			}
			seen[num]++
		}
	}
}
