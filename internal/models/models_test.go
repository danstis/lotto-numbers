package models

import (
	"reflect"
	"testing"
)

func TestLotteryNumbers_Structure(t *testing.T) {
	expected := [][]int{{1, 2, 3, 4, 5, 6}}
	lotteryNumbers := LotteryNumbers{
		Lines: expected,
	}

	if !reflect.DeepEqual(lotteryNumbers.Lines, expected) {
		t.Errorf("LotteryNumbers.Lines = %v, want %v", lotteryNumbers.Lines, expected)
	}
}
