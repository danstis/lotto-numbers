package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/danstis/lotto-numbers/internal/generator" // Import the generator package
	"github.com/danstis/lotto-numbers/internal/models"    // Import the models package
)

// GetLotteryNumbers handles the request to generate lottery numbers.
func GetLotteryNumbers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters with default values
	defaultNumbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	defaultLines := 5
	defaultNumPerLine := 6

	numbersList := defaultNumbersList
	if numbersParam, ok := r.URL.Query()["numbersList"]; ok && len(numbersParam[0]) > 0 {
		var parsedNumbersList []int
		for _, numStr := range strings.Split(numbersParam[0], ",") {
			num, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err != nil {
				// If there's an error, log it and use the default numbersList
				log.Printf("Error parsing numbersList: %v", err)
				parsedNumbersList = nil
				break
			}
			parsedNumbersList = append(parsedNumbersList, num)
		}
		if parsedNumbersList != nil {
			numbersList = parsedNumbersList
		}
	}

	lines := defaultLines
	if linesParam, ok := r.URL.Query()["lines"]; ok && len(linesParam[0]) > 0 {
		if parsedLines, err := strconv.Atoi(linesParam[0]); err == nil {
			lines = parsedLines
		}
	}

	numPerLine := defaultNumPerLine
	if numPerLineParam, ok := r.URL.Query()["numPerLine"]; ok && len(numPerLineParam[0]) > 0 {
		if parsedNumPerLine, err := strconv.Atoi(numPerLineParam[0]); err == nil {
			numPerLine = parsedNumPerLine
		}
	}

	// Call the generateLotteryNumbers function with the defined parameters
	generatedNumbers := generator.GetNumbers(numbersList, lines, numPerLine)

	// Create a LotteryNumbers struct with the generated numbers
	numbers := models.LotteryNumbers{
		Lines: generatedNumbers,
	}

	// Set the Content-Type header and encode the LotteryNumbers struct to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(numbers)
}
