// Package handlers implements HTTP request handlers for the lotto-numbers API.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/danstis/lotto-numbers/internal/version"

	"github.com/danstis/lotto-numbers/internal/generator" // Import the generator package
	"github.com/danstis/lotto-numbers/internal/models"    // Import the models package
)

const maxLines = 100

// GetLotteryNumbers handles the request to generate lottery numbers.
func GetLotteryNumbers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	defaultNumbersList := make([]int, 40)
	for i := range 40 {
		defaultNumbersList[i] = i + 1
	}

	numbersList, err := parseNumbersListQueryParam(r)
	if err != nil {
		sendHTTPError(w, "Bad request", err, http.StatusBadRequest)
		return
	}
	if numbersList == nil {
		numbersList = defaultNumbersList
	}

	lines := parseQueryParamInt(r, "lines", 5)
	numPerLine := parseQueryParamInt(r, "numPerLine", 6)

	// Validate the parameters before calling the generator
	if lines <= 0 || numPerLine <= 0 {
		sendHTTPError(w, "Invalid input: 'lines' and 'numPerLine' must be positive numbers", nil, http.StatusBadRequest)
		return
	}
	if lines > maxLines {
		sendHTTPError(w, fmt.Sprintf("Invalid input: 'lines' must be <= %d", maxLines), nil, http.StatusBadRequest)
		return
	}

	ensureAll := parseQueryParamBool(r, "ensureAllNumbers", false)

	var generatedNumbers [][]int
	if ensureAll {
		// Compute distinct count M and enforce the feasibility rule before calling
		// the generator so the user gets a specific 400 rather than a generic 500.
		distinctCount := countDistinct(numbersList)
		if distinctCount > lines*numPerLine {
			sendHTTPError(w, fmt.Sprintf("cannot fit all %d numbers in %d lines of %d", distinctCount, lines, numPerLine), nil, http.StatusBadRequest)
			return
		}
		generatedNumbers = generator.GetNumbersEnsureAll(numbersList, lines, numPerLine)
	} else {
		generatedNumbers = generator.GetNumbers(numbersList, lines, numPerLine)
	}

	if generatedNumbers == nil {
		sendHTTPError(w, "Error generating numbers: ensure the 'numbersList' contains enough unique numbers", nil, http.StatusInternalServerError)
		return
	}

	// Create a LotteryNumbers struct with the generated numbers
	numbers := models.LotteryNumbers{
		Lines: generatedNumbers,
	}

	// Set the Content-Type header and encode the LotteryNumbers struct to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(numbers); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// VersionHandler writes the current application version to the response.
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(version.Version)); err != nil {
		log.Printf("Error writing version to response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// parseQueryParamBool parses a boolean query parameter, accepting "true" and "1" as true.
func parseQueryParamBool(r *http.Request, param string, defaultValue bool) bool {
	if value, ok := r.URL.Query()[param]; ok && len(value[0]) > 0 {
		v := strings.ToLower(value[0])
		return v == "true" || v == "1"
	}
	return defaultValue
}

// countDistinct returns the number of distinct values in the slice.
func countDistinct(nums []int) int {
	seen := make(map[int]bool, len(nums))
	for _, n := range nums {
		seen[n] = true
	}
	return len(seen)
}

// parseQueryParamInt parses an integer query parameter.
func parseQueryParamInt(r *http.Request, param string, defaultValue int) int {
	if value, ok := r.URL.Query()[param]; ok && len(value[0]) > 0 {
		if intValue, err := strconv.Atoi(value[0]); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// parseNumbersListQueryParam parses the numbersList query parameter.
func parseNumbersListQueryParam(r *http.Request) ([]int, error) {
	if numbersParam, ok := r.URL.Query()["numbersList"]; ok && len(numbersParam[0]) > 0 {
		var parsedNumbersList []int
		for numStr := range strings.SplitSeq(numbersParam[0], ",") {
			num, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err != nil {
				return nil, fmt.Errorf("invalid numbersList parameter: %v", err)
			}
			parsedNumbersList = append(parsedNumbersList, num)
		}
		return parsedNumbersList, nil
	}
	return nil, nil
}

// sendHTTPError logs the error and sends an HTTP error response.
func sendHTTPError(w http.ResponseWriter, errMsg string, err error, statusCode int) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, statusCode)
}
