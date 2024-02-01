package app

import (
	"fmt"
	"strings"

	"github.com/gomig/mig/helpers"
)

type Question struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Default     string   `json:"default"`
	Valids      []string `json:"valids"`
	Placeholder string   `json:"placeholder"`
	Falsy       string   `json:"falsy"`
	Files       []string `json:"files"`
}

// IsValid check if answer is valid
func (q Question) IsValid(answer string) bool {
	if len(q.Valids) == 0 {
		return true
	}

	for _, v := range q.Valids {
		if answer == v {
			return true
		}
	}
	return false
}

// Ask question from input
func (q Question) Ask() string {
	for {
		if len(q.Valids) == 0 {
			fmt.Printf("%s [default: %s]: ", q.Description, q.Default)
		} else {
			fmt.Printf("%s (%s) [default: %s]: ", q.Description, strings.Join(q.Valids, "/"), q.Default)
		}
		if answer := helpers.ReadLine(q.Default); q.IsValid(answer) {
			return answer
		}
		fmt.Println("invalid answer!")
	}
}
