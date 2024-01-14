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

// GetLotteryNumbers handles the request to generate lottery numbers.
func GetLotteryNumbers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	defaultNumbersList := make([]int, 40)
	for i := 0; i < 40; i++ {
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

	// Call the generateLotteryNumbers function with the validated parameters
	generatedNumbers := generator.GetNumbers(numbersList, lines, numPerLine)
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
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(version.Version))
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
		for _, numStr := range strings.Split(numbersParam[0], ",") {
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
