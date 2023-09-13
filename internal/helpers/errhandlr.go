package helpers

import (
	"fmt"
	"os"
)

// Handle error
func Handle(err any) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
