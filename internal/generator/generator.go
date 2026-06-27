// Package generator provides lottery number generation functionality.
package generator

import (
	"fmt"
	"math/rand"
	"sort"
)

func cloneInts(values []int) []int {
	cloned := make([]int, len(values))
	copy(cloned, values)
	return cloned
}

func distinctInts(values []int) []int {
	seen := make(map[int]bool, len(values))
	distinct := make([]int, 0, len(values))
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		distinct = append(distinct, value)
	}
	return distinct
}

func buildLineFromPool(pool, seed []int, numPerLine int) ([]int, bool) {
	inLine := make(map[int]bool, numPerLine)
	for _, value := range seed {
		inLine[value] = true
	}

	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	line := make([]int, len(seed), numPerLine)
	copy(line, seed)
	for _, value := range pool {
		if len(line) == numPerLine {
			break
		}
		if inLine[value] {
			continue
		}
		inLine[value] = true
		line = append(line, value)
	}

	if len(line) != numPerLine {
		return nil, false
	}

	sort.Ints(line)
	return line, true
}

func assignRoundRobin(values []int, lines int) [][]int {
	lineSlots := make([][]int, lines)
	for idx, value := range values {
		lineSlots[idx%lines] = append(lineSlots[idx%lines], value)
	}
	return lineSlots
}

func buildEnsureAllResult(pool, required []int, lines, numPerLine int, enforceUniqueness bool) ([][]int, bool) {
	shuffledRequired := cloneInts(required)
	rand.Shuffle(len(shuffledRequired), func(i, j int) {
		shuffledRequired[i], shuffledRequired[j] = shuffledRequired[j], shuffledRequired[i]
	})

	result := make([][]int, 0, lines)
	usedKeys := make(map[string]bool, lines)
	for _, seed := range assignRoundRobin(shuffledRequired, lines) {
		line, ok := buildLineFromPool(pool, seed, numPerLine)
		if !ok {
			return nil, false
		}
		if !enforceUniqueness {
			result = append(result, line)
			continue
		}

		key := fmt.Sprint(line)
		if usedKeys[key] {
			return nil, false
		}
		usedKeys[key] = true
		result = append(result, line)
	}
	return result, true
}

// GetNumbersEnsureAll generates lottery numbers ensuring every distinct number in
// numbersList appears at least once across the generated lines.
//
// Coverage of all required numbers is the hard guarantee. Line uniqueness is
// best-effort within a bounded retry cap — if uniqueness cannot be reached within
// the cap, the covered lines are returned without uniqueness enforcement.
//
// Feasibility constraint: M (distinct count of numbersList) must satisfy
// M <= lines*numPerLine. Returns nil when this or the basic input constraints
// are violated.
func GetNumbersEnsureAll(numbersList []int, lines, numPerLine int) [][]int {
	if len(numbersList) < numPerLine || lines <= 0 || numPerLine <= 0 {
		return nil
	}

	pool := cloneInts(numbersList)
	required := distinctInts(pool)
	distinctCount := len(required)
	if distinctCount > lines*numPerLine || distinctCount < numPerLine {
		return nil
	}

	// Safety cap for line-uniqueness retries. Coverage is the hard guarantee;
	// uniqueness is best-effort within this bound.
	maxAttempts := lines*(len(pool)+1)*2 + 1
	for range maxAttempts {
		result, ok := buildEnsureAllResult(pool, required, lines, numPerLine, true)
		if ok {
			return result
		}
	}

	// Uniqueness could not be achieved within the cap; return coverage-guaranteed
	// lines without the uniqueness constraint.
	result, ok := buildEnsureAllResult(pool, required, lines, numPerLine, false)
	if ok {
		return result
	}
	return nil
}

// GetNumbers generates a list of lottery numbers based on the given parameters.
//
// Parameters:
// - numbersList: a list of integers representing the available numbers for the lottery.
// - lines: an integer indicating the number of lines of lottery numbers to generate.
// - numPerLine: an integer indicating the number of numbers per line in the lottery.
//
// Returns:
// A 2D slice of integers representing the generated numbers.
func GetNumbers(numbersList []int, lines, numPerLine int) [][]int {
	if len(numbersList) < numPerLine || lines <= 0 || numPerLine <= 0 {
		return nil // Not enough numbers to generate a line, zero lines requested, or non-positive numPerLine
	}

	pool := make([]int, len(numbersList))
	copy(pool, numbersList)

	lotteryNumbers := make([][]int, 0, lines)
	linesMap := make(map[string]bool)
	// Safety cap prevents an infinite loop when lines > C(len(pool), numPerLine).
	maxAttempts := lines * (len(pool) + 1) * 2
	for attempts := 0; len(lotteryNumbers) < lines && attempts < maxAttempts; attempts++ {
		rand.Shuffle(len(pool), func(i, j int) {
			pool[i], pool[j] = pool[j], pool[i]
		})
		uniqueLine := make(map[int]bool)
		line := make([]int, 0, numPerLine)
		for _, num := range pool[:numPerLine] {
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
