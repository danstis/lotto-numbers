// Package models defines data structures used throughout the lotto-numbers application.
package models

// Define your models here. For example:
// LotteryNumbers represents a set of generated lottery numbers.
type LotteryNumbers struct {
	Lines [][]int `json:"lines"`
}
