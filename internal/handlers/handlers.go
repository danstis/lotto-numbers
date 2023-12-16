package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danstis/lotto-numbers/internal/models" // Import the models package
)

// GetLotteryNumbers handles the request to generate lottery numbers.
func GetLotteryNumbers(w http.ResponseWriter, r *http.Request) {
	// Define your parameters for the lottery numbers
	numbersList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	lines := 5       // The number of lines of lottery numbers to generate
	numPerLine := 6  // The number of numbers per line

	// Call the generateLotteryNumbers function with the defined parameters
	generatedNumbers := generateLotteryNumbers(numbersList, lines, numPerLine)

	// Create a LotteryNumbers struct with the generated numbers
	numbers := models.LotteryNumbers{
		Lines: generatedNumbers,
	}

	// Set the Content-Type header and encode the LotteryNumbers struct to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(numbers)
}
