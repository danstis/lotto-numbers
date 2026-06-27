// Package generator provides lottery number generation functionality.
package generator

import (
	"fmt"
	"math/rand"
	"sort"
)

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

	pool := make([]int, len(numbersList))
	copy(pool, numbersList)

	// Dedupe to get the required distinct values.
	seenMap := make(map[int]bool, len(pool))
	required := make([]int, 0, len(pool))
	for _, n := range pool {
		if !seenMap[n] {
			seenMap[n] = true
			required = append(required, n)
		}
	}
	m := len(required)

	// Feasibility: M distinct numbers must fit across lines * numPerLine slots.
	if m > lines*numPerLine {
		return nil
	}
	// Too few distinct values to fill even one line.
	if m < numPerLine {
		return nil
	}

	// Safety cap for line-uniqueness retries. Coverage is the hard guarantee;
	// uniqueness is best-effort within this bound.
	maxAttempts := lines*(len(pool)+1)*2 + 1

	// buildOnce constructs one full set of lines. When enforceUniqueness is true
	// it rejects duplicate lines; when false it skips that check so coverage is
	// always returned even if some lines are identical.
	buildOnce := func(enforceUniqueness bool) ([][]int, bool) {
		rand.Shuffle(m, func(i, j int) {
			required[i], required[j] = required[j], required[i]
		})

		// Assign required values round-robin so every required number lands in
		// exactly one line and no line receives more than ceil(M/lines) values.
		lineSlots := make([][]int, lines)
		for idx, n := range required {
			lineSlots[idx%lines] = append(lineSlots[idx%lines], n)
		}

		result := make([][]int, 0, lines)
		usedKeys := make(map[string]bool)

		for i := range lines {
			inLine := make(map[int]bool, numPerLine)
			for _, n := range lineSlots[i] {
				inLine[n] = true
			}

			rand.Shuffle(len(pool), func(a, b int) {
				pool[a], pool[b] = pool[b], pool[a]
			})

			line := make([]int, len(lineSlots[i]), numPerLine)
			copy(line, lineSlots[i])

			for _, n := range pool {
				if len(line) == numPerLine {
					break
				}
				if !inLine[n] {
					inLine[n] = true
					line = append(line, n)
				}
			}

			if len(line) != numPerLine {
				return nil, false
			}

			sort.Ints(line)

			if enforceUniqueness {
				key := fmt.Sprint(line)
				if usedKeys[key] {
					return nil, false
				}
				usedKeys[key] = true
			}

			result = append(result, line)
		}
		return result, true
	}

	for range maxAttempts {
		if res, ok := buildOnce(true); ok {
			return res
		}
	}

	// Uniqueness could not be achieved within the cap; return coverage-guaranteed
	// lines without the uniqueness constraint.
	res, ok := buildOnce(false)
	if !ok {
		return nil
	}
	return res
}

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
