package app

import (
	"fmt"

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
		return answer != ""
	} else {
		for _, v := range r.Options {
			if answer == v {
				return true
			}
		}
	}
	return false
}

// Ask question from input
func (r Rule) Ask() string {
	for {
		opts := ""
		def := ""
		if len(r.Options) > 0 {
			opts = opts + helpers.Format(" (", helpers.GRAY)
			for i, opt := range r.Options {
				opts = opts + helpers.Format(opt, helpers.ITALIC, helpers.BLUE)
				if i+1 < len(r.Options) {
					opts = opts + helpers.Format("/", helpers.GRAY)
				}
			}
			opts = opts + helpers.Format(")", helpers.GRAY)
		}
		if r.Default != "" {
			def = helpers.Format(" [", helpers.GRAY) +
				helpers.Format(r.Default, helpers.BOLD, helpers.GREEN) +
				helpers.Format("]", helpers.GRAY)
		}
		fmt.Printf("%s%s%s: ", utils.If(r.Description != "", r.Description, r.Name), opts, def)
		if answer := helpers.ReadLine(r.Default); r.IsValid(answer) {
			return answer
		}
		if len(r.Options) == 0 {
			fmt.Println(helpers.Format("Enter value", helpers.RED))
		} else {
			fmt.Println(helpers.Format("Invalid value", helpers.RED))
		}
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
