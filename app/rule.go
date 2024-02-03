package app

import (
	"fmt"
	"strings"

	"github.com/gomig/mig/helpers"
	"github.com/gomig/utils"
)

type Rule struct {
	Name        string              `json:"name"`
	Default     string              `json:"default"`
	Placeholder string              `json:"placeholder"`
	Description string              `json:"desc"`
	Options     []string            `json:"options"`
	Files       map[string][]string `json:"files"`
}

// IsValid check if answer is valid
func (r Rule) IsValid(answer string) bool {
	if len(r.Options) == 0 {
		return true
	}

	for _, v := range r.Options {
		if answer == v {
			return true
		}
	}
	return false
}

// Ask question from input
func (r Rule) Ask() string {
	for {
		if len(r.Options) == 0 {
			fmt.Printf("%s [default: %s]: ", r.Description, r.Default)
		} else {
			fmt.Printf("%s (%s) [default: %s]: ", r.Description, strings.Join(r.Options, "/"), r.Default)
		}
		if answer := helpers.ReadLine(r.Default); r.IsValid(answer) {
			return answer
		}
		fmt.Println("invalid answer!")
	}
}

// Ignores Get list of all files to ignore for condition answer
func (r Rule) Ignores(condition string) []string {
	// parse all files
	files := make([]string, 0)
	for _, f := range r.Files {
		files = append(files, f...)
	}

	// parse included
	included := make([]string, 0)
	if condF, ok := r.Files[condition]; ok {
		included = append(included, condF...)
	}

	// generate result
	res := make([]string, 0)
	for _, file := range files {
		if !utils.Contains(included, file) {
			res = append(res, file)
		}
	}
	return res
}
