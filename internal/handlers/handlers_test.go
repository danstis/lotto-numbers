package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danstis/lotto-numbers/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetLotteryNumbers(t *testing.T) {
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
