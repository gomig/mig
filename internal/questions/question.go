package questions

import (
	"fmt"

	"github.com/gomig/mig/internal/helpers"
)

// Question type
type Question struct {
	Name        string
	Description string
	Default     string
	Valids      []string
}

func (q Question) isValid(answer string) bool {
	if q.Valids == nil || len(q.Valids) == 0 {
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
	answer := ""
	for {
		fmt.Printf("%s [default: %s]: ", q.Description, q.Default)
		answer = helpers.ReadLine(q.Default)
		if q.isValid(answer) {
			break
		}

		fmt.Println("invalid answer!")
	}
	return answer
}
