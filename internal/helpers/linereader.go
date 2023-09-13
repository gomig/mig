package helpers

import (
	"bufio"
	"os"
	"strings"
)

// ReadLine read input
func ReadLine(fallback string) string {
	reader := bufio.NewReader(os.Stdin)
	res, _ := reader.ReadString('\n')
	res = strings.ReplaceAll(res, "\n", "")
	res = strings.ReplaceAll(res, "\r", "")
	res = strings.TrimSpace(res)
	res = strings.ToLower(res)
	res = strings.Trim(res, " ")
	if res == "" {
		return fallback
	}
	return res
}
