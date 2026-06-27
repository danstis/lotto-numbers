package generator

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLotteryNumbers_CorrectNumberOfLines(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

	assert.Equal(t, lines, len(generatedNumbers), "Expected %d lines, got %d", lines, len(generatedNumbers))
}

func TestGenerateLotteryNumbers_CorrectNumbersPerLine(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	lines := 3
	numPerLine := 5

	generatedNumbers := GetNumbers(numbersList, lines, numPerLine)

	assert.NotEmpty(t, generatedNumbers, "No lines were generated, expected %d lines", lines)

	for _, line := range generatedNumbers {
		assert.Equal(t, numPerLine, len(line), "Expected %d numbers per line, got %d", numPerLine, len(line))
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

func TestGenerateLotteryNumbers_NotEnoughNumbers(t *testing.T) {
	tests := []struct {
		name        string
		numbersList []int
		lines       int
		numPerLine  int
		expected    [][]int
	}{
		{
			name:        "Not enough numbers for a single line",
			numbersList: []int{1, 2, 3}, // Only 3 numbers available
			lines:       1,
			numPerLine:  5, // Requires 5 numbers per line
			expected:    nil,
		},
		{
			name:        "Not enough numbers for multiple lines",
			numbersList: []int{1, 2, 3, 4}, // Only 4 numbers available
			lines:       2,
			numPerLine:  5, // Requires 5 numbers per line
			expected:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generatedNumbers := GetNumbers(tc.numbersList, tc.lines, tc.numPerLine)
			assert.Equal(t, tc.expected, generatedNumbers, "Test %s - Expected nil, got %v", tc.name, generatedNumbers)
		})
	}
}
func TestGetNumbers_UniformDistribution(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numPerLine := 5
	iterations := 10000

	freq := make(map[int]int)
	for range iterations {
		result := GetNumbers(numbersList, 1, numPerLine)
		if len(result) == 0 {
			t.Fatal("GetNumbers returned no results")
		}
		for _, num := range result[0] {
			freq[num]++
		}
	}

	// Each number should appear roughly iterations*numPerLine/len(numbersList) times.
	expected := float64(iterations*numPerLine) / float64(len(numbersList))
	// Allow ±15% deviation — chi-square at 99.9% for 9 df is ~27, this is a looser guard.
	tolerance := expected * 0.15
	for num, count := range freq {
		if math.Abs(float64(count)-expected) > tolerance {
			t.Errorf("number %d appeared %d times, expected ~%.0f (±%.0f)", num, count, expected, tolerance)
		}
	}
}

// TestGetNumbers_ConcurrentCallsDiverseResults is the regression test for the
// time-seeded RNG bug. With the old rand.New(rand.NewSource(time.Now().UnixNano()))
// pattern, goroutines that land in the same nanosecond receive identical seeds and
// thus produce identical shuffle sequences — most results collapse to the same
// combination. With the auto-seeded global rand the results must be diverse.
func TestGetNumbers_ConcurrentCallsDiverseResults(t *testing.T) {
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numPerLine := 5
	goroutines := 200

	type result struct{ line string }
	results := make(chan result, goroutines)

	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			r := GetNumbers(numbersList, 1, numPerLine)
			if len(r) > 0 {
				results <- result{fmt.Sprint(r[0])}
			}
		})
	}
	wg.Wait()
	close(results)

	seen := make(map[string]int)
	total := 0
	for r := range results {
		seen[r.line]++
		total++
	}

	// C(10,5) = 252 possible combinations. With 200 concurrent calls and a
	// properly seeded RNG we expect a wide spread; with a time-colliding seed
	// almost all calls collapse to the same combination.
	// Require at least 50 distinct combinations out of 200.
	if len(seen) < 50 {
		t.Errorf("concurrent calls produced only %d distinct combinations out of %d total; "+
			"expected ≥50 (possible RNG seed collision)", len(seen), total)
	}
}

func TestGetNumbers_DoesNotMutateInput(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	input := make([]int, len(original))
	copy(input, original)

	GetNumbers(input, 5, 5)

	assert.Equal(t, original, input, "GetNumbers must not modify the caller's slice")
}

func TestGetNumbersEnsureAll_AllNumbersAppearAtLeastOnce(t *testing.T) {
	// Exact-fit boundary: 24 distinct numbers into exactly 4×6 = 24 slots.
	// The only way to cover all 24 is a perfect partition, so the existing
	// uniform-random GetNumbers essentially never achieves this — making this
	// test a clear demonstrator of the ensure-all guarantee.
	numbersList := make([]int, 24)
	for i := range 24 {
		numbersList[i] = i + 1
	}
	result := GetNumbersEnsureAll(numbersList, 4, 6)
	assert.Len(t, result, 4, "expected 4 lines")

	seen := make(map[int]bool)
	for _, line := range result {
		for _, n := range line {
			seen[n] = true
		}
	}
	for _, n := range numbersList {
		assert.True(t, seen[n], "number %d must appear at least once across the generated lines", n)
	}
}

func TestGetNumbersEnsureAll_CorrectShapeAndUniqueWithinLine(t *testing.T) {
	numbersList := make([]int, 24)
	for i := range 24 {
		numbersList[i] = i + 1
	}
	result := GetNumbersEnsureAll(numbersList, 4, 6)
	assert.Len(t, result, 4, "expected 4 lines")

	for _, line := range result {
		assert.Len(t, line, 6, "each line must contain exactly 6 numbers")
		seen := make(map[int]bool)
		for _, n := range line {
			assert.False(t, seen[n], "duplicate number %d within a line", n)
			seen[n] = true
			assert.True(t, contains(numbersList, n), "number %d not in numbersList", n)
		}
	}
}

func TestGetNumbersEnsureAll_Infeasible_ReturnsNil(t *testing.T) {
	// 25 distinct numbers cannot fit into 4 lines × 6 = 24 slots.
	numbersList := make([]int, 25)
	for i := range 25 {
		numbersList[i] = i + 1
	}
	result := GetNumbersEnsureAll(numbersList, 4, 6)
	assert.Nil(t, result, "expected nil for infeasible input (25 > 4*6)")
}

func TestGetNumbersEnsureAll_FeasibleWithSlack_ReturnsCoverage(t *testing.T) {
	// 12 required numbers across 24 slots — random-fill path is exercised.
	numbersList := make([]int, 12)
	for i := range 12 {
		numbersList[i] = i + 1
	}
	result := GetNumbersEnsureAll(numbersList, 4, 6)
	assert.Len(t, result, 4, "expected 4 lines")

	seen := make(map[int]bool)
	for _, line := range result {
		assert.Len(t, line, 6, "each line must contain exactly 6 numbers")
		for _, n := range line {
			seen[n] = true
		}
	}
	for _, n := range numbersList {
		assert.True(t, seen[n], "number %d must appear at least once", n)
	}
}

func TestGetNumbersEnsureAll_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numbersList []int
		lines       int
		numPerLine  int
	}{
		{
			name:        "pool smaller than numPerLine",
			numbersList: []int{1, 2, 3},
			lines:       1,
			numPerLine:  5,
		},
		{
			name:        "zero lines",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       0,
			numPerLine:  5,
		},
		{
			name:        "negative lines",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       -1,
			numPerLine:  5,
		},
		{
			name:        "zero numPerLine",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       3,
			numPerLine:  0,
		},
		{
			name:        "negative numPerLine",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       3,
			numPerLine:  -1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Nil(t, GetNumbersEnsureAll(tc.numbersList, tc.lines, tc.numPerLine))
		})
	}
}

func TestGetNumbersEnsureAll_DoesNotMutateInput(t *testing.T) {
	original := make([]int, 24)
	for i := range 24 {
		original[i] = i + 1
	}
	input := make([]int, len(original))
	copy(input, original)

	GetNumbersEnsureAll(input, 4, 6)

	assert.Equal(t, original, input, "GetNumbersEnsureAll must not modify the caller's slice")
}

func TestGenerateLotteryNumbers_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numbersList []int
		lines       int
		numPerLine  int
		expected    [][]int
	}{
		{
			name:        "Zero lines",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       0,
			numPerLine:  5,
			expected:    nil,
		},
		{
			name:        "Negative lines",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       -1,
			numPerLine:  5,
			expected:    nil,
		},
		{
			name:        "Negative numbers per line",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       3,
			numPerLine:  -1,
			expected:    nil,
		},
		{
			name:        "Zero numbers per line",
			numbersList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			lines:       3,
			numPerLine:  0,
			expected:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generatedNumbers := GetNumbers(tc.numbersList, tc.lines, tc.numPerLine)
			assert.Equal(t, tc.expected, generatedNumbers)
		})
	}
}
