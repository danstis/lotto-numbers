package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danstis/lotto-numbers/internal/models"
	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/stretchr/testify/assert"
)

type lotteryNumbersExpectation struct {
	wantStatusCode  int
	wantLines       int
	wantNumPerLine  int
	wantSubset      []int
	wantCoverageOf  []int
	wantErrContains string
}

func exerciseLotteryNumbersRequest(t *testing.T, query string) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest("GET", "/lottery-numbers"+query, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	http.HandlerFunc(GetLotteryNumbers).ServeHTTP(rr, req)
	return rr
}

func assertLotteryNumbersError(t *testing.T, rr *httptest.ResponseRecorder, expected lotteryNumbersExpectation) {
	t.Helper()

	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"), "content type is not text/plain; charset=utf-8")
	assert.Contains(t, rr.Body.String(), expected.wantErrContains, fmt.Sprintf("response body should contain '%s'", expected.wantErrContains))
}

func assertLotteryNumbersCoverage(t *testing.T, lines [][]int, expected []int) {
	t.Helper()

	seen := make(map[int]bool)
	for _, line := range lines {
		for _, n := range line {
			seen[n] = true
		}
	}
	for _, n := range expected {
		assert.True(t, seen[n], "number %d must appear at least once across all lines", n)
	}
}

func assertLotteryNumbersSuccess(t *testing.T, rr *httptest.ResponseRecorder, expected lotteryNumbersExpectation) {
	t.Helper()

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

	var lotteryNumbers models.LotteryNumbers
	err := json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
	assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
	assert.Len(t, lotteryNumbers.Lines, expected.wantLines, "there should be the correct number of lines of lottery numbers")
	for _, line := range lotteryNumbers.Lines {
		assert.Len(t, line, expected.wantNumPerLine, "each line should contain the correct number of numbers")
		if expected.wantSubset != nil {
			assert.Subset(t, expected.wantSubset, line, "each line should only contain numbers from the provided numbersList")
		}
	}
	if expected.wantCoverageOf != nil {
		assertLotteryNumbersCoverage(t, lotteryNumbers.Lines, expected.wantCoverageOf)
	}
}

func assertLotteryNumbersResponse(t *testing.T, rr *httptest.ResponseRecorder, expected lotteryNumbersExpectation) {
	t.Helper()

	assert.Equal(t, expected.wantStatusCode, rr.Code, "handler returned wrong status code")
	if expected.wantStatusCode != http.StatusOK {
		assertLotteryNumbersError(t, rr, expected)
		return
	}
	assertLotteryNumbersSuccess(t, rr, expected)
}

func allNumbers(start, end int) []int {
	numbers := make([]int, end-start+1)
	for i := range numbers {
		numbers[i] = start + i
	}
	return numbers
}

func TestGetLotteryNumbers(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectation lotteryNumbersExpectation
	}{
		{
			name:  "Default parameters",
			query: "",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      5,
				wantNumPerLine: 6,
			},
		},
		{
			name:  "Valid numbers list",
			query: "?numbersList=1,2,3,4,5,6,7,8,9,10",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      5,
				wantNumPerLine: 6,
				wantSubset:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
		{
			name:  "Invalid numbers list",
			query: "?numbersList=invalid",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "Bad request",
			},
		},
		{
			name:  "Valid lines parameter",
			query: "?lines=3",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      3,
				wantNumPerLine: 6,
			},
		},
		{
			name:  "Too many lines parameter",
			query: "?lines=101",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "'lines' must be <= 100",
			},
		},
		{
			name:  "Valid numPerLine parameter",
			query: "?numPerLine=5",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      5,
				wantNumPerLine: 5,
			},
		},
		{
			name:  "Non-positive lines parameter",
			query: "?lines=0",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
			},
		},
		{
			name:  "Non-positive numPerLine parameter",
			query: "?numPerLine=0",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
			},
		},
		{
			name:  "Not enough numbers for numPerLine",
			query: "?numbersList=1,2,3&numPerLine=5",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusInternalServerError,
				wantErrContains: "ensure the 'numbersList' contains enough unique numbers",
			},
		},
		{
			name:  "Zero lines requested",
			query: "?lines=0",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
			},
		},
		{
			name:  "Insufficient numbers for numPerLine",
			query: "?numbersList=1,2,3&numPerLine=4",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusInternalServerError,
				wantErrContains: "ensure the 'numbersList' contains enough unique numbers",
			},
		},
		{
			name:  "ensureAllNumbers=true with 24 numbers 4 lines 6 per line",
			query: "?numbersList=1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24&lines=4&numPerLine=6&ensureAllNumbers=true",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      4,
				wantNumPerLine: 6,
				wantCoverageOf: allNumbers(1, 24),
			},
		},
		{
			name:  "ensureAllNumbers=true infeasible 25 numbers 4 lines 6 per line",
			query: "?numbersList=1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25&lines=4&numPerLine=6&ensureAllNumbers=true",
			expectation: lotteryNumbersExpectation{
				wantStatusCode:  http.StatusBadRequest,
				wantErrContains: "cannot fit all 25 numbers",
			},
		},
		{
			name:  "ensureAllNumbers=false behaves as default",
			query: "?ensureAllNumbers=false",
			expectation: lotteryNumbersExpectation{
				wantStatusCode: http.StatusOK,
				wantLines:      5,
				wantNumPerLine: 6,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := exerciseLotteryNumbersRequest(t, tc.query)
			assertLotteryNumbersResponse(t, rr, tc.expectation)
		})
	}
}

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(VersionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedVersion := version.Version // This should be the current version set in version.go
	if rr.Body.String() != expectedVersion {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedVersion)
	}
}
