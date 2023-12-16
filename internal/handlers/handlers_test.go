package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danstis/lotto-numbers/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetLotteryNumbersDefaultParams(t *testing.T) {
	req, err := http.NewRequest("GET", "/lottery-numbers", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLotteryNumbers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
	assert.Len(t, lotteryNumbers.Lines, 5, "there should be 5 lines of lottery numbers")
	for _, line := range lotteryNumbers.Lines {
		assert.Len(t, line, 6, "each line should contain 6 numbers")
	}
}
func TestGetLotteryNumbersWithValidNumbersList(t *testing.T) {
	req, err := http.NewRequest("GET", "/lottery-numbers?numbersList=1,2,3,4,5,6,7,8,9,10", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLotteryNumbers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
	assert.Len(t, lotteryNumbers.Lines, 5, "there should be 5 lines of lottery numbers")
	for _, line := range lotteryNumbers.Lines {
		assert.Len(t, line, 6, "each line should contain 6 numbers")
		assert.Subset(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, line, "each line should only contain numbers from the provided numbersList")
	}
}

func TestGetLotteryNumbersWithInvalidNumbersList(t *testing.T) {
	req, err := http.NewRequest("GET", "/lottery-numbers?numbersList=invalid", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLotteryNumbers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code even with invalid numbersList")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON despite invalid numbersList")
	assert.Len(t, lotteryNumbers.Lines, 5, "there should be 5 lines of lottery numbers by default")
	for _, line := range lotteryNumbers.Lines {
		assert.Len(t, line, 6, "each line should contain 6 numbers by default")
	}
}

func TestGetLotteryNumbersWithValidLinesParam(t *testing.T) {
	req, err := http.NewRequest("GET", "/lottery-numbers?lines=3", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLotteryNumbers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
	assert.Len(t, lotteryNumbers.Lines, 3, "there should be 3 lines of lottery numbers as specified by the lines parameter")
}

func TestGetLotteryNumbersWithValidNumPerLineParam(t *testing.T) {
	req, err := http.NewRequest("GET", "/lottery-numbers?numPerLine=5", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLotteryNumbers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
	for _, line := range lotteryNumbers.Lines {
		assert.Len(t, line, 5, "each line should contain 5 numbers as specified by the numPerLine parameter")
	}
}
